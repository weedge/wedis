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
	ErrEmptyCommand          = errors.New("empty command")
	ErrNotFound              = errors.New("command not found")
	ErrNotAuthenticated      = errors.New("not authenticated")
	ErrAuthenticationFailure = errors.New("authentication failure")
	ErrCmdParams             = errors.New("invalid command param")
	ErrValue                 = errors.New("value is not an integer or out of range")
	ErrSyntax                = errors.New("syntax error")
	ErrOffset                = errors.New("offset bit is not an natural number")
	ErrBool                  = errors.New("value is not 0 or 1")
	ErrScoreOverflow         = errors.New("zset score overflow")
)

const (
	PONG  = redcon.SimpleString("PONG")
	OK    = redcon.SimpleString("OK")
	NOKEY = redcon.SimpleString("NOKEY")
)
