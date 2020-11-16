package logger

import (
	"distate-task/dt/config"
	"fmt"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"path/filepath"
	"runtime"
)

type Logger struct {
	*zap.Logger
}

func NewLogger(cfg *config.Config) (*Logger, error) {
	zapConf := zap.Config{
		Level:       zapLevel(cfg.Logger.GetLevel()),
		Development: cfg.Logger.IsDebugMode(),
		Sampling: &zap.SamplingConfig{
			Initial:    100,
			Thereafter: 100,
		},
		Encoding:         "json",
		EncoderConfig:    zapEncoderConfig(),
		OutputPaths:      cfg.Logger.GetOutput(),
		ErrorOutputPaths: []string{"stderr"},
	}

	zl, err := zapConf.Build()
	if err != nil {
		return nil, err
	}

	logger := &Logger{
		Logger: zl,
	}

	return logger, nil
}

func zapLevel(level string) zap.AtomicLevel {
	switch level {
	case config.DebugLevel:
		return zap.NewAtomicLevelAt(zapcore.DebugLevel)
	case config.ErrorLevel:
		return zap.NewAtomicLevelAt(zapcore.ErrorLevel)
	case config.InfoLevel:
		return zap.NewAtomicLevelAt(zapcore.InfoLevel)
	case config.WarnLevel:
		return zap.NewAtomicLevelAt(zapcore.WarnLevel)
	default:
		panic("unknown log level")
	}
}

func (l *Logger) caller() (caller string) {
	_, file, line, ok := runtime.Caller(3)
	if !ok {
		return
	}
	_, f := filepath.Split(file)
	caller = fmt.Sprintf("%s:%d", f, line)
	return
}

func (l *Logger) fields() []zap.Field {
	zfs := []zap.Field{
		zap.String("caller", l.caller()),
		// TODO: добавить ещё метаинформацию в логи
	}
	return zfs
}

func zapEncoderConfig() zapcore.EncoderConfig {
	return zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "logger",
		MessageKey:     "message",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     zapcore.RFC3339NanoTimeEncoder,
		EncodeDuration: zapcore.NanosDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}
}

func (l *Logger) Info(msg string, fields ...zap.Field) {
	fields = append(fields, l.fields()...)
	l.Logger.Info(msg, fields...)
}

func (l *Logger) Error(msg string, fields ...zap.Field) {
	fields = append(fields, l.fields()...)
	l.Logger.Error(msg, fields...)
}

func (l *Logger) Debug(msg string, fields ...zap.Field) {
	fields = append(fields, l.fields()...)
	l.Logger.Debug(msg, fields...)
}

func (l *Logger) Named(n string) *Logger {
	nl := l.Logger.Named(n)
	c := *l
	c.Logger = nl
	return &c
}
