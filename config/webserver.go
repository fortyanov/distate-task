package config

type WebServer struct {
	Host            string        `envconfig:"host" default:"0.0.0.0" json:"-"`
	Port            string        `envconfig:"port" default:"8080" json:"-"`
	//ReadTimeout     time.Duration `envconfig:"read_timeout" default:"30s" json:"read_timeout"`
	//WriteTimeout    time.Duration `envconfig:"write_timeout" default:"30s" json:"write_timeout"`
	//ShutdownTimeout time.Duration `envconfig:"shutdown_timeout" default:"10s" json:"shutdown_timeout"`
}
