package config

import (
	goleveldb "github.com/weedge/openkv-goleveldb"
	"github.com/weedge/pkg/configparser"
)

type Options struct {
	Server       *ServerOptions          `mapstructure:"server"`
	GoLeveldbCfg goleveldb.LevelDBConfig `mapstructure:"goLeveldbCfg"`
}

// DefaultOptions default opts
func DefaultOptions() *Options {
	return &Options{
		Server:       DefaultServerOptions(),
		GoLeveldbCfg: *goleveldb.DefaultLevelDBConfig(),
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
