package standalone

import (
	"github.com/tidwall/redcon"
	"github.com/weedge/pkg/driver"
)

type RespCmdConn struct {
	*driver.RespConnBase

	srv      *RespCmdService
	isAuthed bool
	redConn  redcon.Conn
}

func (c *RespCmdConn) SetRedConn(redConn redcon.Conn) {
	c.redConn = redConn
}

func (c *RespCmdConn) GetRemoteAddr() string {
	return c.redConn.RemoteAddr()
}
