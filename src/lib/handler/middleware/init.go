package middleware

import "go.opentelemetry.io/otel"

var tracer = otel.Tracer("github.com/kujilabo/cocotola-api/src/lib/handler/middleware")
