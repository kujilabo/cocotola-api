package middleware

import (
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"

	"github.com/kujilabo/cocotola-api/pkg_auth/gateway"
	"github.com/kujilabo/cocotola-api/pkg_lib/log"
)

func NewAuthMiddleware(signingKey []byte) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		logger := log.FromContext(ctx)
		authorization := c.GetHeader("Authorization")
		if !strings.HasPrefix(authorization, "Bearer ") {
			logger.Error("Has not Bearer")
			return
		}

		tokenString := authorization[len("Bearer "):]
		token, err := jwt.ParseWithClaims(tokenString, &gateway.AppUserClaims{}, func(token *jwt.Token) (interface{}, error) {
			return signingKey, nil
		})
		if err != nil {
			logger.WithError(err).Errorf("%v", err)
			return
		}

		if claims, ok := token.Claims.(*gateway.AppUserClaims); ok && token.Valid {
			c.Set("AuthorizedUser", int(claims.AppUserID))
			c.Set("OrganizationID", int(claims.OrganizationID))
			logger.Infof("AuthorizedUser %d", int(claims.AppUserID))
		}
	}
}
