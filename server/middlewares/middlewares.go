package middlewares

import "go.opentelemetry.io/otel"

var tracer = otel.Tracer("github.com/AH-dark/gravatar-with-qq-avatar/internal/server/middlewares")
