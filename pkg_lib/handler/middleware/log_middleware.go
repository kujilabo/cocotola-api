package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/kujilabo/cocotola-api/pkg_lib/log"
)

func NewLogMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tmp, err := uuid.NewRandom()
		if err != nil {
			return
		}
		ctx := log.With(c.Request.Context(), log.Str("request_id", tmp.String()))
		c.Request = c.Request.WithContext(ctx)
		logger := log.FromContext(ctx)
		logger.Infof("uri: %s, method: %s", c.Request.RequestURI, c.Request.Method)
	}
}
