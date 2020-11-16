package config

import (
	"time"

	"github.com/kelseyhightower/envconfig"
)

const CoreEnvironmentPrefix = "dt"

type Config struct {
	Version        string
	BuildDate      string
	Debug          bool          `envconfig:"debug"`
	ProfilerEnable bool          `envconfig:"pprof"`
	StartTimeout   time.Duration `envconfig:"start_timeout" default:"30s"`
	StopTimeout    time.Duration `envconfig:"stop_timeout" default:"60s"`
	Logger         Logger        `envconfig:"logger"`
	Database       Database      `envconfig:"database"`
	WebServer      WebServer     `envconfig:"webserver"`
}

func NewConfigFromEnv() (*Config, error) {
	cfg := &Config{}
	if err := envconfig.Process(CoreEnvironmentPrefix, cfg); err != nil {
		return nil, err
	}

	if cfg.Debug {
		cfg.Logger.Level = "debug"
		cfg.Logger.Debug = true
	}
	return cfg, nil
}
