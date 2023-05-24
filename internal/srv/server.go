package srv

import (
	"context"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/cloudwego/kitex/pkg/klog"
	"github.com/cloudwego/kitex/pkg/utils"
	"github.com/hertz-contrib/obs-opentelemetry/provider"
	"github.com/hertz-contrib/obs-opentelemetry/tracing"
	"github.com/tidwall/redcon"
	"github.com/weedge/wedis/internal/srv/storager"
	"github.com/weedge/wedis/pkg/utils/logutils"
)

type Server struct {
	opts          *ServerOptions
	kitexKVLogger logutils.IKitexZapKVLogger
	// redcon server handler
	mux *redcon.ServeMux
	// redcon server
	redconSrv *redcon.Server
	// dbs
	dbs map[int]*storager.DB
	// dbs map rw lock r>w
	rwMu sync.RWMutex
}

// ServerOptions server options
type ServerOptions struct {
	HttpAddr                  string                 `mapstructure:"httpAddr"`
	LogLevel                  logutils.Level         `mapstructure:"logLevel"`
	ProjectName               string                 `mapstructure:"projectName"`
	LogMeta                   map[string]interface{} `mapstructure:"logMeta"`
	OltpGrpcCollectorEndpoint string                 `mapstructure:"oltpCollectorGrpcEndpoint"`
	RespCmdSrvOpts            RespCmdServiceOptins   `mapstructure:"respCmdSrv"`
}

// DefaultServerOptions default opts
func DefaultServerOptions() *ServerOptions {
	return &ServerOptions{
		//OltpGrpcCollectorEndpoint: ":4317",
		HttpAddr:       ":8110",
		ProjectName:    "wedis",
		LogLevel:       logutils.LevelDebug,
		LogMeta:        map[string]interface{}{},
		RespCmdSrvOpts: *DefaultRespCmdServiceOptins(),
	}
}

type RespCmdServiceOptins struct {
	Addr         string `mapstructure:"addr"`
	AuthPassword string `mapstructure:"authPassword"`

	ReplicaOf string `mapstructure:"replicaOf"`

	Readonly bool `mapstructure:"readOnly"`

	DataDir   string `mapstructure:"dataDir"`
	Databases int    `mapstructure:"databases"`

	DBName       string `mapstructure:"dbName"`
	DBPath       string `mapstructure:"dbPath"`
	DBSyncCommit int    `mapstructure:"dbSyncCommit"`

	ConnReadBufferSize    int `mapstructure:"connReadBufferSize"`
	ConnWriteBufferSize   int `mapstructure:"connWriteBufferSize"`
	ConnKeepaliveInterval int `mapstructure:"connKeepaliveInterval"`

	TTLCheckInterval int `mapstructure:"ttlCheckInterval"`
}

func DefaultRespCmdServiceOptins() *RespCmdServiceOptins {
	return &RespCmdServiceOptins{}
}

// Run reviews server
func (s *Server) Run(ctx context.Context) error {
	klog.Debugf("server opts: %s", s.opts)

	klog.SetLogger(s.kitexKVLogger)
	klog.SetLevel(s.opts.LogLevel.KitexLogLevel())

	s.dbs = make(map[int]*storager.DB)
	s.dbs[0] = storager.New()

	defer s.Stop()

	if s.opts.OltpGrpcCollectorEndpoint != "" {
		tracingProvider := provider.NewOpenTelemetryProvider(
			provider.WithExportEndpoint(s.opts.OltpGrpcCollectorEndpoint),
			provider.WithEnableMetrics(true),
			provider.WithEnableTracing(true),
			provider.WithServiceName(s.opts.ProjectName),
			provider.WithInsecure(),
		)
		defer func(ctx context.Context, p provider.OtelProvider) {
			_ = p.Shutdown(ctx)
		}(ctx, tracingProvider)
	}

	s.InitRespCmdService()

	if len(s.opts.HttpAddr) == 0 {
		sig := make(chan os.Signal, 1)
		signal.Notify(sig, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
		<-sig
		return nil
	}

	tracer, cfg := tracing.NewServerTracer()
	h := server.Default(
		tracer,
		server.WithHostPorts(s.opts.HttpAddr),
	)
	h.Use(tracing.ServerMiddleware(cfg))
	s.SetupRoutes(h)
	h.Spin()

	return nil
}

func (s *Server) Stop() {
	for _, db := range s.dbs {
		if err := db.Close(); err != nil {
			klog.Errorf("close db err: %v", err)
		}
	}

	if err := s.redconSrv.Close(); err != nil {
		klog.Errorf("close RESP cmd server err: %v", err)
	}

	klog.Infof("server stop")
}

func (s *Server) registerRespConnClient() {
	for cmdOp := range SupportedCommands {
		s.mux.HandleFunc(cmdOp, func(conn redcon.Conn, cmd redcon.Command) {
			cmdOp := utils.SliceByteToString(cmd.Args[0])
			params := [][]byte{}
			if len(cmd.Args) > 0 {
				params = cmd.Args[1:]
			}

			switch cmdOp {
			case "quit":
				err := conn.Close()
				if err != nil {
					klog.Errorf("resp cmd quit connect close err: %s", err.Error())
				}
				return
			}

			cli, ok := conn.Context().(*ConnClient)
			if !ok {
				klog.Errorf("resp cmd connect client init err")
				return
			}
			ctx := context.WithValue(context.Background(), RespCmdCtxKey, conn.Context())
			res, err := cli.DoCmd(ctx, cmdOp, params)
			if err != nil {
				conn.WriteError(err.Error())
				return
			}
			conn.WriteAny(res)
		})
	}
}

func (s *Server) InitConnClient(dbIdx int) *ConnClient {
	if dbIdx < 0 || dbIdx >= s.opts.RespCmdSrvOpts.Databases {
		dbIdx = 0
	}
	cli := new(ConnClient)
	cli.SetSrv(s)
	s.rwMu.RLock()
	cli.SetDb(s.dbs[dbIdx])
	s.rwMu.RUnlock()
	return cli
}

func (s *Server) InitRespCmdService() {
	s.registerRespConnClient()

	//RESP cmd tcp server
	s.redconSrv = redcon.NewServer(s.opts.RespCmdSrvOpts.Addr, s.mux.ServeRESP,
		func(conn redcon.Conn) bool {
			// use this function to accept (return true) or deny the connection (return false).
			// set ctx
			klog.Infof("accept: %s", conn.RemoteAddr())
			conn.SetContext(s.InitConnClient(0))
			return true
		},
		func(conn redcon.Conn, err error) {
			// this is called when the connection has been closed
			klog.Infof("closed: %s, err: %v", conn.RemoteAddr(), err)
		},
	)
	if s.opts.RespCmdSrvOpts.ConnKeepaliveInterval > 0 {
		s.redconSrv.SetIdleClose(time.Duration(s.opts.RespCmdSrvOpts.ConnKeepaliveInterval))
	}

	go func() {
		err := s.redconSrv.ListenAndServe()
		if err != nil {
			klog.Fatal(err)
		}
	}()
}
