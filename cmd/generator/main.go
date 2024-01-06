package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/AH-dark/bytestring"
	md5simd "github.com/minio/md5-simd"
	"github.com/uptrace/opentelemetry-go-extra/otelzap"
	"go.opentelemetry.io/otel"
	"go.uber.org/fx"
	"go.uber.org/zap"

	"github.com/AH-dark/gravatar-with-qq-avatar/database/dal"
	"github.com/AH-dark/gravatar-with-qq-avatar/entry"
	"github.com/AH-dark/gravatar-with-qq-avatar/pkg/cryptor"
)

var ctx = context.Background()
var tracer = otel.Tracer("github.com/AH-dark/gravatar-with-qq-avatar/cmd/generator")

var (
	from = int64(10000)
	to   = int64(9999999999)
)

func init() {
	flag.Int64Var(&from, "from", from, "from")
	flag.Int64Var(&to, "to", to, "to")
	flag.Parse()
}

func generate(ctx context.Context, md5QQMappingRepo dal.MD5QQMappingRepo) error {
	ctx, span := tracer.Start(ctx, "cmd.generator.generate")
	defer span.End()

	svr := md5simd.NewServer()
	defer svr.Close()

	for i := from; i <= to; i++ {
		md5Hash := cryptor.Md5WithServer(svr, bytestring.StringToBytes(fmt.Sprintf("%d@qq.com", i)))

		if err := md5QQMappingRepo.InsertMapping(ctx, i, bytestring.BytesToString(md5Hash)); err != nil {
			otelzap.L().Ctx(ctx).Error("insert qq avatar failed", zap.Error(err))
			return err
		}

		if i%10000 == 0 {
			otelzap.L().Ctx(ctx).Info("insert qq avatar", zap.Int64("qq_id", i))
		}
	}

	return nil
}

func main() {
	app := fx.New(
		fx.Supply(fx.Annotate(ctx, fx.As(new(context.Context)))),
		fx.Supply(fx.Annotate("generator", fx.ResultTags(`name:"serviceName"`))),
		entry.AppEntries(),
		fx.Invoke(generate),
	)

	if err := app.Start(ctx); err != nil {
		panic(err)
	}

	if err := app.Stop(ctx); err != nil {
		panic(err)
	}
}
