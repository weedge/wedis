package srv

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/cloudwego/kitex/pkg/klog"
	"github.com/cloudwego/kitex/pkg/utils"
	"github.com/hertz-contrib/obs-opentelemetry/provider"
	"github.com/hertz-contrib/obs-opentelemetry/tracing"
	"github.com/tidwall/redcon"
	"github.com/weedge/wedis/internal/srv/config"
	"github.com/weedge/wedis/internal/srv/storager"
	"github.com/weedge/wedis/pkg/utils/logutils"
)

type Server struct {
	opts          *config.ServerOptions
	kitexKVLogger logutils.IKitexZapKVLogger
	// redcon server handler
	mux *redcon.ServeMux
	// redcon server
	redconSrv *redcon.Server
	// store
	store *storager.Storager
}

// Run reviews server
func (s *Server) Run(ctx context.Context) error {
	klog.Debugf("server opts: %s", s.opts)

	klog.SetLogger(s.kitexKVLogger)
	klog.SetLevel(s.opts.LogLevel.KitexLogLevel())

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
	s.InitStore()

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
	if s.store != nil {
		if err := s.store.Close(); err != nil {
			klog.Errorf("close store err: %s", err.Error())
		}
	}

	if s.redconSrv != nil {
		if err := s.redconSrv.Close(); err != nil {
			klog.Errorf("close RESP cmd server err: %s", err.Error())
		}
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
	if dbIdx < 0 || dbIdx >= s.opts.StoreOpts.Databases {
		dbIdx = 0
	}
	cli := new(ConnClient)
	cli.SetSrv(s)
	db, err := s.store.Select(dbIdx)
	if err != nil {
		return nil
	}
	cli.SetDb(db)
	return cli
}

func (s *Server) InitRespCmdService() {
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

	s.registerRespConnClient()

	listenErrSignal := make(chan error)
	go func() {
		err := s.redconSrv.ListenServeAndSignal(listenErrSignal)
		if err != nil {
			klog.Fatal(err)
		}
	}()
	err := <-listenErrSignal
	if err != nil {
		klog.Errorf("resp cmd server listen err:%s", err.Error())
		return
	}
	klog.Infof("resp cmd server listening on address=%s", s.opts.RespCmdSrvOpts.Addr)
}

func (s *Server) InitStore() {
	store, err := storager.Open(&s.opts.StoreOpts)
	if err != nil {
		klog.Fatalf("open store err:%s", err.Error())
	}
	s.store = store
}
