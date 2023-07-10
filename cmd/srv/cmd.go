package srv

import (
	"context"

	"github.com/cloudwego/kitex/pkg/klog"
	"github.com/spf13/cobra"
	"github.com/weedge/wedis/internal/srv"
)

func NewCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "srv",
		Short: "start wedis server",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			opts, err := srv.ParserCfg()
			if err != nil {
				return err
			}
			klog.Infof("config:%+v", opts)

			// register store local engine
			srv.RegisterGoleveldb(opts)
			srv.RegisterMemGoleveldb(opts)

			// register resp cmd kv storager
			srv.RegisterXdisStorager(opts)
			srv.RegisterXdisTikv(opts)

			// register resp cmd service
			srv.RegisterStandaloneRespCmdSrv(opts)

			server, err := srv.NewServer(ctx, opts)
			if err != nil {
				return err
			}
			return server.Run(ctx)
		},
	}
}
