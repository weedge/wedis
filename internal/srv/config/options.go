package config

import (
	goleveldb "github.com/weedge/openkv-goleveldb"
	"github.com/weedge/pkg/configparser"
	"github.com/weedge/xdis-storager/config"
	tikvconf "github.com/weedge/xdis-tikv/config"
)

type Options struct {
	Server       ServerOptions            `mapstructure:"server"`
	StoreCfg     config.StorgerOptions    `mapstructure:"storeCfg"`
	TikvCfg      tikvconf.StoragerOptions `mapstructure:"tikvStoreCfg"`
	GoLeveldbCfg goleveldb.LevelDBConfig  `mapstructure:"goLeveldbCfg"`
}

// DefaultOptions default opts
func DefaultOptions() *Options {
	return &Options{
		Server:       *DefaultServerOptions(),
		GoLeveldbCfg: *goleveldb.DefaultLevelDBConfig(),
		StoreCfg:     *config.DefaultStoragerOptions(),
		TikvCfg:      *tikvconf.DefaultStoragerOptions(),
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
