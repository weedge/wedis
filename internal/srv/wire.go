//go:build wireinject
// +build wireinject

package srv

import (
	"context"

	"github.com/google/wire"
	"github.com/tidwall/redcon"

	goleveldb "github.com/weedge/openkv-goleveldb"
	"github.com/weedge/pkg/configparser"
	driver "github.com/weedge/pkg/driver/openkv"
	"github.com/weedge/pkg/utils/logutils"
	"github.com/weedge/wedis/internal/srv/config"
	"github.com/weedge/wedis/internal/srv/storager"
)

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
		),

		wire.FieldsOf(new(*config.ServerOptions), "LogLevel", "LogMeta", "StoreOpts"),
		logutils.NewkitexZapKVLogger,
		redcon.NewServeMux,
		storager.Open,

		wire.Struct(new(Server), "opts", "kitexKVLogger", "mux", "store"),
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
		wire.Bind(new(driver.IStore), new(*goleveldb.Store)),
		driver.Register,
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
		wire.Bind(new(driver.IStore), new(*goleveldb.Store)),
		driver.Register,
	))
}

func ProvideOpts(
	c goleveldb.Option,
) []goleveldb.Option {
	return []goleveldb.Option{
		c,
	}
}
