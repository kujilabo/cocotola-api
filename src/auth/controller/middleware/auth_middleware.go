package middleware

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"

	"github.com/kujilabo/cocotola-api/src/auth/gateway"
	"github.com/kujilabo/cocotola-api/src/lib/log"
)

func NewAuthMiddleware(signingKey []byte) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		logger := log.FromContext(ctx)
		authorization := c.GetHeader("Authorization")
		if !strings.HasPrefix(authorization, "Bearer ") {
			logger.Warn("invalid header. Bearer not found")
			return
		}

		tokenString := authorization[len("Bearer "):]
		token, err := jwt.ParseWithClaims(tokenString, &gateway.AppUserClaims{}, func(token *jwt.Token) (interface{}, error) {
			return signingKey, nil
		})
		if err != nil {
			logger.WithError(err).Warnf("invalid token. err: %v", err)
			return
		}

		if claims, ok := token.Claims.(*gateway.AppUserClaims); ok && token.Valid {
			c.Set("AuthorizedUser", int(claims.AppUserID))
			c.Set("OrganizationID", int(claims.OrganizationID))
			c.Set("Role", claims.Role)

			logger.Infof("uri: %s, user: %d, role: %s", c.Request.RequestURI, int(claims.AppUserID), claims.Role)
		} else {
			logger.Warnf("invalid token")
			return
		}
	}
}
