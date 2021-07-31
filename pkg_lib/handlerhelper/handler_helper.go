package handlerhelper

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/kujilabo/cocotola-api/pkg_lib/log"
	user "github.com/kujilabo/cocotola-api/pkg_user/domain"
)

func HandleSecuredFunction(c *gin.Context, fn func(organizationID user.OrganizationID, operatorID user.AppUserID) error, errorHandle func(c *gin.Context, err error) bool) {
	ctx := c.Request.Context()
	logger := log.FromContext(ctx)
	organizationID := user.OrganizationID((c.GetInt("OrganizationID")))
	operatorID := user.AppUserID(uint(c.GetInt("AuthorizedUser")))

	logger.Infof("OperatorID: %d, OrganizationID: %d", operatorID, organizationID)
	if err := fn(organizationID, operatorID); err != nil {
		if handled := errorHandle(c, err); !handled {
			c.Status(http.StatusInternalServerError)
		}
	}
}

type IDResponse struct {
	ID uint
}
