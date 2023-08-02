package config

import (
	"github.com/weedge/pkg/utils/logutils"
)

// ServerOptions server options
type ServerOptions struct {
	HttpAddr                  string                 `mapstructure:"httpAddr"`
	LogLevel                  logutils.Level         `mapstructure:"logLevel"`
	ProjectName               string                 `mapstructure:"projectName"`
	LogMeta                   map[string]interface{} `mapstructure:"logMeta"`
	OltpGrpcCollectorEndpoint string                 `mapstructure:"oltpCollectorGrpcEndpoint"`
	PProfPort                 string                 `mapstructure:"pprofPort"`
}

// DefaultServerOptions default opts
func DefaultServerOptions() *ServerOptions {
	return &ServerOptions{
		//OltpGrpcCollectorEndpoint: ":4317",
		//HttpAddr:       ":8110",
		ProjectName: "wedis",
		LogLevel:    logutils.LevelDebug,
		LogMeta:     map[string]interface{}{},
		PProfPort:   "2222",
	}
}
