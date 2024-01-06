package logging

import (
	"context"
	"github.com/samber/lo"

	"github.com/spf13/viper"
	"github.com/uptrace/opentelemetry-go-extra/otelzap"
	"go.opentelemetry.io/otel"
	"go.uber.org/fx"
	"go.uber.org/fx/fxevent"
	"go.uber.org/zap"
)

var tracer = otel.Tracer("github.com/AH-dark/gravatar-with-qq-avatar/internal/logging")

type Logger struct {
	fx.Out
	OtelLogger *otelzap.Logger
	ZapLogger  *zap.Logger
}

func NewLogger(ctx context.Context, viper *viper.Viper, options []zap.Option) (Logger, error) {
	ctx, span := tracer.Start(ctx, "logging.NewLogger")
	defer span.End()

	zapLogger, err := lo.If(viper.GetBool("debug"), zap.NewDevelopment).Else(zap.NewProduction)(options...)
	if err != nil {
		span.RecordError(err)
		otelzap.L().Ctx(ctx).Error("failed to create logger", zap.Error(err))
		return Logger{}, err
	}

	otelLogger := otelzap.New(zapLogger,
		otelzap.WithCaller(viper.GetBool("log.caller")),
		otelzap.WithTraceIDField(viper.GetBool("log.trace_id")),
	)

	return Logger{
		OtelLogger: otelLogger,
		ZapLogger:  zapLogger,
	}, nil
}

func UseLogger(zapLogger *zap.Logger, otelLogger *otelzap.Logger) {
	zap.ReplaceGlobals(zapLogger)
	otelzap.ReplaceGlobals(otelLogger)
}

func FxLogger(logger *otelzap.Logger) fxevent.Logger {
	return &fxevent.ZapLogger{Logger: logger.Logger}
}
