package srv

import (
	"context"
	"strconv"

	"github.com/weedge/pkg/utils"
)

func init() {
	RegisterCmd("bitcount", bitcount)
	RegisterCmd("bitop", bitop)
	RegisterCmd("bitpos", bitpos)
	RegisterCmd("getbit", getbit)
	RegisterCmd("setbit", setbit)
}

func bitcount(ctx context.Context, c *ConnClient, cmdParams [][]byte) (res interface{}, err error) {
	if len(cmdParams) == 0 {
		return nil, ErrCmdParams
	}
	if len(cmdParams) > 3 {
		return nil, ErrCmdParams
	}

	start, end, err := parseBitRange(cmdParams[1:])
	if err != nil {
		return
	}

	res, err = c.db.DBBitmap().BitCount(cmdParams[0], start, end)
	return
}

func parseBitRange(args [][]byte) (start int, end int, err error) {
	start = 0
	end = -1
	if len(args) > 0 {
		if start, err = strconv.Atoi(string(args[0])); err != nil {
			return
		}
	}

	if len(args) == 2 {
		if end, err = strconv.Atoi(string(args[1])); err != nil {
			return
		}
	}

	return
}

func bitop(ctx context.Context, c *ConnClient, cmdParams [][]byte) (res interface{}, err error) {
	if len(cmdParams) < 3 {
		return nil, ErrCmdParams
	}

	op := utils.Bytes2String(cmdParams[0])
	destKey := cmdParams[1]
	srcKeys := cmdParams[2:]

	res, err = c.db.DBBitmap().BitOP(op, destKey, srcKeys...)

	return
}

func bitpos(ctx context.Context, c *ConnClient, cmdParams [][]byte) (res interface{}, err error) {
	if len(cmdParams) < 2 {
		return nil, ErrCmdParams
	}

	bit, err := strconv.Atoi(utils.Bytes2String(cmdParams[1]))
	if err != nil {
		return
	}

	start, end, err := parseBitRange(cmdParams[2:])
	if err != nil {
		return
	}

	res, err = c.db.DBBitmap().BitPos(cmdParams[0], bit, start, end)
	return
}

func getbit(ctx context.Context, c *ConnClient, cmdParams [][]byte) (res interface{}, err error) {
	if len(cmdParams) != 2 {
		return nil, ErrCmdParams
	}

	offset, err := strconv.Atoi(utils.Bytes2String(cmdParams[1]))
	if err != nil {
		return
	}

	res, err = c.db.DBBitmap().GetBit(cmdParams[0], offset)
	return
}

func setbit(ctx context.Context, c *ConnClient, cmdParams [][]byte) (res interface{}, err error) {
	if len(cmdParams) != 3 {
		return nil, ErrCmdParams
	}

	offset, err := strconv.Atoi(utils.Bytes2String(cmdParams[1]))
	if err != nil {
		return
	}

	value, err := strconv.Atoi(utils.Bytes2String(cmdParams[2]))
	if err != nil {
		return
	}

	res, err = c.db.DBBitmap().SetBit(cmdParams[0], offset, value)
	return
}
