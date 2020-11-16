package config

const (
	DebugLevel = "debug"
	ErrorLevel = "error"
	InfoLevel  = "info"
	WarnLevel  = "warn"
)

const (
	DefaultLevel              = InfoLevel
	//DefaultTimeEncoder        = "epoch"
	DefaultOutputPath         = "stdout"
)

type Logger struct {
	Debug              bool     `envconfig:"debug" json:"debug"`
	Level              string   `envconfig:"level" json:"level"`
	Output             []string `envconfig:"output" json:"-"`
	//TimeEncoder        string   `envconfig:"time_encoder" json:"time_encoder"`
}

func (c Logger) IsDebugMode() bool {
	return c.Debug
}

func (c Logger) GetLevel() string {
	if len(c.Level) == 0 {
		return DefaultLevel
	}
	return c.Level
}

func (c Logger) GetOutput() []string {
	if len(c.Output) == 0 {
		return []string{DefaultOutputPath}
	}
	return c.Output
}