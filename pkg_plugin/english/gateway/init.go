package gateway

import "go.opentelemetry.io/otel"

var tracer = otel.Tracer("github.com/kujilabo/cocotola-api/pkg_app/common/gateway")
