package utils

import "go.opentelemetry.io/otel"

var tracer = otel.Tracer("github.com/AH-dark/gravatar-with-qq-avatar/pkg/utils")
