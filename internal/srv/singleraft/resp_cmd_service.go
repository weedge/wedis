package singleraft

import (
	"context"

	"github.com/tidwall/redcon"
	"github.com/weedge/pkg/driver"
	"github.com/weedge/wedis/internal/srv/singleraft/config"
)

type RespCmdService struct {
	opts *config.RespCmdServiceOptions
	// redcon server handler
	mux *redcon.ServeMux
	// redcon server
	redconSrv *redcon.Server
	// storager
	store driver.IStorager
}

// Start service
func (s *RespCmdService) Start(ctx context.Context) (err error) {
	return
}

// InitRespConn init resp connect session by select db index,
// return IRespConn interface
func (s *RespCmdService) InitRespConn(ctx context.Context, dbIdx int) driver.IRespConn {
	if dbIdx < 0 {
		dbIdx = 0
	}

	conn := &RespCmdConn{RespConnBase: &driver.RespConnBase{}, srv: s, isAuthed: false}
	db, err := s.store.Select(ctx, dbIdx)
	if err != nil {
		return nil
	}
	conn.SetDb(db)
	return conn
}

// Close resp service
func (s *RespCmdService) Close() (err error) {
	return
}

// Name
func (s *RespCmdService) Name() driver.RespServiceName {
	return config.RegisterRespSrvModeName
}

// SetStorager
func (s *RespCmdService) SetStorager(store driver.IStorager) {
}
