package srv

import (
	"context"

	"github.com/tidwall/redcon"
	"github.com/weedge/pkg/driver"
	"github.com/weedge/pkg/utils"
)

func init() {
	RegisterCmd(CmdTypeHash, "hexists", hexists)
	RegisterCmd(CmdTypeHash, "hget", hget)
	RegisterCmd(CmdTypeHash, "hgetall", hgetall)
	RegisterCmd(CmdTypeHash, "hincrby", hincrby)
	RegisterCmd(CmdTypeHash, "hkeys", hkeys)
	RegisterCmd(CmdTypeHash, "hlen", hlen)
	RegisterCmd(CmdTypeHash, "hmget", hmget)
	RegisterCmd(CmdTypeHash, "hmset", hmset)
	RegisterCmd(CmdTypeHash, "hset", hset)
	RegisterCmd(CmdTypeHash, "hvals", hvals)

	//del for hash
	RegisterCmd(CmdTypeHash, "hmclear", hmclear)
	//exists for hash
	RegisterCmd(CmdTypeHash, "hkeyexists", hkeyexists)
	//expire for hash
	RegisterCmd(CmdTypeHash, "hexpire", hexpire)
	//expireat for hash
	RegisterCmd(CmdTypeHash, "hexpireat", hexpireat)
	//ttl for hash
	RegisterCmd(CmdTypeHash, "httl", httl)
	//persist for hash
	RegisterCmd(CmdTypeHash, "hpersist", hpersist)
}

func hexists(ctx context.Context, c *ConnClient, cmdParams [][]byte) (res interface{}, err error) {
	if len(cmdParams) < 2 {
		err = ErrCmdParams
		return
	}

	v, err := c.db.DBHash().HGet(ctx, cmdParams[0], cmdParams[1])
	if err != nil {
		return
	}
	if v == nil {
		res = redcon.SimpleInt(0)
		return
	}

	res = redcon.SimpleInt(1)
	return
}

func hget(ctx context.Context, c *ConnClient, cmdParams [][]byte) (res interface{}, err error) {
	if len(cmdParams) < 2 {
		err = ErrCmdParams
		return
	}

	v, err := c.db.DBHash().HGet(ctx, cmdParams[0], cmdParams[1])
	if len(v) == 0 {
		return nil, nil
	}
	res = v
	return
}

func hgetall(ctx context.Context, c *ConnClient, cmdParams [][]byte) (res interface{}, err error) {
	if len(cmdParams) != 1 {
		err = ErrCmdParams
		return
	}

	data, err := c.db.DBHash().HGetAll(ctx, cmdParams[0])
	if err != nil {
		return
	}

	tmp := [][]byte{}
	for _, item := range data {
		tmp = append(tmp, item.Field)
		tmp = append(tmp, item.Value)
	}
	res = tmp

	return
}

func hincrby(ctx context.Context, c *ConnClient, cmdParams [][]byte) (res interface{}, err error) {
	if len(cmdParams) != 3 {
		err = ErrCmdParams
		return
	}
	delta, err := utils.StrInt64(cmdParams[2], nil)
	if err != nil {
		err = ErrValue
		return
	}

	res, err = c.db.DBHash().HIncrBy(ctx, cmdParams[0], cmdParams[1], delta)
	return
}

func hkeys(ctx context.Context, c *ConnClient, cmdParams [][]byte) (res interface{}, err error) {
	if len(cmdParams) != 1 {
		err = ErrCmdParams
		return
	}

	res, err = c.db.DBHash().HKeys(ctx, cmdParams[0])
	return
}

func hlen(ctx context.Context, c *ConnClient, cmdParams [][]byte) (res interface{}, err error) {
	if len(cmdParams) != 1 {
		err = ErrCmdParams
		return
	}

	res, err = c.db.DBHash().HLen(ctx, cmdParams[0])
	return
}

func hmget(ctx context.Context, c *ConnClient, cmdParams [][]byte) (res interface{}, err error) {
	if len(cmdParams) < 2 {
		err = ErrCmdParams
		return
	}

	res, err = c.db.DBHash().HMget(ctx, cmdParams[0], cmdParams[1:]...)
	return
}

func hmset(ctx context.Context, c *ConnClient, cmdParams [][]byte) (res interface{}, err error) {
	if len(cmdParams) < 3 || len(cmdParams[1:])%2 != 0 {
		err = ErrCmdParams
		return
	}

	args := cmdParams[1:]
	kvs := make([]driver.FVPair, len(args)/2)
	for i := 0; i < len(kvs); i++ {
		kvs[i].Field = args[2*i]
		kvs[i].Value = args[2*i+1]
	}

	if err = c.db.DBHash().HMset(ctx, cmdParams[0], kvs...); err != nil {
		return
	}

	res = OK
	return
}

func hset(ctx context.Context, c *ConnClient, cmdParams [][]byte) (res interface{}, err error) {
	if len(cmdParams) != 3 {
		err = ErrCmdParams
		return
	}

	res, err = c.db.DBHash().HSet(ctx, cmdParams[0], cmdParams[1], cmdParams[2])
	return
}

func hvals(ctx context.Context, c *ConnClient, cmdParams [][]byte) (res interface{}, err error) {
	if len(cmdParams) != 1 {
		err = ErrCmdParams
		return
	}

	res, err = c.db.DBHash().HValues(ctx, cmdParams[0])
	return
}

func hmclear(ctx context.Context, c *ConnClient, cmdParams [][]byte) (res interface{}, err error) {
	if len(cmdParams) == 0 {
		err = ErrCmdParams
		return
	}

	res, err = c.db.DBHash().Del(ctx, cmdParams...)
	return
}

func hkeyexists(ctx context.Context, c *ConnClient, cmdParams [][]byte) (res interface{}, err error) {
	if len(cmdParams) != 1 {
		err = ErrCmdParams
		return
	}

	res, err = c.db.DBHash().Exists(ctx, cmdParams[0])
	return
}

func hexpire(ctx context.Context, c *ConnClient, cmdParams [][]byte) (res interface{}, err error) {
	if len(cmdParams) != 2 {
		err = ErrCmdParams
		return
	}

	d, err := utils.StrInt64(cmdParams[1], nil)
	if err != nil {
		err = ErrValue
		return
	}

	res, err = c.db.DBHash().Expire(ctx, cmdParams[0], d)
	return
}

func hexpireat(ctx context.Context, c *ConnClient, cmdParams [][]byte) (res interface{}, err error) {
	if len(cmdParams) != 2 {
		err = ErrCmdParams
		return
	}

	d, err := utils.StrInt64(cmdParams[1], nil)
	if err != nil {
		err = ErrValue
		return
	}

	res, err = c.db.DBHash().ExpireAt(ctx, cmdParams[0], d)
	return
}

func httl(ctx context.Context, c *ConnClient, cmdParams [][]byte) (res interface{}, err error) {
	if len(cmdParams) != 1 {
		err = ErrCmdParams
		return
	}

	res, err = c.db.DBHash().TTL(ctx, cmdParams[0])
	return
}

func hpersist(ctx context.Context, c *ConnClient, cmdParams [][]byte) (res interface{}, err error) {
	if len(cmdParams) != 1 {
		err = ErrCmdParams
		return
	}

	res, err = c.db.DBHash().Persist(ctx, cmdParams[0])
	return
}
