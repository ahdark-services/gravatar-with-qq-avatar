package settings

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"strconv"
	"sync"
	"time"

	"github.com/AH-dark/bytestring"
	"github.com/uptrace/opentelemetry-go-extra/otelzap"
	clientv3 "go.etcd.io/etcd/client/v3"
	"go.etcd.io/etcd/client/v3/namespace"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

var tracer = otel.Tracer("github.com/AH-dark/gravatar-with-qq-support/internal/settings")

type implement struct {
	fx.In  `ignore-unexported:"true"`
	Client *clientv3.Client

	kv      clientv3.KV
	watcher clientv3.Watcher
	store   *sync.Map
}

func NewSettings(s implement) Settings {
	s.kv = namespace.NewKV(s.Client.KV, "settings:")
	s.watcher = namespace.NewWatcher(s.Client.Watcher, "settings:")
	s.store = &sync.Map{}

	return &s
}

func (svc *implement) handleEvent(key string, watchResp clientv3.WatchResponse) {
	ctx, span := tracer.Start(context.Background(), "Settings.handleEvent", trace.WithAttributes(
		attribute.String("setting.key", key),
	))
	defer span.End()

	if watchResp.Err() != nil {
		otelzap.L().Ctx(ctx).Error("failed to watch for changes", zap.Error(watchResp.Err()), zap.String("key", key))
		return
	}

	event := watchResp.Events[len(watchResp.Events)-1]
	otelzap.L().Ctx(ctx).Debug("setting changed", zap.String("key", key), zap.String("value", bytestring.BytesToString(event.Kv.Value)))
	switch event.Type {
	case clientv3.EventTypePut:
		svc.store.Store(key, bytestring.BytesToString(event.Kv.Value))
	case clientv3.EventTypeDelete:
		svc.store.Delete(key)
	default:
		otelzap.L().Ctx(ctx).Warn("unknown event type", zap.String("key", key), zap.String("type", event.Type.String()))
	}
}

func (svc *implement) getSetting(ctx context.Context, key string) (string, error) {
	ctx, span := tracer.Start(ctx, "Settings.getSetting", trace.WithAttributes(
		attribute.String("setting.key", key),
	))
	defer span.End()

	// Fetching from cache
	span.AddEvent("fetching-from-cache")
	if v, ok := svc.store.Load(key); ok {
		if s, ok := v.(string); ok {
			otelzap.L().Ctx(ctx).Debug("setting found in cache", zap.String("key", key), zap.String("value", s))
			return s, nil
		}
	}

	// Fetching from etcd
	span.AddEvent("fetching-from-etcd")
	resp, err := svc.kv.Get(ctx, key, clientv3.WithLastRev()...)
	if err != nil {
		span.RecordError(err)
		otelzap.L().Ctx(ctx).Error("failed to get setting from etcd", zap.Error(err), zap.String("key", key))
		return "", fmt.Errorf("failed to get setting from etcd: %w", err)
	}

	// If key is not present in etcd
	if resp.Count == 0 {
		span.SetStatus(codes.Error, "key not found")
		otelzap.L().Ctx(ctx).Error("key not found", zap.Error(err), zap.String("key", key))
		return "", errors.New("key not found")
	}

	// Store in cache
	span.AddEvent("store-in-cache")
	svc.store.Store(key, bytestring.BytesToString(resp.Kvs[0].Value))

	// Watch for changes
	span.AddEvent("watch-for-changes")
	go func(key string) {
		defer svc.store.Delete(key)

		for watchResp := range svc.watcher.Watch(ctx, key) {
			svc.handleEvent(key, watchResp)
		}
	}(key)

	return bytestring.BytesToString(resp.Kvs[0].Value), nil
}

func (svc *implement) GetString(ctx context.Context, key string) (string, error) {
	ctx, span := tracer.Start(ctx, "Settings.GetString", trace.WithAttributes(
		attribute.String("setting.key", key),
	))
	defer span.End()

	s, err := svc.getSetting(ctx, key)
	if err != nil {
		span.RecordError(err)
		otelzap.L().Ctx(ctx).Error("failed to get setting", zap.Error(err), zap.String("key", key))
		return "", fmt.Errorf("failed to get setting: %w", err)
	}

	return s, nil
}

func (svc *implement) GetInt64(ctx context.Context, key string) (int64, error) {
	ctx, span := tracer.Start(ctx, "Settings.GetInt64", trace.WithAttributes(
		attribute.String("setting.key", key),
	))
	defer span.End()

	s, err := svc.getSetting(ctx, key)
	if err != nil {
		span.RecordError(err)
		otelzap.L().Ctx(ctx).Error("failed to get setting", zap.Error(err), zap.String("key", key))
		return 0, fmt.Errorf("failed to get setting: %w", err)
	}

	i, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		span.RecordError(err)
		otelzap.L().Ctx(ctx).Error("failed to parse setting", zap.Error(err), zap.String("key", key))
		return 0, fmt.Errorf("failed to parse setting: %w", err)
	}

	return i, nil
}

func (svc *implement) GetUint64(ctx context.Context, key string) (uint64, error) {
	ctx, span := tracer.Start(ctx, "Settings.GetUint64", trace.WithAttributes(
		attribute.String("setting.key", key),
	))
	defer span.End()

	s, err := svc.getSetting(ctx, key)
	if err != nil {
		span.RecordError(err)
		otelzap.L().Ctx(ctx).Error("failed to get setting", zap.Error(err), zap.String("key", key))
		return 0, fmt.Errorf("failed to get setting: %w", err)
	}

	i, err := strconv.ParseUint(s, 10, 64)
	if err != nil {
		span.RecordError(err)
		otelzap.L().Ctx(ctx).Error("failed to parse setting", zap.Error(err), zap.String("key", key))
		return 0, fmt.Errorf("failed to parse setting: %w", err)
	}

	return i, nil
}

func (svc *implement) GetBool(ctx context.Context, key string) (bool, error) {
	ctx, span := tracer.Start(ctx, "Settings.GetBool", trace.WithAttributes(
		attribute.String("setting.key", key),
	))
	defer span.End()

	s, err := svc.getSetting(ctx, key)
	if err != nil {
		span.RecordError(err)
		otelzap.L().Ctx(ctx).Error("failed to get setting", zap.Error(err), zap.String("key", key))
		return false, fmt.Errorf("failed to get setting: %w", err)
	}

	b, err := strconv.ParseBool(s)
	if err != nil {
		span.RecordError(err)
		otelzap.L().Ctx(ctx).Error("failed to parse setting", zap.Error(err), zap.String("key", key))
		return false, fmt.Errorf("failed to parse setting: %w", err)
	}

	return b, nil
}

func (svc *implement) GetTimeDuration(ctx context.Context, key string) (time.Duration, error) {
	ctx, span := tracer.Start(ctx, "Settings.GetTimeDuration", trace.WithAttributes(
		attribute.String("setting.key", key),
	))
	defer span.End()

	s, err := svc.getSetting(ctx, key)
	if err != nil {
		span.RecordError(err)
		otelzap.L().Ctx(ctx).Error("failed to get setting", zap.Error(err), zap.String("key", key))
		return 0, fmt.Errorf("failed to get setting: %w", err)
	}

	d, err := time.ParseDuration(s)
	if err != nil {
		span.RecordError(err)
		otelzap.L().Ctx(ctx).Error("failed to parse setting", zap.Error(err), zap.String("key", key))
		return 0, fmt.Errorf("failed to parse setting: %w", err)
	}

	return d, nil
}

func (svc *implement) GetUrl(ctx context.Context, key string) (*url.URL, error) {
	ctx, span := tracer.Start(ctx, "Settings.GetUrl", trace.WithAttributes(
		attribute.String("setting.key", key),
	))
	defer span.End()

	s, err := svc.getSetting(ctx, key)
	if err != nil {
		span.RecordError(err)
		otelzap.L().Ctx(ctx).Error("failed to get setting", zap.Error(err), zap.String("key", key))
		return nil, fmt.Errorf("failed to get setting: %w", err)
	}

	u, err := url.Parse(s)
	if err != nil {
		span.RecordError(err)
		otelzap.L().Ctx(ctx).Error("failed to parse setting", zap.Error(err), zap.String("key", key))
		return nil, fmt.Errorf("failed to parse setting: %w", err)
	}

	return u, nil
}

func (svc *implement) ListStrings(ctx context.Context, keys []string) (map[string]string, error) {
	ctx, span := tracer.Start(ctx, "Settings.ListStrings", trace.WithAttributes(
		attribute.StringSlice("setting.keys", keys),
	))
	defer span.End()

	data := make(map[string]string, len(keys))
	for _, key := range keys {
		s, err := svc.getSetting(ctx, key)
		if err != nil {
			span.RecordError(err)
			otelzap.L().Ctx(ctx).Error("failed to get setting", zap.Error(err), zap.String("key", key))
			return nil, fmt.Errorf("failed to get setting: %w, key %s", err, key)
		}

		data[key] = s
		otelzap.L().Ctx(ctx).Debug("setting found", zap.String("key", key), zap.String("value", s))
	}

	return data, nil
}

func (svc *implement) SaveString(ctx context.Context, key string, value string) error {
	ctx, span := tracer.Start(ctx, "Settings.SaveString")
	defer span.End()

	resp, err := svc.kv.Put(ctx, key, value)
	if err != nil {
		span.RecordError(err)
		otelzap.L().Ctx(ctx).Error("failed to update setting", zap.Error(err), zap.String("key", key), zap.String("value", value))
		return fmt.Errorf("failed to update setting: %w", err)
	}

	otelzap.L().Ctx(ctx).Debug("setting updated", zap.String("key", key), zap.String("value", value), zap.Int64("revision", resp.Header.Revision))

	return nil
}
