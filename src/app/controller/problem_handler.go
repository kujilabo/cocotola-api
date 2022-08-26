package controller

import (
	"context"
	"errors"
	"io"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"golang.org/x/xerrors"

	"github.com/kujilabo/cocotola-api/src/app/controller/converter"
	"github.com/kujilabo/cocotola-api/src/app/controller/entity"
	"github.com/kujilabo/cocotola-api/src/app/domain"
	"github.com/kujilabo/cocotola-api/src/app/service"
	studentU "github.com/kujilabo/cocotola-api/src/app/usecase/student"
	libD "github.com/kujilabo/cocotola-api/src/lib/domain"
	"github.com/kujilabo/cocotola-api/src/lib/ginhelper"
	"github.com/kujilabo/cocotola-api/src/lib/log"
	controllerhelper "github.com/kujilabo/cocotola-api/src/user/controller/helper"
	userD "github.com/kujilabo/cocotola-api/src/user/domain"
)

type ProblemHandler interface {
	FindProblems(c *gin.Context)

	FindAllProblems(c *gin.Context)

	FindProblemsByProblemIDs(c *gin.Context)

	FindProblemByID(c *gin.Context)

	AddProblem(c *gin.Context)

	// FindProblemIDs(c *gin.Context)

	ImportProblems(c *gin.Context)

	UpdateProblem(c *gin.Context)

	RemoveProblem(c *gin.Context)
}

type problemHandler struct {
	studentUsecaseProblem studentU.StudentUsecaseProblem
	newIterator           func(ctx context.Context, workbookID domain.WorkbookID, problemType string, reader io.Reader) (service.ProblemAddParameterIterator, error)
}

func NewProblemHandler(studentUsecaseProblem studentU.StudentUsecaseProblem, newIterator func(ctx context.Context, workbookID domain.WorkbookID, problemType string, reader io.Reader) (service.ProblemAddParameterIterator, error)) ProblemHandler {
	return &problemHandler{
		studentUsecaseProblem: studentUsecaseProblem,
		newIterator:           newIterator,
	}
}

func (h *problemHandler) FindProblems(c *gin.Context) {
	ctx := c.Request.Context()
	logger := log.FromContext(ctx)
	logger.Info("FindProblems")

	controllerhelper.HandleSecuredFunction(c, func(organizationID userD.OrganizationID, operatorID userD.AppUserID) error {
		workbookID, err := ginhelper.GetUintFromPath(c, "workbookID")
		if err != nil {
			c.Status(http.StatusBadRequest)
			return nil
		}

		param := entity.ProblemFindParameter{}
		if err := c.ShouldBindJSON(&param); err != nil {
			c.Status(http.StatusBadRequest)
			return nil
		}

		parameter, err := converter.ToProblemSearchCondition(ctx, &param, domain.WorkbookID(workbookID))
		if err != nil {
			return err
		}

		result, err := h.studentUsecaseProblem.FindProblemsByWorkbookID(ctx, organizationID, operatorID, domain.WorkbookID(workbookID), parameter)
		if err != nil {
			return err
		}

		response, err := converter.ToProblemFindResponse(ctx, result)
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

	controllerhelper.HandleSecuredFunction(c, func(organizationID userD.OrganizationID, operatorID userD.AppUserID) error {
		workbookID, err := ginhelper.GetUintFromPath(c, "workbookID")
		if err != nil {
			c.Status(http.StatusBadRequest)
			return nil
		}

		result, err := h.studentUsecaseProblem.FindAllProblemsByWorkbookID(ctx, organizationID, operatorID, domain.WorkbookID(workbookID))
		if err != nil {
			return err
		}

		response, err := converter.ToProblemFindAllResponse(ctx, result)
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

	controllerhelper.HandleSecuredFunction(c, func(organizationID userD.OrganizationID, operatorID userD.AppUserID) error {
		workbookID, err := ginhelper.GetUintFromPath(c, "workbookID")
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

		result, err := h.studentUsecaseProblem.FindProblemsByProblemIDs(ctx, organizationID, operatorID, domain.WorkbookID(workbookID), parameter)
		if err != nil {
			return err
		}

		response, err := converter.ToProblemFindResponse(ctx, result)
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

	controllerhelper.HandleSecuredFunction(c, func(organizationID userD.OrganizationID, operatorID userD.AppUserID) error {
		id, err := h.toProblemSelectParameter1(c)
		if err != nil {
			c.Status(http.StatusBadRequest)
			return nil
		}

		result, err := h.studentUsecaseProblem.FindProblemByID(ctx, organizationID, operatorID, id)
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

// 	controllerhelper.HandleSecuredFunction(c, func(organizationID userD.OrganizationID, operatorID userD.AppUserID) error {
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

	controllerhelper.HandleSecuredFunction(c, func(organizationID userD.OrganizationID, operatorID userD.AppUserID) error {
		workbookID, err := ginhelper.GetUintFromPath(c, "workbookID")
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

		problemID, err := h.studentUsecaseProblem.AddProblem(ctx, organizationID, operatorID, parameter)
		if err != nil {
			return xerrors.Errorf("failed to AddProblem. param: %+v, err: %w", parameter, err)
		}

		c.JSON(http.StatusOK, gin.H{"id": problemID})
		return nil
	}, h.errorHandle)
}

func (h *problemHandler) UpdateProblem(c *gin.Context) {
	ctx := c.Request.Context()
	logger := log.FromContext(ctx)
	logger.Infof("UpdateProblem")

	controllerhelper.HandleSecuredFunction(c, func(organizationID userD.OrganizationID, operatorID userD.AppUserID) error {
		id, err := h.toProblemSelectParameter2(c)
		if err != nil {
			c.Status(http.StatusBadRequest)
			return nil
		}

		param := entity.ProblemUpdateParameter{}
		if err := c.BindJSON(&param); err != nil {
			logger.Infof("failed to BindJSON. err: %v", err)
			return nil
		}

		parameter, err := converter.ToProblemUpdateParameter(&param)
		if err != nil {
			return xerrors.Errorf("failed to ToProblemUpdateParameter. param: %+v, err: %w", parameter, err)
		}

		if err := h.studentUsecaseProblem.UpdateProblem(ctx, organizationID, operatorID, id, parameter); err != nil {
			return xerrors.Errorf("failed to UpdateProblem. param: %+v, err: %w", parameter, err)
		}

		c.Status(http.StatusOK)
		return nil
	}, h.errorHandle)
}

func (h *problemHandler) RemoveProblem(c *gin.Context) {
	ctx := c.Request.Context()
	logger := log.FromContext(ctx)

	controllerhelper.HandleSecuredFunction(c, func(organizationID userD.OrganizationID, operatorID userD.AppUserID) error {
		logger.Infof("RemvoeProblem. organizationID: %d, operatorID: %d", organizationID, operatorID)

		id, err := h.toProblemSelectParameter2(c)
		if err != nil {
			c.Status(http.StatusBadRequest)
			return nil
		}

		if err := h.studentUsecaseProblem.RemoveProblem(ctx, organizationID, operatorID, id); err != nil {
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

	controllerhelper.HandleSecuredFunction(c, func(organizationID userD.OrganizationID, operatorID userD.AppUserID) error {
		workbookID, err := ginhelper.GetUintFromPath(c, "workbookID")
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

		newIterator := func(workbookID domain.WorkbookID, problemType string) (service.ProblemAddParameterIterator, error) {
			return h.newIterator(ctx, workbookID, problemType, multipartFile)
		}

		if err := h.studentUsecaseProblem.ImportProblems(ctx, organizationID, operatorID, domain.WorkbookID(workbookID), newIterator); err != nil {
			return xerrors.Errorf("failed to ImportProblems. err: %w", err)
		}

		c.Status(http.StatusOK)
		return nil
	}, h.errorHandle)
}

func (h *problemHandler) toProblemSelectParameter1(c *gin.Context) (service.ProblemSelectParameter1, error) {
	workbookID, err := ginhelper.GetUintFromPath(c, "workbookID")
	if err != nil {
		return nil, libD.ErrInvalidArgument
	}
	problemID, err := ginhelper.GetUintFromPath(c, "problemID")
	if err != nil {
		return nil, libD.ErrInvalidArgument
	}
	param, err := service.NewProblemSelectParameter1(domain.WorkbookID(workbookID), domain.ProblemID(problemID))
	if err != nil {
		return nil, libD.ErrInvalidArgument
	}

	return param, nil
}

func (h *problemHandler) toProblemSelectParameter2(c *gin.Context) (service.ProblemSelectParameter2, error) {
	workbookID, err := ginhelper.GetUintFromPath(c, "workbookID")
	if err != nil {
		return nil, libD.ErrInvalidArgument
	}
	problemID, err := ginhelper.GetUintFromPath(c, "problemID")
	if err != nil {
		return nil, libD.ErrInvalidArgument
	}
	version, err := ginhelper.GetIntFromQuery(c, "version")
	if err != nil {
		return nil, libD.ErrInvalidArgument
	}
	param, err := service.NewProblemSelectParameter2(domain.WorkbookID(workbookID), domain.ProblemID(problemID), version)
	if err != nil {
		return nil, libD.ErrInvalidArgument
	}

	return param, nil
}

func (h *problemHandler) errorHandle(c *gin.Context, err error) bool {
	ctx := c.Request.Context()
	logger := log.FromContext(ctx)
	var pluginError = &domain.PluginError{}
	if errors.Is(err, service.ErrProblemAlreadyExists) {
		c.JSON(http.StatusConflict, gin.H{"message": "Problem already exists"})
		return true
	} else if errors.Is(err, service.ErrWorkbookNotFound) {
		c.JSON(http.StatusNotFound, gin.H{"message": err.Error()})
		return true
	} else if errors.Is(err, service.ErrProblemNotFound) {
		c.JSON(http.StatusNotFound, gin.H{"message": err.Error()})
		return true
	} else if errors.As(err, &pluginError) {
		h := gin.H{
			"code":     pluginError.ErrorCode,
			"message":  pluginError.Error(),
			"messages": pluginError.ErrorMessages,
		}
		switch strings.ToLower(string(pluginError.ErrorType)) {
		case "client":
			c.JSON(http.StatusBadRequest, h)
		case "server":
			c.JSON(http.StatusInternalServerError, h)
		default:
			c.JSON(http.StatusInternalServerError, h)
		}
		return true
	}
	logger.Errorf("problemHandler error: %+v", err)
	return false
}
