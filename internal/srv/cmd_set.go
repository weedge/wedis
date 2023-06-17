package srv

import (
	"context"

	"github.com/weedge/pkg/utils"
)

func init() {
	RegisterCmd(CmdTypeSet, "sadd", sadd)
	RegisterCmd(CmdTypeSet, "scard", scard)
	RegisterCmd(CmdTypeSet, "sdiff", sdiff)
	RegisterCmd(CmdTypeSet, "sdiffstore", sdiffstore)
	RegisterCmd(CmdTypeSet, "sinter", sinter)
	RegisterCmd(CmdTypeSet, "sinterstore", sinterstore)
	RegisterCmd(CmdTypeSet, "sismember", sismember)
	RegisterCmd(CmdTypeSet, "smembers", smembers)
	RegisterCmd(CmdTypeSet, "srem", srem)
	RegisterCmd(CmdTypeSet, "sunion", sunion)
	RegisterCmd(CmdTypeSet, "sunionstore", sunionstore)

	// del
	RegisterCmd(CmdTypeSet, "smclear", smclear)
	// expire
	RegisterCmd(CmdTypeSet, "sexpire", sexpire)
	// expireat
	RegisterCmd(CmdTypeSet, "sexpireat", sexpireat)
	// ttl
	RegisterCmd(CmdTypeSet, "sttl", sttl)
	// persist
	RegisterCmd(CmdTypeSet, "spersist", spersist)
	// exists
	RegisterCmd(CmdTypeSet, "skeyexists", skeyexists)
}

func sadd(ctx context.Context, c *ConnClient, cmdParams [][]byte) (res interface{}, err error) {
	if len(cmdParams) < 2 {
		err = ErrCmdParams
		return
	}

	res, err = c.db.DBSet().SAdd(ctx, cmdParams[0], cmdParams[1:]...)
	return
}

func scard(ctx context.Context, c *ConnClient, cmdParams [][]byte) (res interface{}, err error) {
	if len(cmdParams) != 1 {
		err = ErrCmdParams
		return
	}

	res, err = c.db.DBSet().SCard(ctx, cmdParams[0])
	return
}

func sdiff(ctx context.Context, c *ConnClient, cmdParams [][]byte) (res interface{}, err error) {
	if len(cmdParams) == 0 {
		err = ErrCmdParams
		return
	}

	res, err = c.db.DBSet().SDiff(ctx, cmdParams...)
	return
}

func sdiffstore(ctx context.Context, c *ConnClient, cmdParams [][]byte) (res interface{}, err error) {
	if len(cmdParams) < 2 {
		err = ErrCmdParams
		return
	}

	res, err = c.db.DBSet().SDiffStore(ctx, cmdParams[0], cmdParams[1:]...)
	return
}

func sinter(ctx context.Context, c *ConnClient, cmdParams [][]byte) (res interface{}, err error) {
	if len(cmdParams) == 0 {
		err = ErrCmdParams
		return
	}

	res, err = c.db.DBSet().SInter(ctx, cmdParams...)
	return
}

func sinterstore(ctx context.Context, c *ConnClient, cmdParams [][]byte) (res interface{}, err error) {
	if len(cmdParams) < 2 {
		err = ErrCmdParams
		return
	}

	res, err = c.db.DBSet().SInterStore(ctx, cmdParams[0], cmdParams[1:]...)
	return
}

func sismember(ctx context.Context, c *ConnClient, cmdParams [][]byte) (res interface{}, err error) {
	if len(cmdParams) != 2 {
		err = ErrCmdParams
		return
	}

	res, err = c.db.DBSet().SIsMember(ctx, cmdParams[0], cmdParams[1])
	return
}

func smembers(ctx context.Context, c *ConnClient, cmdParams [][]byte) (res interface{}, err error) {
	if len(cmdParams) != 1 {
		err = ErrCmdParams
		return
	}

	res, err = c.db.DBSet().SMembers(ctx, cmdParams[0])
	return
}

func srem(ctx context.Context, c *ConnClient, cmdParams [][]byte) (res interface{}, err error) {
	if len(cmdParams) < 2 {
		err = ErrCmdParams
		return
	}

	res, err = c.db.DBSet().SRem(ctx, cmdParams[0], cmdParams[1:]...)
	return
}

func sunion(ctx context.Context, c *ConnClient, cmdParams [][]byte) (res interface{}, err error) {
	if len(cmdParams) == 0 {
		err = ErrCmdParams
		return
	}

	res, err = c.db.DBSet().SUnion(ctx, cmdParams...)
	return
}

func sunionstore(ctx context.Context, c *ConnClient, cmdParams [][]byte) (res interface{}, err error) {
	if len(cmdParams) < 2 {
		err = ErrCmdParams
		return
	}

	res, err = c.db.DBSet().SUnionStore(ctx, cmdParams[0], cmdParams[1:]...)
	return
}

func smclear(ctx context.Context, c *ConnClient, cmdParams [][]byte) (res interface{}, err error) {
	if len(cmdParams) == 0 {
		err = ErrCmdParams
		return
	}

	res, err = c.db.DBSet().Del(ctx, cmdParams...)
	return
}

func skeyexists(ctx context.Context, c *ConnClient, cmdParams [][]byte) (res interface{}, err error) {
	if len(cmdParams) != 1 {
		err = ErrCmdParams
		return
	}

	res, err = c.db.DBSet().Exists(ctx, cmdParams[0])
	return
}

func sexpire(ctx context.Context, c *ConnClient, cmdParams [][]byte) (res interface{}, err error) {
	if len(cmdParams) != 2 {
		err = ErrCmdParams
		return
	}

	d, err := utils.StrInt64(cmdParams[1], nil)
	if err != nil {
		err = ErrValue
		return
	}

	res, err = c.db.DBSet().Expire(ctx, cmdParams[0], d)
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

	res, err = c.db.DBSet().ExpireAt(ctx, cmdParams[0], d)
	return
}

func sttl(ctx context.Context, c *ConnClient, cmdParams [][]byte) (res interface{}, err error) {
	if len(cmdParams) != 1 {
		err = ErrCmdParams
		return
	}

	res, err = c.db.DBSet().TTL(ctx, cmdParams[0])
	return
}

func spersist(ctx context.Context, c *ConnClient, cmdParams [][]byte) (res interface{}, err error) {
	if len(cmdParams) != 1 {
		err = ErrCmdParams
		return
	}

	res, err = c.db.DBSet().Persist(ctx, cmdParams[0])
	return
}
