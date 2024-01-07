package avatar

import (
	"github.com/imroc/req/v3"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

func WithTracer(rt req.RoundTripper) req.RoundTripFunc {
	return func(req *req.Request) (resp *req.Response, err error) {
		ctx := req.Context()

		ctx, span := tracer.Start(ctx, "req.httpRequest", trace.WithAttributes(
			attribute.String("http.url", req.URL.String()),
			attribute.String("http.method", req.Method),
			attribute.String("http.req.header", req.HeaderToString()),
		))
		defer span.End()

		if len(req.Body) > 0 {
			span.SetAttributes(
				attribute.String("http.req.body", string(req.Body)),
			)
		}

		resp, err = rt.RoundTrip(req)
		if err != nil {
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())
		}

		if resp.Response != nil {
			span.SetAttributes(
				attribute.Int("http.status_code", resp.StatusCode),
				attribute.String("http.resp.header", resp.HeaderToString()),
			)
		}

		return
	}
}
