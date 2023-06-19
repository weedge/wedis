package config

import (
	"github.com/weedge/pkg/utils/logutils"
	storager "github.com/weedge/xdis-storager"
)

// ServerOptions server options
type ServerOptions struct {
	HttpAddr                  string                 `mapstructure:"httpAddr"`
	LogLevel                  logutils.Level         `mapstructure:"logLevel"`
	ProjectName               string                 `mapstructure:"projectName"`
	LogMeta                   map[string]interface{} `mapstructure:"logMeta"`
	OltpGrpcCollectorEndpoint string                 `mapstructure:"oltpCollectorGrpcEndpoint"`
	StoragerName              string                 `mapstructure:"storagerName"`

	RespCmdSrvOpts RespCmdServiceOptins     `mapstructure:"respCmdSrv"`
}

// DefaultServerOptions default opts
func DefaultServerOptions() *ServerOptions {
	return &ServerOptions{
		//OltpGrpcCollectorEndpoint: ":4317",
		//HttpAddr:       ":8110",
		ProjectName:  "wedis",
		LogLevel:     logutils.LevelDebug,
		LogMeta:      map[string]interface{}{},
		StoragerName: storager.RegisterStoragerName,

		RespCmdSrvOpts: *DefaultRespCmdServiceOptins(),
	}
}

type RespCmdServiceOptins struct {
	Addr                  string `mapstructure:"addr"`
	AuthPassword          string `mapstructure:"authPassword"`
	ConnKeepaliveInterval int    `mapstructure:"connKeepaliveInterval"`
	Mode                  string `mapstructure:"mode"`
}

func DefaultRespCmdServiceOptins() *RespCmdServiceOptins {
	return &RespCmdServiceOptins{Mode: ModeStandalone}
}
