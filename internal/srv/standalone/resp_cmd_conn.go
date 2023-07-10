package standalone

import (
	"github.com/weedge/pkg/driver"
)

type AuthRespSrvConn struct {
	*driver.RespConnBase

	srv      *RespCmdService
	isAuthed bool
}
