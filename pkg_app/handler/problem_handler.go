package handler

import (
	"context"
	"errors"
	"io"
	"net/http"
	"strings"

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

type ProblemHandler interface {
	FindProblems(c *gin.Context)

	FindAllProblems(c *gin.Context)

	FindProblemsByProblemIDs(c *gin.Context)

	FindProblemByID(c *gin.Context)

	AddProblem(c *gin.Context)

	// FindProblemIDs(c *gin.Context)

	ImportProblems(c *gin.Context)

	// UpdateProblem(c *gin.Context)

	RemoveProblem(c *gin.Context)
}

type problemHandler struct {
	problemService application.ProblemService
	newIterator    func(ctx context.Context, workbookID domain.WorkbookID, problemType string, reader io.Reader) (domain.ProblemAddParameterIterator, error)
}

func NewProblemHandler(problemService application.ProblemService, newIterator func(ctx context.Context, workbookID domain.WorkbookID, problemType string, reader io.Reader) (domain.ProblemAddParameterIterator, error)) ProblemHandler {
	return &problemHandler{
		problemService: problemService,
		newIterator:    newIterator,
	}
}

func (h *problemHandler) FindProblems(c *gin.Context) {
	ctx := c.Request.Context()
	logger := log.FromContext(ctx)
	logger.Info("FindProblems")

	handlerhelper.HandleSecuredFunction(c, func(organizationID user.OrganizationID, operatorID user.AppUserID) error {
		workbookID, err := ginhelper.GetUint(c, "workbookID")
		if err != nil {
			c.Status(http.StatusBadRequest)
			return nil
		}

		param := entity.ProblemSearchParameter{}
		if err := c.ShouldBindJSON(&param); err != nil {
			c.Status(http.StatusBadRequest)
			return nil
		}

		parameter, err := converter.ToProblemSearchCondition(ctx, &param, domain.WorkbookID(workbookID))
		if err != nil {
			return err
		}

		result, err := h.problemService.FindProblemsByWorkbookID(ctx, organizationID, operatorID, domain.WorkbookID(workbookID), parameter)
		if err != nil {
			return err
		}

		response, err := converter.ToProblemSearchResponse(ctx, result)
		if err != nil {
			return err
		}

		c.JSON(http.StatusOK, response)
		return nil
	}, h.errorHandle)
}

func (h *problemHandler) FindAllProblems(c *gin.Context) {
	ctx := c.Request.Context()
	logger := log.FromContext(ctx)
	logger.Info("FindAllProblems")

	handlerhelper.HandleSecuredFunction(c, func(organizationID user.OrganizationID, operatorID user.AppUserID) error {
		workbookID, err := ginhelper.GetUint(c, "workbookID")
		if err != nil {
			c.Status(http.StatusBadRequest)
			return nil
		}

		result, err := h.problemService.FindAllProblemsByWorkbookID(ctx, organizationID, operatorID, domain.WorkbookID(workbookID))
		if err != nil {
			return err
		}

		response, err := converter.ToProblemSearchResponse(ctx, result)
		if err != nil {
			return err
		}

		c.JSON(http.StatusOK, response)
		return nil
	}, h.errorHandle)
}

func (h *problemHandler) FindProblemsByProblemIDs(c *gin.Context) {
	ctx := c.Request.Context()
	logger := log.FromContext(ctx)
	logger.Info("FindProblemsByProblemIDs")

	handlerhelper.HandleSecuredFunction(c, func(organizationID user.OrganizationID, operatorID user.AppUserID) error {
		workbookID, err := ginhelper.GetUint(c, "workbookID")
		if err != nil {
			c.Status(http.StatusBadRequest)
			return nil
		}

		param := entity.ProblemIDsParameter{}
		if err := c.ShouldBindJSON(&param); err != nil {
			c.Status(http.StatusBadRequest)
			return nil
		}

		parameter, err := converter.ToProblemIDsCondition(ctx, &param, domain.WorkbookID(workbookID))
		if err != nil {
			return err
		}

		result, err := h.problemService.FindProblemsByProblemIDs(ctx, organizationID, operatorID, domain.WorkbookID(workbookID), parameter)
		if err != nil {
			return err
		}

		response, err := converter.ToProblemSearchResponse(ctx, result)
		if err != nil {
			return err
		}

		c.JSON(http.StatusOK, response)
		return nil
	}, h.errorHandle)
}

func (h *problemHandler) FindProblemByID(c *gin.Context) {
	ctx := c.Request.Context()
	logger := log.FromContext(ctx)
	logger.Info("FindProblemByID")

	handlerhelper.HandleSecuredFunction(c, func(organizationID user.OrganizationID, operatorID user.AppUserID) error {
		workbookID, err := ginhelper.GetUint(c, "workbookID")
		if err != nil {
			c.Status(http.StatusBadRequest)
			return nil
		}

		problemID, err := ginhelper.GetUint(c, "problemID")
		if err != nil {
			c.Status(http.StatusBadRequest)
			return nil
		}

		result, err := h.problemService.FindProblemByID(ctx, organizationID, operatorID, domain.WorkbookID(workbookID), domain.ProblemID(problemID))
		if err != nil {
			return err
		}

		response, err := converter.ToProblemResponse(ctx, result)
		if err != nil {
			return err
		}

		c.JSON(http.StatusOK, response)
		return nil
	}, h.errorHandle)
}

// func (h *problemHandler) FindProblemIDs(c *gin.Context) {
// 	ctx := c.Request.Context()
// 	logger := log.FromContext(ctx)
// 	logger.Info("FindProblemIDs")

// 	handlerhelper.HandleSecuredFunction(c, func(organizationID user.OrganizationID, operatorID user.AppUserID) error {
// 		workbookID, err := ginhelper.GetUint(c, "workbookID")
// 		if err != nil {
// 			c.Status(http.StatusBadRequest)
// 			return nil
// 		}

// 		// result, err := h.problemService.FindProblemIDs(ctx, organizationID, operatorID, domain.WorkbookID(workbookID))
// 		// if err != nil {
// 		// 	return err
// 		// }
// 		h.studyService.Find

// 		response, err := converter.ToProblemWithLevelList(ctx, result)
// 		if err != nil {
// 			return err
// 		}

// 		c.JSON(http.StatusOK, response)
// 		return nil
// 	}, h.errorHandle)
// }

func (h *problemHandler) AddProblem(c *gin.Context) {
	ctx := c.Request.Context()
	logger := log.FromContext(ctx)
	logger.Infof("AddProblem")

	handlerhelper.HandleSecuredFunction(c, func(organizationID user.OrganizationID, operatorID user.AppUserID) error {
		workbookID, err := ginhelper.GetUint(c, "workbookID")
		if err != nil {
			c.Status(http.StatusBadRequest)
			return nil
		}

		param := entity.ProblemAddParameter{}
		if err := c.BindJSON(&param); err != nil {
			logger.Infof("failed to BindJSON. err: %v", err)
			return nil
		}

		parameter, err := converter.ToProblemAddParameter(domain.WorkbookID(workbookID), &param)
		if err != nil {
			return err
		}

		problemID, err := h.problemService.AddProblem(ctx, organizationID, operatorID, parameter)
		if err != nil {
			return xerrors.Errorf("failed to AddProblem. param: %+v, err: %w", parameter, err)
		}

		c.JSON(http.StatusOK, gin.H{"id": problemID})
		return nil
	}, h.errorHandle)
}

func (h *problemHandler) RemoveProblem(c *gin.Context) {
	ctx := c.Request.Context()
	logger := log.FromContext(ctx)

	handlerhelper.HandleSecuredFunction(c, func(organizationID user.OrganizationID, operatorID user.AppUserID) error {
		logger.Infof("RemvoeProblem. organizationID: %d, operatorID: %d", organizationID, operatorID)

		workbookID, err := ginhelper.GetUint(c, "workbookID")
		if err != nil {
			c.Status(http.StatusBadRequest)
			return nil
		}

		problemID, err := ginhelper.GetUint(c, "problemID")
		if err != nil {
			c.Status(http.StatusBadRequest)
			return nil
		}

		version, err := ginhelper.GetIntFromQuery(c, "version")
		if err != nil {
			c.Status(http.StatusBadRequest)
			return nil
		}

		if err := h.problemService.RemoveProblem(ctx, organizationID, operatorID, domain.WorkbookID(workbookID), domain.ProblemID(problemID), version); err != nil {
			return xerrors.Errorf("failed to RemoveProblem. err: %w", err)
		}

		c.Status(http.StatusNoContent)
		return nil
	}, h.errorHandle)
}

func (h *problemHandler) ImportProblems(c *gin.Context) {
	ctx := c.Request.Context()
	logger := log.FromContext(ctx)
	logger.Infof("ImportProblems")

	handlerhelper.HandleSecuredFunction(c, func(organizationID user.OrganizationID, operatorID user.AppUserID) error {
		workbookID, err := ginhelper.GetUint(c, "workbookID")
		if err != nil {
			c.Status(http.StatusBadRequest)
			return nil
		}

		contentType := c.GetHeader("Content-Type")
		if !strings.HasPrefix(contentType, "multipart/form-data") {
			logger.Warnf("contentType: %s", contentType)
			c.Status(http.StatusBadRequest)
			return nil
		}

		file, err := c.FormFile("file")
		if err != nil {
			if errors.Is(err, http.ErrMissingFile) {
				logger.Warnf("err: %+v", err)
				c.Status(http.StatusBadRequest)
				return nil
			}
			return err
		}

		logger.Infof("fileName: %s", file.Filename)
		multipartFile, err := file.Open()
		if err != nil {
			return xerrors.Errorf("failed to file.Open. err: %w", err)
		}
		defer multipartFile.Close()

		newIterator := func(workbookID domain.WorkbookID, problemType string) (domain.ProblemAddParameterIterator, error) {
			return h.newIterator(ctx, workbookID, problemType, multipartFile)
		}

		if err := h.problemService.ImportProblems(ctx, organizationID, operatorID, domain.WorkbookID(workbookID), newIterator); err != nil {
			return xerrors.Errorf("failed to ImportProblems. err: %w", err)
		}

		c.Status(http.StatusOK)
		return nil
	}, h.errorHandle)
}

func (h *problemHandler) errorHandle(c *gin.Context, err error) bool {
	ctx := c.Request.Context()
	logger := log.FromContext(ctx)
	if xerrors.Is(err, domain.ErrProblemAlreadyExists) {
		c.JSON(http.StatusConflict, gin.H{"message": "Problem already exists"})
		return true
	} else if xerrors.Is(err, domain.ErrWorkbookNotFound) {
		c.JSON(http.StatusNotFound, gin.H{"message": err.Error()})
		return true
	}
	logger.Errorf("workbookHandler error:%v", err)
	return false
}
