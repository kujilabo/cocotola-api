package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/xerrors"

	"github.com/kujilabo/cocotola-api/pkg_app/application"
	"github.com/kujilabo/cocotola-api/pkg_app/domain"
	"github.com/kujilabo/cocotola-api/pkg_app/handler/converter"
	"github.com/kujilabo/cocotola-api/pkg_app/handler/entity"
	"github.com/kujilabo/cocotola-api/pkg_lib/ginhelper"
	"github.com/kujilabo/cocotola-api/pkg_lib/handlerhelper"
	"github.com/kujilabo/cocotola-api/pkg_lib/log"
	user "github.com/kujilabo/cocotola-api/pkg_user/domain"
)

type StudyHandler interface {
	FindRecordbook(c *gin.Context)
	SetStudyResult(c *gin.Context)
}

type studyHandler struct {
	studyService application.StudyService
}

func NewStudyHandler(studyService application.StudyService) StudyHandler {
	return &studyHandler{
		studyService: studyService,
	}
}

func (h *studyHandler) FindRecordbook(c *gin.Context) {
	ctx := c.Request.Context()
	logger := log.FromContext(ctx)
	logger.Info("FindRecordbook")

	handlerhelper.HandleSecuredFunction(c, func(organizationID user.OrganizationID, operatorID user.AppUserID) error {
		workbookID, err := ginhelper.GetUint(c, "workbookID")
		if err != nil {
			c.Status(http.StatusBadRequest)
			return nil
		}
		studyType := ginhelper.GetString(c, "studyType")

		result, err := h.studyService.FindResults(ctx, organizationID, operatorID, domain.WorkbookID(workbookID), studyType)
		if err != nil {
			return err
		}

		response, err := converter.ToProblemWithLevelList(ctx, result)
		if err != nil {
			return err
		}

		c.JSON(http.StatusOK, response)
		return nil
	}, h.errorHandle)
}

func (h *studyHandler) SetStudyResult(c *gin.Context) {
	ctx := c.Request.Context()
	logger := log.FromContext(ctx)
	logger.Info("SetStudyResult")

	handlerhelper.HandleSecuredFunction(c, func(organizationID user.OrganizationID, operatorID user.AppUserID) error {
		workbookID, err := ginhelper.GetUint(c, "workbookID")
		if err != nil {
			c.Status(http.StatusBadRequest)
			return nil
		}
		studyType := ginhelper.GetString(c, "studyType")
		problemID, err := ginhelper.GetUint(c, "problemID")
		if err != nil {
			return err
		}

		param := entity.StudyResultParameter{}
		if err := c.ShouldBindJSON(&param); err != nil {
			c.Status(http.StatusBadRequest)
			return nil
		}

		// recordbook, err := h.studyService.FindRecordbook(ctx, organizationID, operatorID, domain.WorkbookID(workbookID), studyType)
		// if err != nil {
		// 	return err
		// }

		if err := h.studyService.SetResult(ctx, organizationID, operatorID, domain.WorkbookID(workbookID), studyType, domain.ProblemID(problemID), param.Result); err != nil {
			return xerrors.Errorf("failed to SetResult. err: %w", err)
		}

		c.Status(http.StatusOK)
		return nil
	}, h.errorHandle)
}

func (h *studyHandler) errorHandle(c *gin.Context, err error) bool {
	ctx := c.Request.Context()
	logger := log.FromContext(ctx)
	if xerrors.Is(err, domain.ErrProblemAlreadyExists) {
		c.JSON(http.StatusConflict, gin.H{"message": "Problem already exists"})
		return true
	} else if xerrors.Is(err, domain.ErrWorkbookNotFound) {
		c.JSON(http.StatusNotFound, gin.H{"message": err.Error()})
		return true
	}
	logger.Errorf("studyHandler error:%v", err)
	return false
}
