package srv

import (
	"errors"

	"github.com/tidwall/redcon"
)

type CtxKey int

const (
	RespCmdCtxKey CtxKey = iota
)

var (
	ErrEmptyCommand          = errors.New("ERR empty command")
	ErrNotFound              = errors.New("ERR command not found")
	ErrNotAuthenticated      = errors.New("ERR not authenticated")
	ErrAuthenticationFailure = errors.New("ERR authentication failure")
	ErrCmdParams             = errors.New("ERR wrong number of arguments")
	ErrValue                 = errors.New("ERR value is not an integer or out of range")
	ErrSyntax                = errors.New("ERR syntax error")
	ErrOffset                = errors.New("ERR offset bit is not an natural number")
	ErrBool                  = errors.New("ERR value is not 0 or 1")
	ErrScoreOverflow         = errors.New("ERR zset score overflow")

	ErrProtocalVer  = errors.New("ERR Protocol version is not an integer or out of range")
	ErrUnsupportVer = errors.New("NOPROTO unsupported protocol version")
	ErrNOPwd        = errors.New("ERR Client sent AUTH, but no password is set")
	ErrInvalidPwd   = errors.New("ERR invalid password")
)

const (
	PONG  = redcon.SimpleString("PONG")
	OK    = redcon.SimpleString("OK")
	NOKEY = redcon.SimpleString("NOKEY")
)
