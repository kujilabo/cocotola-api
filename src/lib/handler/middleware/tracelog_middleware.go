package middleware

import (
	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel/trace"

	"github.com/kujilabo/cocotola-api/src/lib/log"
)

func NewTraceLogMiddleware(appName string) gin.HandlerFunc {
	return func(c *gin.Context) {
		sc := trace.SpanFromContext(c.Request.Context()).SpanContext()
		if !sc.TraceID().IsValid() || !sc.SpanID().IsValid() {
			return
		}
		otTraceID := sc.TraceID().String()

		ctx := log.With(c.Request.Context(), log.Str("request_id", otTraceID))

		savedCtx := ctx
		defer func() {
			c.Request = c.Request.WithContext(savedCtx)
		}()

		ctx, span := tracer.Start(ctx, "TraceLog")
		defer span.End()

		c.Request = c.Request.WithContext(ctx)

		// serve the request to the next middleware
		c.Next()
	}
}
