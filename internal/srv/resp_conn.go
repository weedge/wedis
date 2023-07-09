package srv

import (
	"github.com/weedge/pkg/driver"
)

type AuthRespSrvConn struct {
	*driver.RespConnBase

	srv      *Server
	isAuthed bool
}
