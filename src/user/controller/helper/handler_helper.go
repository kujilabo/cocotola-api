package handlerhelper

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/kujilabo/cocotola-api/src/lib/log"
	"github.com/kujilabo/cocotola-api/src/user/domain"
)

func HandleSecuredFunction(c *gin.Context, fn func(organizationID domain.OrganizationID, operatorID domain.AppUserID) error, errorHandle func(c *gin.Context, err error) bool) {
	ctx := c.Request.Context()
	logger := log.FromContext(ctx)
	organizationID := domain.OrganizationID((c.GetInt("OrganizationID")))
	operatorID := domain.AppUserID(uint(c.GetInt("AuthorizedUser")))

	logger.Infof("OperatorID: %d, OrganizationID: %d", operatorID, organizationID)
	if err := fn(organizationID, operatorID); err != nil {
		if handled := errorHandle(c, err); !handled {
			c.Status(http.StatusInternalServerError)
		}
	}
}

func HandleRoleFunction(c *gin.Context, targetRole string, fn func(organizationID domain.OrganizationID, operatorID domain.AppUserID) error, errorHandle func(c *gin.Context, err error) bool) {
	ctx := c.Request.Context()
	logger := log.FromContext(ctx)
	organizationID := domain.OrganizationID((c.GetInt("OrganizationID")))
	operatorID := domain.AppUserID(uint(c.GetInt("AuthorizedUser")))
	role := c.GetString("Role")

	logger.Infof("OperatorID: %d, OrganizationID: %d, Role: %s", operatorID, organizationID, role)
	if role != targetRole {
		c.Status(http.StatusForbidden)
		return
	}

	if err := fn(organizationID, operatorID); err != nil {
		if handled := errorHandle(c, err); !handled {
			c.Status(http.StatusInternalServerError)
		}
	}
}

func HandleFunction(c *gin.Context, fn func(organizationID domain.OrganizationID, operatorID domain.AppUserID) error, errorHandle func(c *gin.Context, err error) bool) {
	ctx := c.Request.Context()
	logger := log.FromContext(ctx)
	organizationID := domain.OrganizationID((c.GetInt("OrganizationID")))
	operatorID := domain.AppUserID(uint(c.GetInt("AuthorizedUser")))
	role := c.GetString("Role")
	logger.Infof("OperatorID: %d, OrganizationID: %d, Role: %s", operatorID, organizationID, role)

	if err := fn(organizationID, operatorID); err != nil {
		if handled := errorHandle(c, err); !handled {
			c.Status(http.StatusInternalServerError)
		}
	}
}

type IDResponse struct {
	ID uint
}
