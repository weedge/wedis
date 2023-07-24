package srv

import (
	"context"

	"github.com/cloudwego/kitex/pkg/klog"
	"github.com/spf13/cobra"
	"github.com/weedge/pkg/driver"
	openkvDriver "github.com/weedge/pkg/driver/openkv"
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
			srv.ConfigOpts = opts

			// register log store local engine
			srv.RegisterLogStoreGoleveldb(opts)

			// register kv store local engine
			srv.RegisterGoleveldb(opts)
			srv.RegisterMemGoleveldb(opts)

			// register resp cmd kv storager
			srv.RegisterXdisStorager(opts)
			srv.RegisterXdisTikv(opts)

			// register resp cmd service
			srv.RegisterStandaloneRespCmdSrv(opts)
			srv.RegisterReplicaMasterSlaveRespCmdSrv(opts)

			klog.Infof("register resp cmd services: %+v, current used service: %s",
				driver.ListRespCmdSrvs(), opts.RespCmdSrvName)
			klog.Infof("register storagers: %+v, current used storager: %s",
				driver.ListStoragers(), opts.StoragerName)
			klog.Infof("register store engines: %+v, current used kvstore engine: %s, logstore engine: %s",
				openkvDriver.ListStores(),
				opts.StoreCfg.KVStoreName,
				opts.RplMasterSlaveRespCmdSrvCfg.LogStoreOpenkvCfg.KVStoreName)

			server, err := srv.NewServer(ctx, opts)
			if err != nil {
				return err
			}
			return server.Run(ctx)
		},
	}
}
