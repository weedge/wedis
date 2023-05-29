package srv

import (
	"context"

	"github.com/weedge/pkg/driver"
	"github.com/weedge/pkg/utils"
)

func init() {
	RegisterCmd("hdel", hdel)
	RegisterCmd("hexists", hexists)
	RegisterCmd("hget", hget)
	RegisterCmd("hgetall", hgetall)
	RegisterCmd("hincrby", hincrby)
	RegisterCmd("hkeys", hkeys)
	RegisterCmd("hlen", hlen)
	RegisterCmd("hmget", hmget)
	RegisterCmd("hmset", hmset)
	RegisterCmd("hset", hset)
	RegisterCmd("hvals", hvals)

	//del for hash
	RegisterCmd("hmclear", hmclear)
	//exists for hash
	RegisterCmd("hkeyexists", hkeyexists)
	//expire for hash
	RegisterCmd("hexpire", hexpire)
	//expireat for hash
	RegisterCmd("hexpireat", hexpireat)
	//ttl for hash
	RegisterCmd("httl", httl)
	//persist for hash
	RegisterCmd("hpersist", hpersist)
}

func hdel(ctx context.Context, c *ConnClient, cmdParams [][]byte) (res interface{}, err error) {
	if len(cmdParams) < 2 {
		err = ErrCmdParams
		return
	}

	res, err = c.db.DBHash().HDel(cmdParams[0], cmdParams[1:]...)
	return
}

func hexists(ctx context.Context, c *ConnClient, cmdParams [][]byte) (res interface{}, err error) {
	if len(cmdParams) < 2 {
		err = ErrCmdParams
		return
	}

	v, err := c.db.DBHash().HGet(cmdParams[0], cmdParams[1])
	if err != nil {
		return
	}
	if v == nil {
		res = 0
	}

	return
}

func hget(ctx context.Context, c *ConnClient, cmdParams [][]byte) (res interface{}, err error) {
	if len(cmdParams) < 2 {
		err = ErrCmdParams
		return
	}

	v, err := c.db.DBHash().HGet(cmdParams[0], cmdParams[1])
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

	res, err = c.db.DBHash().HGetAll(cmdParams[0])
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

	res, err = c.db.DBHash().HIncrBy(cmdParams[0], cmdParams[1], delta)
	return
}

func hkeys(ctx context.Context, c *ConnClient, cmdParams [][]byte) (res interface{}, err error) {
	if len(cmdParams) != 1 {
		err = ErrCmdParams
		return
	}

	res, err = c.db.DBHash().HKeys(cmdParams[0])
	return
}

func hlen(ctx context.Context, c *ConnClient, cmdParams [][]byte) (res interface{}, err error) {
	if len(cmdParams) != 1 {
		err = ErrCmdParams
		return
	}

	res, err = c.db.DBHash().HLen(cmdParams[0])
	return
}

func hmget(ctx context.Context, c *ConnClient, cmdParams [][]byte) (res interface{}, err error) {
	if len(cmdParams) < 2 {
		err = ErrCmdParams
		return
	}

	res, err = c.db.DBHash().HMget(cmdParams[0], cmdParams[1:]...)
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

	if err = c.db.DBHash().HMset(cmdParams[0], kvs...); err != nil {
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

	res, err = c.db.DBHash().HSet(cmdParams[0], cmdParams[1], cmdParams[2])
	return
}

func hvals(ctx context.Context, c *ConnClient, cmdParams [][]byte) (res interface{}, err error) {
	if len(cmdParams) != 1 {
		err = ErrCmdParams
		return
	}

	res, err = c.db.DBHash().HValues(cmdParams[0])
	return
}

func hmclear(ctx context.Context, c *ConnClient, cmdParams [][]byte) (res interface{}, err error) {
	if len(cmdParams) == 0 {
		err = ErrCmdParams
		return
	}

	res, err = c.db.DBHash().Del(cmdParams...)
	return
}

func hkeyexists(ctx context.Context, c *ConnClient, cmdParams [][]byte) (res interface{}, err error) {
	if len(cmdParams) != 1 {
		err = ErrCmdParams
		return
	}

	res, err = c.db.DBHash().Exists(cmdParams[0])
	return
}

func hexpire(ctx context.Context, c *ConnClient, cmdParams [][]byte) (res interface{}, err error) {
	if len(cmdParams) != 2 {
		err = ErrCmdParams
		return
	}

	d, err := utils.StrInt64(cmdParams[2], nil)
	if err != nil {
		err = ErrValue
		return
	}

	res, err = c.db.DBHash().Expire(cmdParams[0], d)
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

	res, err = c.db.DBHash().ExpireAt(cmdParams[0], d)
	return
}

func httl(ctx context.Context, c *ConnClient, cmdParams [][]byte) (res interface{}, err error) {
	if len(cmdParams) != 1 {
		err = ErrCmdParams
		return
	}

	res, err = c.db.DBHash().TTL(cmdParams[0])
	return
}

func hpersist(ctx context.Context, c *ConnClient, cmdParams [][]byte) (res interface{}, err error) {
	if len(cmdParams) != 1 {
		err = ErrCmdParams
		return
	}

	res, err = c.db.DBHash().Persist(cmdParams[0])
	return
}
