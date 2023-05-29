package srv

import (
	"context"

	"github.com/weedge/pkg/utils"
)

func init() {
	RegisterCmd("sadd", sadd)
	RegisterCmd("scard", scard)
	RegisterCmd("sdiff", sdiff)
	RegisterCmd("sdiffstore", sdiffstore)
	RegisterCmd("sinter", sinter)
	RegisterCmd("sinterstore", sinterstore)
	RegisterCmd("sismember", sismember)
	RegisterCmd("smembers", smembers)
	RegisterCmd("srem", srem)
	RegisterCmd("sunion", sunion)
	RegisterCmd("sunionstore", sunionstore)

	// del
	RegisterCmd("smclear", smclear)
	// expire
	RegisterCmd("sexpire", sexpire)
	// expireat
	RegisterCmd("sexpireat", sexpireat)
	// ttl
	RegisterCmd("sttl", sttl)
	// persist
	RegisterCmd("spersist", spersist)
	// exists
	RegisterCmd("skeyexists", skeyexists)
}

func sadd(ctx context.Context, c *ConnClient, cmdParams [][]byte) (res interface{}, err error) {
	if len(cmdParams) < 2 {
		err = ErrCmdParams
		return
	}

	res, err = c.db.DBSet().SAdd(cmdParams[0], cmdParams[1:]...)
	return
}

func scard(ctx context.Context, c *ConnClient, cmdParams [][]byte) (res interface{}, err error) {
	if len(cmdParams) != 1 {
		err = ErrCmdParams
		return
	}

	res, err = c.db.DBSet().SCard(cmdParams[0])
	return
}

func sdiff(ctx context.Context, c *ConnClient, cmdParams [][]byte) (res interface{}, err error) {
	if len(cmdParams) == 0 {
		err = ErrCmdParams
		return
	}

	res, err = c.db.DBSet().SDiff(cmdParams...)
	return
}

func sdiffstore(ctx context.Context, c *ConnClient, cmdParams [][]byte) (res interface{}, err error) {
	if len(cmdParams) < 2 {
		err = ErrCmdParams
		return
	}

	res, err = c.db.DBSet().SDiffStore(cmdParams[0], cmdParams[1:]...)
	return
}

func sinter(ctx context.Context, c *ConnClient, cmdParams [][]byte) (res interface{}, err error) {
	if len(cmdParams) == 0 {
		err = ErrCmdParams
		return
	}

	res, err = c.db.DBSet().SInter(cmdParams...)
	return
}

func sinterstore(ctx context.Context, c *ConnClient, cmdParams [][]byte) (res interface{}, err error) {
	if len(cmdParams) < 2 {
		err = ErrCmdParams
		return
	}

	res, err = c.db.DBSet().SInterStore(cmdParams[0], cmdParams[1:]...)
	return
}

func sismember(ctx context.Context, c *ConnClient, cmdParams [][]byte) (res interface{}, err error) {
	if len(cmdParams) != 2 {
		err = ErrCmdParams
		return
	}

	res, err = c.db.DBSet().SIsMember(cmdParams[0], cmdParams[1])
	return
}

func smembers(ctx context.Context, c *ConnClient, cmdParams [][]byte) (res interface{}, err error) {
	if len(cmdParams) != 1 {
		err = ErrCmdParams
		return
	}

	res, err = c.db.DBSet().SMembers(cmdParams[0])
	return
}

func srem(ctx context.Context, c *ConnClient, cmdParams [][]byte) (res interface{}, err error) {
	if len(cmdParams) < 2 {
		err = ErrCmdParams
		return
	}

	res, err = c.db.DBSet().SRem(cmdParams[0], cmdParams[1:]...)
	return
}

func sunion(ctx context.Context, c *ConnClient, cmdParams [][]byte) (res interface{}, err error) {
	if len(cmdParams) == 0 {
		err = ErrCmdParams
		return
	}

	res, err = c.db.DBSet().SUnion(cmdParams...)
	return
}

func sunionstore(ctx context.Context, c *ConnClient, cmdParams [][]byte) (res interface{}, err error) {
	if len(cmdParams) < 2 {
		err = ErrCmdParams
		return
	}

	res, err = c.db.DBSet().SUnionStore(cmdParams[0], cmdParams[1:]...)
	return
}

func smclear(ctx context.Context, c *ConnClient, cmdParams [][]byte) (res interface{}, err error) {
	if len(cmdParams) == 0 {
		err = ErrCmdParams
		return
	}

	res, err = c.db.DBSet().Del(cmdParams...)
	return
}

func skeyexists(ctx context.Context, c *ConnClient, cmdParams [][]byte) (res interface{}, err error) {
	if len(cmdParams) != 1 {
		err = ErrCmdParams
		return
	}

	res, err = c.db.DBSet().Exists(cmdParams[0])
	return
}

func sexpire(ctx context.Context, c *ConnClient, cmdParams [][]byte) (res interface{}, err error) {
	if len(cmdParams) != 2 {
		err = ErrCmdParams
		return
	}

	d, err := utils.StrInt64(cmdParams[2], nil)
	if err != nil {
		err = ErrValue
		return
	}

	res, err = c.db.DBSet().Expire(cmdParams[0], d)
	return
}

func sexpireat(ctx context.Context, c *ConnClient, cmdParams [][]byte) (res interface{}, err error) {
	if len(cmdParams) != 2 {
		err = ErrCmdParams
		return
	}

	d, err := utils.StrInt64(cmdParams[1], nil)
	if err != nil {
		err = ErrValue
		return
	}

	res, err = c.db.DBSet().ExpireAt(cmdParams[0], d)
	return
}

func sttl(ctx context.Context, c *ConnClient, cmdParams [][]byte) (res interface{}, err error) {
	if len(cmdParams) != 1 {
		err = ErrCmdParams
		return
	}

	res, err = c.db.DBSet().TTL(cmdParams[0])
	return
}

func spersist(ctx context.Context, c *ConnClient, cmdParams [][]byte) (res interface{}, err error) {
	if len(cmdParams) != 1 {
		err = ErrCmdParams
		return
	}

	res, err = c.db.DBSet().Persist(cmdParams[0])
	return
}
