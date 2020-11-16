package db

import (
	"context"
	"distate-task/dt/config"
	"distate-task/dt/logger"
	"github.com/jackc/pgx/v4/pgxpool"
	"go.uber.org/fx"
)

func Module() fx.Option {
	return fx.Options(
		fx.Provide(
			func(cfg *config.Config, log *logger.Logger) (*Connection, error) {
				pool, err := pgxpool.Connect(context.Background(), cfg.Database.DSN)
				if err != nil {
					return nil, err
				}
				log.Info("Created pgx connection pool")

				conn := &Connection{Pool: pool}
				return conn, nil
			},
		),
		fx.Invoke(
			func(lf fx.Lifecycle, dbConn *Connection, cfg *config.Config, log *logger.Logger) {
				lf.Append(fx.Hook{
					OnStop: func(ctx context.Context) error {
						dbConn.Pool.Close()
						return nil
					},
				})
			},
		),
	)
}
