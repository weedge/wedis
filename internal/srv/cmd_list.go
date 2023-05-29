package srv

import (
	"bytes"
	"context"
	"strconv"
	"time"

	"github.com/weedge/pkg/utils"
)

func init() {
	RegisterCmd("blpop", blpop)
	RegisterCmd("brpop", brpop)
	RegisterCmd("lindex", lindex)
	RegisterCmd("llen", llen)
	RegisterCmd("lpop", lpop)
	RegisterCmd("lrange", lrange)
	RegisterCmd("lset", lset)
	RegisterCmd("lpush", lpush)
	RegisterCmd("rpop", rpop)
	RegisterCmd("rpush", rpush)
	RegisterCmd("brpoplpush", brpoplpush)
	RegisterCmd("rpoplpush", rpoplpush)

	//del for list
	RegisterCmd("lmclear", lmclear)
	//exists for list
	RegisterCmd("lkeyexists", lkeyexists)
	//expire for list
	RegisterCmd("lexpire", lexpire)
	//expireat for list
	RegisterCmd("lexpireat", lexpireat)
	//ttl for list
	RegisterCmd("lttl", lttl)
	//persist for list
	RegisterCmd("lpersist", lpersist)
}

func lmclear(ctx context.Context, c *ConnClient, cmdParams [][]byte) (res interface{}, err error) {
	if len(cmdParams) == 0 {
		err = ErrCmdParams
		return
	}

	res, err = c.db.DBList().Del(cmdParams...)
	return
}

func blpop(ctx context.Context, c *ConnClient, cmdParams [][]byte) (res interface{}, err error) {
	if len(cmdParams) < 2 {
		err = ErrCmdParams
		return
	}

	t, err := strconv.ParseFloat(utils.Bytes2String(cmdParams[len(cmdParams)-1]), 64)
	if err != nil {
		return
	}
	timeout := time.Duration(t * float64(time.Second))

	res, err = c.db.DBList().BLPop(cmdParams[:len(cmdParams)-1], timeout)
	return
}

func brpop(ctx context.Context, c *ConnClient, cmdParams [][]byte) (res interface{}, err error) {
	if len(cmdParams) < 2 {
		err = ErrCmdParams
		return
	}

	t, err := strconv.ParseFloat(utils.Bytes2String(cmdParams[len(cmdParams)-1]), 64)
	if err != nil {
		return
	}
	timeout := time.Duration(t * float64(time.Second))

	res, err = c.db.DBList().BRPop(cmdParams[:len(cmdParams)-1], timeout)
	return
}

func lindex(ctx context.Context, c *ConnClient, cmdParams [][]byte) (res interface{}, err error) {
	if len(cmdParams) != 2 {
		err = ErrCmdParams
		return
	}

	i, err := utils.StrInt64(cmdParams[1], nil)
	if err != nil {
		err = ErrValue
		return
	}

	res, err = c.db.DBList().LIndex(cmdParams[0], int32(i))
	return
}

func llen(ctx context.Context, c *ConnClient, cmdParams [][]byte) (res interface{}, err error) {
	if len(cmdParams) != 0 {
		err = ErrCmdParams
		return
	}

	res, err = c.db.DBList().LLen(cmdParams[0])
	return
}

func lpop(ctx context.Context, c *ConnClient, cmdParams [][]byte) (res interface{}, err error) {
	if len(cmdParams) != 1 {
		err = ErrCmdParams
		return
	}

	res, err = c.db.DBList().LPop(cmdParams[0])
	return
}

func lrange(ctx context.Context, c *ConnClient, cmdParams [][]byte) (res interface{}, err error) {
	if len(cmdParams) != 3 {
		err = ErrCmdParams
		return
	}

	start, err := utils.StrInt64(cmdParams[1], nil)
	if err != nil {
		err = ErrValue
		return
	}
	end, err := utils.StrInt64(cmdParams[2], nil)
	if err != nil {
		err = ErrValue
		return
	}

	res, err = c.db.DBList().LRange(cmdParams[0], int32(start), int32(end))
	return
}

func lset(ctx context.Context, c *ConnClient, cmdParams [][]byte) (res interface{}, err error) {
	if len(cmdParams) != 3 {
		err = ErrCmdParams
		return
	}
	i, err := utils.StrInt64(cmdParams[1], nil)
	if err != nil {
		err = ErrValue
		return
	}

	if err = c.db.DBList().LSet(cmdParams[0], int32(i), cmdParams[2]); err != nil {
		return
	}

	res = OK
	return
}

func lpush(ctx context.Context, c *ConnClient, cmdParams [][]byte) (res interface{}, err error) {
	if len(cmdParams) < 2 {
		err = ErrCmdParams
		return
	}

	res, err = c.db.DBList().LPush(cmdParams[0], cmdParams[1:]...)
	return
}

func rpop(ctx context.Context, c *ConnClient, cmdParams [][]byte) (res interface{}, err error) {
	if len(cmdParams) != 1 {
		err = ErrCmdParams
		return
	}

	res, err = c.db.DBList().RPop(cmdParams[0])
	return
}

func rpush(ctx context.Context, c *ConnClient, cmdParams [][]byte) (res interface{}, err error) {
	if len(cmdParams) < 2 {
		err = ErrCmdParams
		return
	}

	res, err = c.db.DBList().RPush(cmdParams[0], cmdParams[1:]...)
	return
}

func brpoplpush(ctx context.Context, c *ConnClient, cmdParams [][]byte) (res interface{}, err error) {
	if len(cmdParams) != 3 {
		err = ErrCmdParams
		return
	}

	source, dest := cmdParams[0], cmdParams[1]
	t, err := strconv.ParseFloat(utils.Bytes2String(cmdParams[len(cmdParams)-1]), 64)
	if err != nil {
		return
	}
	timeout := time.Duration(t * float64(time.Second))

	ttl := int64(-1)
	// source dest equal, same list, get ttl
	if bytes.Equal(source, dest) {
		ttl, err = c.db.DBList().TTL(source)
		if err != nil {
			return
		}
	}

	kvdata, err := c.db.DBList().BRPop([][]byte{source}, timeout)
	if err != nil {
		return
	}
	if kvdata == nil {
		return
	}
	if len(kvdata) < 2 {
		return
	}

	vdata, ok := kvdata[1].([]byte)
	if !ok {
		err = ErrValue
		return
	}

	// lpush err rpush back
	if _, err = c.db.DBList().LPush(dest, vdata); err != nil {
		c.db.DBList().RPush(source, vdata)
		return
	}

	// reset tll
	if ttl != -1 {
		c.db.DBList().Expire(source, ttl)
	}

	res = vdata
	return
}

func rpoplpush(ctx context.Context, c *ConnClient, cmdParams [][]byte) (res interface{}, err error) {
	if len(cmdParams) != 2 {
		err = ErrCmdParams
		return
	}

	source, dest := cmdParams[0], cmdParams[1]
	ttl := int64(-1)
	// source dest equal, same list, get ttl
	if bytes.Equal(source, dest) {
		ttl, err = c.db.DBList().TTL(source)
		if err != nil {
			return
		}
	}

	data, err := c.db.DBList().RPop(source)
	if err != nil {
		return
	}
	if data == nil {
		return
	}

	// lpush err rpush back
	if _, err = c.db.DBList().LPush(dest, data); err != nil {
		c.db.DBList().RPush(source, data)
		return
	}

	// reset tll
	if ttl != -1 {
		c.db.DBList().Expire(source, ttl)
	}

	res = data
	return
}

func lkeyexists(ctx context.Context, c *ConnClient, cmdParams [][]byte) (res interface{}, err error) {
	if len(cmdParams) != 1 {
		err = ErrCmdParams
		return
	}

	res, err = c.db.DBList().Exists(cmdParams[0])
	return
}

func lexpire(ctx context.Context, c *ConnClient, cmdParams [][]byte) (res interface{}, err error) {
	if len(cmdParams) != 2 {
		err = ErrCmdParams
		return
	}

	d, err := utils.StrInt64(cmdParams[2], nil)
	if err != nil {
		err = ErrValue
		return
	}

	res, err = c.db.DBList().Expire(cmdParams[0], d)
	return
}

func lexpireat(ctx context.Context, c *ConnClient, cmdParams [][]byte) (res interface{}, err error) {
	if len(cmdParams) != 2 {
		err = ErrCmdParams
		return
	}

	d, err := utils.StrInt64(cmdParams[1], nil)
	if err != nil {
		err = ErrValue
		return
	}

	res, err = c.db.DBList().ExpireAt(cmdParams[0], d)
	return
}

func lttl(ctx context.Context, c *ConnClient, cmdParams [][]byte) (res interface{}, err error) {
	if len(cmdParams) != 1 {
		err = ErrCmdParams
		return
	}

	res, err = c.db.DBList().TTL(cmdParams[0])
	return
}

func lpersist(ctx context.Context, c *ConnClient, cmdParams [][]byte) (res interface{}, err error) {
	if len(cmdParams) != 1 {
		err = ErrCmdParams
		return
	}

	res, err = c.db.DBList().Persist(cmdParams[0])
	return
}
