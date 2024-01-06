package config

import (
	"context"

	"github.com/google/uuid"
	"github.com/spf13/viper"
	"github.com/uptrace/opentelemetry-go-extra/otelzap"
	"go.opentelemetry.io/otel"
	"go.uber.org/zap"
)

var tracer = otel.Tracer("github.com/AH-dark/gravatar-with-qq-support/internal/config")

var defaultConfig = map[string]interface{}{
	"name":        "gravatar-with-qq-support",
	"namespace":   "gravatar-with-qq-support",
	"version":     "0.0.1",
	"instance_id": uuid.New(),
}

func NewViper(ctx context.Context) (*viper.Viper, error) {
	ctx, span := tracer.Start(ctx, "config.NewViper")
	defer span.End()

	vip := viper.GetViper()
	vip.AutomaticEnv()

	vip.SetConfigName("config")
	vip.SetConfigType("yaml")
	vip.AddConfigPath(".")
	vip.AddConfigPath("./config")
	vip.AddConfigPath("/etc/gravatar-with-qq-support")
	vip.AddConfigPath("$HOME/.gravatar-with-qq-support")

	for k, v := range defaultConfig {
		vip.SetDefault(k, v)
	}

	if err := vip.ReadInConfig(); err != nil {
		span.RecordError(err)
		otelzap.L().Ctx(ctx).Error("failed to read config", zap.Error(err))
		return nil, err
	}

	return vip, nil
}
