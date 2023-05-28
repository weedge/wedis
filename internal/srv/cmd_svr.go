package srv

import (
	"context"
)

func init() {
	RegisterCmd("ping", ping)
}

func ping(ctx context.Context, c *ConnClient, cmdParams [][]byte) (res interface{}, err error) {
	res = "PONG"
	return
}
