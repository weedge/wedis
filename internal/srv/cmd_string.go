package srv

import "context"

func init() {
	RegisterCmd("get", get)
	RegisterCmd("set", set)
}

func get(ctx context.Context, c *ConnClient, cmdParams [][]byte) (res interface{}, err error) {
	if len(cmdParams) != 1 {
		err = ErrCmdParams
		return
	}

	v, err := c.db.DBString().GetSlice(cmdParams[0])
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

	if err = c.db.DBString().Set(cmdParams[0], cmdParams[1]); err != nil {
		return
	}

	return OK, nil
}

func append(ctx context.Context, c *ConnClient, cmdParams [][]byte) (res interface{}, err error) {
	return
}

func decr(ctx context.Context, c *ConnClient, cmdParams [][]byte) (res interface{}, err error) {
	return
}

func decrby(ctx context.Context, c *ConnClient, cmdParams [][]byte) (res interface{}, err error) {
	return
}

func del(ctx context.Context, c *ConnClient, cmdParams [][]byte) (res interface{}, err error) {
	return
}

func exists(ctx context.Context, c *ConnClient, cmdParams [][]byte) (res interface{}, err error) {
	return
}

func getrange(ctx context.Context, c *ConnClient, cmdParams [][]byte) (res interface{}, err error) {
	return
}

func getset(ctx context.Context, c *ConnClient, cmdParams [][]byte) (res interface{}, err error) {
	return
}

func incr(ctx context.Context, c *ConnClient, cmdParams [][]byte) (res interface{}, err error) {
	return
}

func incrby(ctx context.Context, c *ConnClient, cmdParams [][]byte) (res interface{}, err error) {
	return
}

func mget(ctx context.Context, c *ConnClient, cmdParams [][]byte) (res interface{}, err error) {
	return
}

func mset(ctx context.Context, c *ConnClient, cmdParams [][]byte) (res interface{}, err error) {
	return
}

func setnx(ctx context.Context, c *ConnClient, cmdParams [][]byte) (res interface{}, err error) {
	return
}

func setex(ctx context.Context, c *ConnClient, cmdParams [][]byte) (res interface{}, err error) {
	return
}

func setrange(ctx context.Context, c *ConnClient, cmdParams [][]byte) (res interface{}, err error) {
	return
}

func strlen(ctx context.Context, c *ConnClient, cmdParams [][]byte) (res interface{}, err error) {
	return
}

func expire(ctx context.Context, c *ConnClient, cmdParams [][]byte) (res interface{}, err error) {
	return
}

func expireat(ctx context.Context, c *ConnClient, cmdParams [][]byte) (res interface{}, err error) {
	return
}

func ttl(ctx context.Context, c *ConnClient, cmdParams [][]byte) (res interface{}, err error) {
	return
}

func persist(ctx context.Context, c *ConnClient, cmdParams [][]byte) (res interface{}, err error) {
	return
}
