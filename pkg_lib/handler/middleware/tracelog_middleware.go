package middleware

import (
	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel/attribute"
)

func NewTraceLogMiddleware(appName string) gin.HandlerFunc {
	return func(c *gin.Context) {
		requestID := c.GetString("request_id")

		savedCtx := c.Request.Context()
		defer func() {
			c.Request = c.Request.WithContext(savedCtx)
		}()

		ctx, span := tracer.Start(c.Request.Context(), "TraceLog")
		defer span.End()

		span.SetAttributes(attribute.String(appName+".request_id", requestID))

		c.Request = c.Request.WithContext(ctx)

		// serve the request to the next middleware
		c.Next()
	}
}
