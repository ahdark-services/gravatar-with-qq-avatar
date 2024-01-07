package middlewares

import (
	"context"
	"github.com/cloudwego/hertz/pkg/app"
)

const HeaderRequestId = "X-Request-Id"

func RequestId() app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
		ctx, span := tracer.Start(ctx, "server.middlewares.RequestId")
		defer span.End()

		c.Header(HeaderRequestId, span.SpanContext().TraceID().String())
		c.Next(ctx)
	}
}
