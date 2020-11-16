package config

type WebServer struct {
	Host string `envconfig:"host" default:"0.0.0.0" json:"-"`
	Port string `envconfig:"port" default:"8080" json:"-"`
}
