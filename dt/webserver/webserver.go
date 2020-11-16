package webserver

import (
	"distate-task/config"
	"distate-task/dt/db"
	"distate-task/dt/logger"
	"github.com/fasthttp/router"
	"github.com/oklog/run"
	"github.com/valyala/fasthttp"
	"net"
)

type WebServer struct {
	Config config.WebServer
	Log    *logger.FxLogger
	ln     net.Listener
	fs     *fasthttp.Server
	router *router.Router
	debug  bool
	conn   *db.Connection
}

func (s *WebServer) Run() (err error) {
	s.fs = &fasthttp.Server{
		Handler:            s.router.Handler,
		Name:               "dt-server",
		ReadBufferSize:     1024 * 2,
		MaxConnsPerIP:      5,
		MaxRequestsPerConn: 100,
		MaxRequestBodySize: 1024 * 4,
		Concurrency:        3000,
		Logger:             s.Log,
	}

	var g run.Group
	g.Add(func() error {
		return s.fs.Serve(s.ln)
	}, func(e error) {
		_ = s.ln.Close()
	})
	return g.Run()
}

func (s *WebServer) Close() (err error) {
	return s.ln.Close()
}

func (s *WebServer) Shutdown() (err error) {
	return s.fs.Shutdown()
}
