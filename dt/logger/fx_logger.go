package logger

import "go.uber.org/zap"

type FxLogger struct {
	*zap.SugaredLogger
}

func (fxl FxLogger) Printf(format string, v ...interface{}) {
	fxl.Infof(format, v...)
}
