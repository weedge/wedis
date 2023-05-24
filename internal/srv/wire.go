//go:build wireinject
// +build wireinject

package srv

import (
	"context"

	"github.com/google/wire"
	"github.com/tidwall/redcon"
	"github.com/weedge/wedis/pkg/configparser"
	"github.com/weedge/wedis/pkg/utils/logutils"
)

// NewServer build server with wire, dependency obj inject, so init random
func NewServer(ctx context.Context) (*Server, error) {
	panic(wire.Build(
		configparser.Default,
		Configure,
		wire.FieldsOf(new(*Options),
			"Server",
		),

		wire.FieldsOf(new(*ServerOptions), "LogLevel", "LogMeta"),
		logutils.NewkitexZapKVLogger,
		redcon.NewServeMux,

		wire.Struct(new(Server), "opts", "kitexKVLogger", "mux"),
	))
}
