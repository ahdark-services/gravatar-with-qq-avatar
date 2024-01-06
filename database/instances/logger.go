package instances

import (
	"go.uber.org/zap"
)

type StdLogger struct {
	logger *zap.Logger
}

func (l *StdLogger) Print(v ...interface{}) {
	l.logger.Sugar().Debug(v...)
}

func (l *StdLogger) Printf(format string, v ...interface{}) {
	l.logger.Sugar().Debugf(format, v...)
}

func (l *StdLogger) Println(v ...interface{}) {
	l.logger.Sugar().Debugln(v...)
}
