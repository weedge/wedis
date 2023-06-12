package srv

import (
	"context"
	"errors"
	"strings"

	"github.com/weedge/pkg/driver"
)

type cmdHandle func(ctx context.Context, h *ConnClient, cmdParams [][]byte) (interface{}, error)

var RegisteredCommands = map[string]cmdHandle{}

func RegisterCmd(cmd string, handle cmdHandle) {
	RegisteredCommands[cmd] = handle
}

type ConnClient struct {
	srv      *Server
	db       driver.IDB
	isAuthed bool
	name     string
}

func (c *ConnClient) SetSrv(srv *Server) {
	c.srv = srv
}

func (c *ConnClient) SetDb(db driver.IDB) {
	c.db = db
}

func (c *ConnClient) SetConnName(name string) {
	c.name = name
}
func (c *ConnClient) Name() (name string) {
	return c.name
}

func (c *ConnClient) DoCmd(ctx context.Context, cmd string, cmdParams [][]byte) (res interface{}, err error) {
	cmd = strings.ToLower(strings.TrimSpace(cmd))
	f, ok := RegisteredCommands[cmd]
	if !ok {
		err = errors.New("ERR unknown command '" + cmd + "'")
		return
	}

	res, err = f(ctx, c, cmdParams)
	if err != nil {
		return
	}

	return
}
