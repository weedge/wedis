//go:build wireinject
// +build wireinject

package srv

import (
	"context"

	"github.com/google/wire"
	goleveldb "github.com/weedge/openkv-goleveldb"
	"github.com/weedge/pkg/configparser"
	"github.com/weedge/pkg/driver"
	openkvDriver "github.com/weedge/pkg/driver/openkv"
	"github.com/weedge/pkg/utils/logutils"
	"github.com/weedge/wedis/internal/srv/config"
	standalone "github.com/weedge/xdis-standalone"
	replica "github.com/weedge/xdis-replica-storager"
	storager "github.com/weedge/xdis-storager"
	xdistikv "github.com/weedge/xdis-tikv"
)

// notice: multiple bindings
// https://github.com/google/wire/blob/main/docs/faq.md#what-if-my-dependency-graph-has-two-dependencies-of-the-same-type
// build server with wire, dependency obj inject, so init random

// ParserCfg config paser
func ParserCfg() (*config.Options, error) {
	panic(wire.Build(
		configparser.Default,
		config.Configure,
	))
}

// NewServer server init
func NewServer(context.Context, *config.Options) (*Server, error) {
	panic(wire.Build(
		wire.FieldsOf(new(*config.Options),
			"Server",
			"RespCmdSrvName",
			"StoragerName",
		),

		wire.FieldsOf(new(*config.ServerOptions), "LogLevel", "LogMeta"),
		logutils.NewkitexZapKVLogger,
		driver.GetRespCmdSrv,
		driver.GetStorager,

		wire.Struct(new(Server), "opts", "kitexKVLogger", "respSrv", "store"),
	))
}

func RegisterStandaloneRespCmdSrv(*config.Options) error {
	panic(wire.Build(
		wire.FieldsOf(new(*config.Options),
			"StandaloneRespCmdSrvCfg",
		),

		standalone.New,
		wire.Bind(new(driver.IRespService), new(*standalone.RespCmdService)),
		driver.RegisterRespCmdSrv,
	))
}

func RegisterReplicaMasterSlaveRespCmdSrv(*config.Options) error {
	panic(wire.Build(
		wire.FieldsOf(new(*config.Options),
			"StandaloneRespCmdSrvCfg",
			"RplMasterSlaveRespCmdSrvCfg",
		),

		replica.New,
		wire.Bind(new(driver.IRespService), new(*replica.RespCmdService)),
		driver.RegisterRespCmdSrv,
	))
}

// RegisterXdisStorager register storager xdis-storager
func RegisterXdisStorager(*config.Options) error {
	panic(wire.Build(
		wire.FieldsOf(new(*config.Options),
			"StoreCfg",
		),

		storager.New,
		wire.Bind(new(driver.IStorager), new(*storager.Storager)),
		driver.RegisterStorager,
	))
}

// RegisterXdisTikv register storager xdis-tikv
func RegisterXdisTikv(*config.Options) error {
	panic(wire.Build(
		wire.FieldsOf(new(*config.Options),
			"TikvCfg",
		),

		xdistikv.New,
		wire.Bind(new(driver.IStorager), new(*xdistikv.Storager)),
		driver.RegisterStorager,
	))
}

// RegisterGoleveldb register kv store engine goleveldb
func RegisterGoleveldb(*config.Options) error {
	panic(wire.Build(
		wire.FieldsOf(new(*config.Options),
			"GoLeveldbCfg",
		),

		goleveldb.WithConfig,
		ProvideOpts,
		wire.Value(goleveldb.StoreTypeDB),
		goleveldb.New,
		wire.Bind(new(openkvDriver.IStore), new(*goleveldb.Store)),
		openkvDriver.Register,
	))
}

// RegisterMemGoleveldb regisiter kv store engine memory goleveldb
func RegisterMemGoleveldb(*config.Options) error {
	panic(wire.Build(
		wire.FieldsOf(new(*config.Options),
			"GoLeveldbCfg",
		),

		goleveldb.WithConfig,
		ProvideOpts,
		wire.Value(goleveldb.StoreTypeMemory),
		goleveldb.New,
		wire.Bind(new(openkvDriver.IStore), new(*goleveldb.Store)),
		openkvDriver.Register,
	))
}

func ProvideOpts(
	c goleveldb.Option,
) []goleveldb.Option {
	return []goleveldb.Option{
		c,
	}
}
