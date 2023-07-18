package singleraft

import "github.com/weedge/pkg/driver"

type RespCmdConn struct {
	*driver.RespConnBase

	srv      *RespCmdService
	isAuthed bool
}
