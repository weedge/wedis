package srv

import (
	"context"
	"strconv"
	"strings"

	"github.com/weedge/pkg/driver"
	"github.com/weedge/pkg/utils"
	"github.com/weedge/wedis/internal/srv/config"
	standalone "github.com/weedge/xdis-standalone"
)

var (
	ConfigOpts *config.Options
)

func init() {
	driver.RegisterCmd(driver.CmdTypeSrv, "config", configCmd)
}

func configCmd(ctx context.Context, c driver.IRespConn, cmdParams [][]byte) (res interface{}, err error) {
	if len(cmdParams) < 1 {
		return nil, standalone.ErrCmdParams
	}

	switch strings.ToLower(utils.Bytes2String(cmdParams[0])) {
	case "get":
		return getCfgCmd(ctx, c, cmdParams)
	default:
		return nil, standalone.ErrCmdParams
	}
}

func getCfgCmd(ctx context.Context, c driver.IRespConn, cmdParams [][]byte) (res interface{}, err error) {
	if len(cmdParams) != 2 {
		return nil, standalone.ErrCmdParams
	}

	data := []any{}
	switch strings.ToLower(utils.Bytes2String(cmdParams[1])) {
	case "databases":
		data = append(data,
			"databases",
			strconv.AppendInt(nil, int64(ConfigOpts.StoreCfg.Databases), 10),
		)
	case "maxmemory":
		data = append(data, "maxmemory", 0)
	}
	res = data

	return
}
