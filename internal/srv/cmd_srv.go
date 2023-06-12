package srv

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/tidwall/redcon"
	"github.com/weedge/pkg/utils"
)

func init() {
	RegisterCmd("client", client)
	RegisterCmd("echo", echo)
	RegisterCmd("hello", hello)
	RegisterCmd("ping", ping)
	RegisterCmd("select", selectCmd)
	RegisterCmd("flushdb", flushdb)
	RegisterCmd("flushall", flushall)
}

func authUser(ctx context.Context, c *ConnClient, pwd string) (err error) {
	return
}

func client(ctx context.Context, c *ConnClient, cmdParams [][]byte) (res interface{}, err error) {
	if len(cmdParams) < 1 {
		return nil, ErrCmdParams
	}
	op := strings.ToLower(utils.Bytes2String(cmdParams[0]))
	switch op {
	case "getname":
		res = c.name
	default:
		//todo
	}

	return
}

// just hello cmd, no resp protocol change
func hello(ctx context.Context, c *ConnClient, cmdParams [][]byte) (res interface{}, err error) {
	if len(cmdParams) > 6 {
		op := cmdParams[len(cmdParams)-1]
		return nil, fmt.Errorf("%s in %s option '%s'", ErrSyntax.Error(), "HELLO", op)
	}

	/*
		//golang map u know :)
		data := map[string]any{
			"server": "redis",
			"proto":  redcon.SimpleInt(2),
			"mode":   c.srv.opts.RespCmdSrvOpts.Mode,
		}
	*/
	data := []any{
		"server", "redis",
		"proto", redcon.SimpleInt(2),
		"mode", c.srv.opts.RespCmdSrvOpts.Mode,
	}
	res = data
	if len(cmdParams) == 0 {
		return
	}

	protocalVer, err := strconv.ParseInt(utils.Bytes2String(cmdParams[0]), 10, 64)
	if err != nil {
		return nil, ErrProtocalVer
	}
	if protocalVer < 2 || protocalVer > 3 {
		return nil, ErrUnsupportVer
	}

	for nextArg := 1; nextArg < len(cmdParams); nextArg++ {
		moreArgs := len(cmdParams) - nextArg - 1
		op := strings.ToLower(utils.Bytes2String(cmdParams[nextArg]))
		if op == "auth" && moreArgs > 0 && moreArgs%2 == 0 {
			nextArg++
			if strings.ToLower(utils.Bytes2String(cmdParams[nextArg])) != "default" {
				return nil, ErrInvalidPwd
			}
			nextArg++
			pwd := utils.Bytes2String(cmdParams[nextArg])
			if err = authUser(ctx, c, pwd); err != nil {
				return nil, err
			}
			//println("auth", nextArg)
		} else if op == "setname" && moreArgs > 0 {
			nextArg++
			c.SetConnName(utils.Bytes2String(cmdParams[nextArg]))
			//println("setname", nextArg)
		} else {
			//println("other", moreArgs, nextArg, op)
			if moreArgs == 3 {
				op = "setname"
			}
			return nil, fmt.Errorf("%s in %s option '%s'", ErrSyntax.Error(), "HELLO", op)
		}
	}

	return
}

func ping(ctx context.Context, c *ConnClient, cmdParams [][]byte) (res interface{}, err error) {
	if len(cmdParams) > 1 {
		return nil, ErrCmdParams
	}
	if len(cmdParams) == 1 {
		res = cmdParams[0]
		return
	}
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
