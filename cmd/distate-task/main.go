package main

import (
	"distate-task/dt"
	"distate-task/dt/config"
	"distate-task/dt/logger"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"net/http"
)

var (
	// Provisioned by ldflags
	version   string
	buildDate string
)

func main() {
	conf, err := config.NewConfigFromEnv()
	if err != nil {
		panic(err)
	}

	conf.Version = version
	conf.BuildDate = buildDate

	log, err := logger.NewLogger(conf)
	if err != nil {
		panic(err)
	}

	if conf.ProfilerEnable {
		runProfiler(log)
	}

	log = log.Named(dt.Name)
	log.Info("Start distate-task")

	defer dt.Recover(log)
	fx.New(
		dt.Provide(conf, log),
	).Run()
}

func runProfiler(log *logger.Logger) {
	go func() {
		if err := http.ListenAndServe(":6060", nil); err != nil {
			log.Error("pprof http server can't start", zap.Error(err))
		}
	}()
}
