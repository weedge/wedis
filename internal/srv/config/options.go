package config

import (
	goleveldb "github.com/weedge/openkv-goleveldb"
	"github.com/weedge/pkg/configparser"
	"github.com/weedge/pkg/driver"
	replicaconf "github.com/weedge/xdis-replica-storager/config"
	standaloneconf "github.com/weedge/xdis-standalone/config"
	storager "github.com/weedge/xdis-storager"
	"github.com/weedge/xdis-storager/config"
	tikvconf "github.com/weedge/xdis-tikv/config"
)

type Options struct {
	StoragerName                string                               `mapstructure:"storagerName"`
	RespCmdSrvName              driver.RespServiceName               `mapstructure:"respCmdSrvName"`
	Server                      ServerOptions                        `mapstructure:"server"`
	StandaloneRespCmdSrvCfg     standaloneconf.RespCmdServiceOptions `mapstructure:"standaloneRespCmdSrvCfg"`
	RplMasterSlaveRespCmdSrvCfg replicaconf.RespCmdServiceOptions    `mapstructure:"rplMasterSlaveRespCmdSrvCfg"`
	StoreCfg                    config.StorgerOptions                `mapstructure:"storeCfg"`
	TikvCfg                     tikvconf.StoragerOptions             `mapstructure:"tikvStoreCfg"`
	GoLeveldbCfg                goleveldb.LevelDBConfig              `mapstructure:"goLeveldbCfg"`
}

// DefaultOptions default opts
func DefaultOptions() *Options {
	return &Options{
		StoragerName:                storager.RegisterStoragerName,
		RespCmdSrvName:              standaloneconf.RegisterRespSrvModeName,
		Server:                      *DefaultServerOptions(),
		GoLeveldbCfg:                *goleveldb.DefaultLevelDBConfig(),
		StoreCfg:                    *config.DefaultStoragerOptions(),
		TikvCfg:                     *tikvconf.DefaultStoragerOptions(),
		StandaloneRespCmdSrvCfg:     *standaloneconf.DefaultRespCmdServiceOptions(),
		RplMasterSlaveRespCmdSrvCfg: *replicaconf.DefaultRespCmdServiceOptions(),
	}
}

// Configure inject config
func Configure(configProvider configparser.Provider) (*Options, error) {
	opt := DefaultOptions()

	cp, err := configProvider.Get()
	if err != nil {
		return nil, err
	}

	if err = cp.UnmarshalExact(opt); err != nil {
		return nil, err
	}

	return opt, nil
}
