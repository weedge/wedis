package srv

import (
	"context"
	"strconv"

	"github.com/weedge/pkg/utils"
)

func init() {
	RegisterCmd("echo", echo)
	RegisterCmd("ping", ping)
	RegisterCmd("select", selectCmd)
	RegisterCmd("flushdb", flushdb)
	RegisterCmd("flushall", flushall)
}

func ping(ctx context.Context, c *ConnClient, cmdParams [][]byte) (res interface{}, err error) {
	res = PONG
	return
}

func echo(ctx context.Context, c *ConnClient, cmdParams [][]byte) (res interface{}, err error) {
	if len(cmdParams) != 1 {
		return nil, ErrCmdParams
	}

	res = cmdParams[0]
	return
}

func selectCmd(ctx context.Context, c *ConnClient, cmdParams [][]byte) (res interface{}, err error) {
	if len(cmdParams) != 1 {
		return nil, ErrCmdParams
	}

	index, err := strconv.Atoi(utils.Bytes2String(cmdParams[0]))
	if err != nil {
		return
	}

	db, err := c.srv.store.Select(ctx, index)
	if err != nil {
		return
	}
	c.db = db

	res = OK
	return
}

func flushdb(ctx context.Context, c *ConnClient, cmdParams [][]byte) (res interface{}, err error) {
	_, err = c.db.FlushDB(ctx)
	if err != nil {
		return
	}

	res = OK
	return
}

func flushall(ctx context.Context, c *ConnClient, cmdParams [][]byte) (res interface{}, err error) {
	err = c.srv.store.FlushAll(ctx)
	if err != nil {
		return
	}

	res = OK
	return
}
