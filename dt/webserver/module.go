package webserver

import (
	"context"
	"distate-task/config"
	"distate-task/dt/db"
	"distate-task/dt/logger"
	"github.com/fasthttp/router"
	"github.com/valyala/fasthttp/reuseport"
	fx "go.uber.org/fx"
	"net"
)

func Module() fx.Option {
	return fx.Options(
		fx.Provide(
			func(cfg *config.Config, l *logger.Logger, conn *db.Connection) (s *WebServer, err error) {
				s = &WebServer{
					Config: cfg.WebServer,
					Log:    &logger.FxLogger{SugaredLogger: l.Named("WebServer").Sugar()},
					router: router.New(),
					debug:  cfg.Debug,
					conn: conn,
				}
				s.RegisterRoutes()
				s.ln, err = reuseport.Listen("tcp4", net.JoinHostPort(s.Config.Host, s.Config.Port))
				if err != nil {
					return nil, err
				}
				return s, nil
			},
		),
		fx.Invoke(
			func(lf fx.Lifecycle, s *WebServer, cfg *config.Config, log *logger.Logger) {
				lf.Append(fx.Hook{
					OnStart: func(ctx context.Context) error {
						go s.Run()
						return nil
					},
					OnStop: func(ctx context.Context) error {
						err := s.Shutdown()
						if err != nil {
							return err
						}
						return s.Close()
					},
				})
			},
		),
	)
}
