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
