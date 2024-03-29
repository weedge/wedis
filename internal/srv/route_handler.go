package srv

import (
	"bytes"
	"context"
	"net/http"
	"strconv"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/weedge/pkg/driver"
)

func (s *Server) SetupRoutes(h *server.Hertz) {
	SetupProbeRoutes(h)
	SetupCmdRoutes(h, s)
}

func SetupProbeRoutes(h *server.Hertz) {
	// ready
	h.GET("/readiness", func(ctx context.Context, c *app.RequestContext) {
		c.JSON(http.StatusOK, "server is readiness")
	})

	// liveness
	h.GET("/health", func(ctx context.Context, c *app.RequestContext) {
		c.JSON(http.StatusOK, "server is running")
	})
}

func SetupCmdRoutes(h *server.Hertz, s *Server) {
	h.GET("/cmds", func(ctx context.Context, c *app.RequestContext) {
		c.JSON(http.StatusOK, driver.RegisteredCmdSet)
	})

	// content-type : multipart/form-data
	h.POST("/:db/:cmd", func(ctx context.Context, c *app.RequestContext) {
		cmd := c.Param("cmd")
		cmdParams := c.FormValue("params")
		params := bytes.Split(cmdParams, []byte(" "))
		db := c.Param("db")
		dbIdx, err := strconv.Atoi(db)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusOK, err.Error())
			return
		}

		if cmd == "select" {
			c.AbortWithStatusJSON(http.StatusOK, "select unsupport for http api")
			return
		}

		cli := s.respSrv.InitRespConn(ctx, dbIdx)
		res, err := cli.DoCmd(ctx, cmd, params)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusOK, err.Error())
			return
		}

		c.JSON(http.StatusOK, res)
	})

}
