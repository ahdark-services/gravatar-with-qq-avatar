package dcron_logger

import (
	"github.com/libi/dcron/dlog"
	"go.uber.org/zap"
)

type ZapLoggerImpl struct {
	log *zap.Logger
}

var _ dlog.Logger = (*ZapLoggerImpl)(nil)

func NewZapLogger(log *zap.Logger) *ZapLoggerImpl {
	return &ZapLoggerImpl{log: log}
}

func (l *ZapLoggerImpl) Printf(s string, a ...any) {
	l.log.Sugar().Infof(s, a...)
}

func (l *ZapLoggerImpl) Infof(s string, a ...any) {
	l.log.Sugar().Infof(s, a...)
}

func (l *ZapLoggerImpl) Warnf(s string, a ...any) {
	l.log.Sugar().Warnf(s, a...)
}

func (l *ZapLoggerImpl) Errorf(s string, a ...any) {
	l.log.Sugar().Errorf(s, a...)
}
