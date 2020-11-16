package dt

import (
	"distate-task/config"
	"distate-task/dt/db"
	"distate-task/dt/logger"
	"distate-task/dt/webserver"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

const Name = "distate-task"

func Provide(cfg *config.Config, log *logger.Logger) fx.Option {
	return fx.Options(
		fx.StartTimeout(cfg.StartTimeout),
		fx.StopTimeout(cfg.StopTimeout),

		fx.Logger(
			logger.FxLogger{SugaredLogger: log.Named("fx").Sugar()},
		),

		fx.Provide(
			func() *config.Config {
				return cfg
			},
			func() *logger.Logger {
				return log
			},
		),

		db.Module(),
		webserver.Module(),

		fx.Invoke(
			func(c *config.Config, log *logger.Logger) {
				log.Info("distate-task started",
					zap.String("buildDate", c.BuildDate))
			},
		),
	)
}

func Recover(log *logger.Logger) {
	if err := recover(); err != nil {
		log.Error("app recover error", zap.Any("error", err))
	}
}
