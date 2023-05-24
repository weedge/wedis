package srv

import (
	"context"
	"errors"

	"github.com/cloudwego/kitex/pkg/klog"
	"github.com/weedge/wedis/internal/srv/storager"
)

type cmdHandle func(h *ConnClient, params [][]byte) (interface{}, error)

var SupportedCommands = map[string]cmdHandle{
	"ping": Ping,
	"get":  Get,
}

type ConnClient struct {
	srv *Server
	db  *storager.DB
}

func (c *ConnClient) SetSrv(srv *Server) {
	c.srv = srv
}

func (c *ConnClient) SetDb(db *storager.DB) {
	c.db = db
}

func (c *ConnClient) DoCmd(ctx context.Context, cmd string, cmdParams [][]byte) (res interface{}, err error) {
	klog.Debugf("cmd:%s cmdParams:%s len:%d", cmd, cmdParams, len(cmdParams))
	f, ok := SupportedCommands[cmd]
	if !ok {
		err = errors.New("ERR unknown command '" + cmd + "'")
		return
	}

	res, err = f(c, cmdParams)
	if err != nil {
		return
	}

	return
}
