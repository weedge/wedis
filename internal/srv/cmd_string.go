package srv

import (
	"context"
	"strconv"

	"github.com/weedge/pkg/driver"
	"github.com/weedge/pkg/utils"
)

func init() {
	RegisterCmd(CmdTypeString,"append", appendCmd)
	RegisterCmd(CmdTypeString,"decr", decr)
	RegisterCmd(CmdTypeString,"decrby", decrby)
	RegisterCmd(CmdTypeString,"get", get)
	RegisterCmd(CmdTypeString,"getrange", getrange)
	RegisterCmd(CmdTypeString,"getset", getset)
	RegisterCmd(CmdTypeString,"incr", incr)
	RegisterCmd(CmdTypeString,"incrby", incrby)
	RegisterCmd(CmdTypeString,"mget", mget)
	RegisterCmd(CmdTypeString,"mset", mset)
	RegisterCmd(CmdTypeString,"set", set)
	RegisterCmd(CmdTypeString,"setnx", setnx)
	RegisterCmd(CmdTypeString,"setex", setex)
	RegisterCmd(CmdTypeString,"setnxex", setnxex)
	RegisterCmd(CmdTypeString,"setxxex", setxxex)
	RegisterCmd(CmdTypeString,"setrange", setrange)
	RegisterCmd(CmdTypeString,"strlen", strlen)

	// just for string type key
	RegisterCmd(CmdTypeString,"del", del)
	RegisterCmd(CmdTypeString,"exists", exists)
	RegisterCmd(CmdTypeString,"expire", expire)
	RegisterCmd(CmdTypeString,"expireat", expireat)
	RegisterCmd(CmdTypeString,"persist", persist)
	RegisterCmd(CmdTypeString,"ttl", ttl)
}

func get(ctx context.Context, c *ConnClient, cmdParams [][]byte) (res interface{}, err error) {
	if len(cmdParams) != 1 {
		err = ErrCmdParams
		return
	}

	v, err := c.db.DBString().GetSlice(ctx, cmdParams[0])
	if err != nil {
		return
	}
	if v == nil {
		return
	}

	res = v.Data()
	v.Free()
	return
}

func set(ctx context.Context, c *ConnClient, cmdParams [][]byte) (res interface{}, err error) {
	if len(cmdParams) != 2 {
		err = ErrCmdParams
		return
	}

	if err = c.db.DBString().Set(ctx, cmdParams[0], cmdParams[1]); err != nil {
		return
	}

	return OK, nil
}

func appendCmd(ctx context.Context, c *ConnClient, cmdParams [][]byte) (res interface{}, err error) {
	if len(cmdParams) != 2 {
		err = ErrCmdParams
		return
	}

	res, err = c.db.DBString().Append(ctx, cmdParams[0], cmdParams[1])
	return
}

func decr(ctx context.Context, c *ConnClient, cmdParams [][]byte) (res interface{}, err error) {
	if len(cmdParams) != 1 {
		err = ErrCmdParams
		return
	}

	res, err = c.db.DBString().Decr(ctx, cmdParams[0])
	return
}

func decrby(ctx context.Context, c *ConnClient, cmdParams [][]byte) (res interface{}, err error) {
	if len(cmdParams) != 2 {
		err = ErrCmdParams
		return
	}
	delta, err := utils.StrInt64(cmdParams[1], nil)
	if err != nil {
		err = ErrValue
		return
	}

	res, err = c.db.DBString().DecrBy(ctx, cmdParams[0], delta)
	return
}

func del(ctx context.Context, c *ConnClient, cmdParams [][]byte) (res interface{}, err error) {
	if len(cmdParams) == 0 {
		err = ErrCmdParams
		return
	}

	res, err = c.db.DBString().Del(ctx, cmdParams...)
	return
}

func exists(ctx context.Context, c *ConnClient, cmdParams [][]byte) (res interface{}, err error) {
	if len(cmdParams) != 1 {
		err = ErrCmdParams
		return
	}

	res, err = c.db.DBString().Exists(ctx, cmdParams[0])
	return
}

func getrange(ctx context.Context, c *ConnClient, cmdParams [][]byte) (res interface{}, err error) {
	if len(cmdParams) != 3 {
		err = ErrCmdParams
		return
	}

	start, err := strconv.Atoi(string(cmdParams[1]))
	if err != nil {
		err = ErrValue
		return
	}

	end, err := strconv.Atoi(string(cmdParams[2]))
	if err != nil {
		err = ErrValue
		return
	}

	res, err = c.db.DBString().GetRange(ctx, cmdParams[0], start, end)
	return
}

func getset(ctx context.Context, c *ConnClient, cmdParams [][]byte) (res interface{}, err error) {
	if len(cmdParams) != 2 {
		err = ErrCmdParams
		return
	}

	res, err = c.db.DBString().GetSet(ctx, cmdParams[0], cmdParams[1])
	return
}

func incr(ctx context.Context, c *ConnClient, cmdParams [][]byte) (res interface{}, err error) {
	if len(cmdParams) != 1 {
		err = ErrCmdParams
		return
	}

	res, err = c.db.DBString().Incr(ctx, cmdParams[0])
	return
}

func incrby(ctx context.Context, c *ConnClient, cmdParams [][]byte) (res interface{}, err error) {
	if len(cmdParams) != 2 {
		err = ErrCmdParams
		return
	}

	delta, err := utils.StrInt64(cmdParams[1], nil)
	if err != nil {
		err = ErrValue
		return
	}

	res, err = c.db.DBString().IncrBy(ctx, cmdParams[0], delta)
	return
}

func mget(ctx context.Context, c *ConnClient, cmdParams [][]byte) (res interface{}, err error) {
	if len(cmdParams) == 0 {
		err = ErrCmdParams
		return
	}

	return c.db.DBString().MGet(ctx, cmdParams...)
}

func mset(ctx context.Context, c *ConnClient, cmdParams [][]byte) (res interface{}, err error) {
	if len(cmdParams) == 0 || len(cmdParams)%2 != 0 {
		err = ErrCmdParams
		return
	}

	kvs := make([]driver.KVPair, len(cmdParams)/2)
	for i := 0; i < len(kvs); i++ {
		kvs[i].Key = cmdParams[2*i]
		kvs[i].Value = cmdParams[2*i+1]
	}

	err = c.db.DBString().MSet(ctx, kvs...)
	if err != nil {
		return
	}
	res = OK

	return
}

func setnx(ctx context.Context, c *ConnClient, cmdParams [][]byte) (res interface{}, err error) {
	if len(cmdParams) != 2 {
		err = ErrCmdParams
		return
	}

	res, err = c.db.DBString().SetNX(ctx, cmdParams[0], cmdParams[1])
	return
}

func setex(ctx context.Context, c *ConnClient, cmdParams [][]byte) (res interface{}, err error) {
	if len(cmdParams) != 3 {
		err = ErrCmdParams
		return
	}

	sec, err := utils.StrInt64(cmdParams[1], nil)
	if err != nil {
		err = ErrValue
		return
	}

	err = c.db.DBString().SetEX(ctx, cmdParams[0], sec, cmdParams[2])
	if err != nil {
		return
	}

	res = OK
	return
}

func setnxex(ctx context.Context, c *ConnClient, cmdParams [][]byte) (res interface{}, err error) {
	if len(cmdParams) != 3 {
		err = ErrCmdParams
		return
	}

	sec, err := utils.StrInt64(cmdParams[1], nil)
	if err != nil {
		err = ErrValue
		return
	}

	err = c.db.DBString().SetNXEX(ctx, cmdParams[0], sec, cmdParams[2])
	if err != nil {
		return
	}

	res = OK
	return
}

func setxxex(ctx context.Context, c *ConnClient, cmdParams [][]byte) (res interface{}, err error) {
	if len(cmdParams) != 3 {
		err = ErrCmdParams
		return
	}

	sec, err := utils.StrInt64(cmdParams[1], nil)
	if err != nil {
		err = ErrValue
		return
	}

	err = c.db.DBString().SetXXEX(ctx, cmdParams[0], sec, cmdParams[2])
	if err != nil {
		return
	}

	res = OK
	return
}

func setrange(ctx context.Context, c *ConnClient, cmdParams [][]byte) (res interface{}, err error) {
	if len(cmdParams) != 3 {
		err = ErrCmdParams
		return
	}

	offset, err := strconv.Atoi(string(cmdParams[1]))
	if err != nil {
		err = ErrValue
		return
	}

	res, err = c.db.DBString().SetRange(ctx, cmdParams[0], offset, cmdParams[2])
	return
}

func strlen(ctx context.Context, c *ConnClient, cmdParams [][]byte) (res interface{}, err error) {
	if len(cmdParams) != 1 {
		err = ErrCmdParams
		return
	}

	res, err = c.db.DBString().StrLen(ctx, cmdParams[0])
	return
}

func expire(ctx context.Context, c *ConnClient, cmdParams [][]byte) (res interface{}, err error) {
	if len(cmdParams) != 2 {
		err = ErrCmdParams
		return
	}

	duration, err := utils.StrInt64(cmdParams[1], nil)
	if err != nil {
		err = ErrValue
		return
	}

	res, err = c.db.DBString().Expire(ctx, cmdParams[0], duration)
	return
}

func expireat(ctx context.Context, c *ConnClient, cmdParams [][]byte) (res interface{}, err error) {
	if len(cmdParams) != 2 {
		err = ErrCmdParams
		return
	}

	when, err := utils.StrInt64(cmdParams[1], nil)
	if err != nil {
		err = ErrValue
		return
	}

	res, err = c.db.DBString().ExpireAt(ctx, cmdParams[0], when)
	return
}

func ttl(ctx context.Context, c *ConnClient, cmdParams [][]byte) (res interface{}, err error) {
	if len(cmdParams) != 1 {
		err = ErrCmdParams
		return
	}

	res, err = c.db.DBString().TTL(ctx, cmdParams[0])
	return
}

func persist(ctx context.Context, c *ConnClient, cmdParams [][]byte) (res interface{}, err error) {
	if len(cmdParams) != 1 {
		err = ErrCmdParams
		return
	}

	res, err = c.db.DBString().Persist(ctx, cmdParams[0])
	return
}
