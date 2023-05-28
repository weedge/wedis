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
	RespCmdSrvOpts            RespCmdServiceOptins   `mapstructure:"respCmdSrv"`
	StoreOpts                 StorgerOptions         `mapstructure:"storeOpts"`
}

// DefaultServerOptions default opts
func DefaultServerOptions() *ServerOptions {
	return &ServerOptions{
		//OltpGrpcCollectorEndpoint: ":4317",
		//HttpAddr:       ":8110",
		ProjectName:    "wedis",
		LogLevel:       logutils.LevelDebug,
		LogMeta:        map[string]interface{}{},
		RespCmdSrvOpts: *DefaultRespCmdServiceOptins(),
		StoreOpts:      *DefaultStoragerOptions(),
	}
}

type RespCmdServiceOptins struct {
	Addr                  string `mapstructure:"addr"`
	AuthPassword          string `mapstructure:"authPassword"`
	ConnKeepaliveInterval int    `mapstructure:"connKeepaliveInterval"`
}

func DefaultRespCmdServiceOptins() *RespCmdServiceOptins {
	return &RespCmdServiceOptins{}
}

type StorgerOptions struct {
	DataDir          string `mapstructure:"dataDir"`
	Databases        int    `mapstructure:"databases"`
	KVStoreName      string `mapstructure:"kvStoreName"`
	DBPath           string `mapstructure:"dbPath"`
	DBSyncCommit     int    `mapstructure:"dbSyncCommit"`
	TTLCheckInterval int    `mapstructure:"ttlCheckInterval"`
}

func DefaultStoragerOptions() *StorgerOptions {
	return &StorgerOptions{
		DataDir:     DefaultDataDir,
		Databases:   DefaultDatabases,
		KVStoreName: DefaultKVStoreName,
	}
}
