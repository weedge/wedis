package srv

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/cloudwego/kitex/pkg/klog"
	"github.com/hertz-contrib/obs-opentelemetry/provider"
	"github.com/hertz-contrib/obs-opentelemetry/tracing"
	"github.com/weedge/pkg/driver"
	openkvDriver "github.com/weedge/pkg/driver/openkv"
	"github.com/weedge/pkg/utils/logutils"
	"github.com/weedge/wedis/internal/srv/config"
)

type Server struct {
	opts          *config.ServerOptions
	kitexKVLogger logutils.IKitexZapKVLogger
	// storage
	store driver.IStorager
	// RESPService
	respSrv driver.IRespService
}

// Run reviews server
func (s *Server) Run(ctx context.Context) error {
	klog.SetLogger(s.kitexKVLogger)
	klog.SetLevel(s.opts.LogLevel.KitexLogLevel())

	klog.Infof("register storagers: %+v", driver.ListStoragers())
	klog.Infof("register store engines: %+v", openkvDriver.ListStores())
	klog.Infof("server opts: %+v", s.opts)

	defer s.Stop()

	// open storager
	if err := s.store.Open(ctx); err != nil {
		return err
	}

	// set storager, start RESP service
	s.respSrv.SetStorager(s.store)
	if err := s.respSrv.Start(ctx); err != nil {
		return err
	}

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

	if s.opts.HttpAddr == "" {
		signalToNotify := []os.Signal{syscall.SIGINT, syscall.SIGHUP, syscall.SIGTERM}
		if signal.Ignored(syscall.SIGHUP) {
			signalToNotify = []os.Signal{syscall.SIGINT, syscall.SIGTERM}
		}
		sig := make(chan os.Signal, 1)
		signal.Notify(sig, signalToNotify...)
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
		s.store = nil
	}

	if s.respSrv != nil {
		if err := s.respSrv.Close(); err != nil {
			klog.Errorf("close RESP cmd server err: %s", err.Error())
		}
		s.respSrv = nil
	}

	klog.Infof("server stop")
}
