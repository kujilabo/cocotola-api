package handler

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/kujilabo/cocotola-api/pkg_app/domain"
	"github.com/kujilabo/cocotola-api/pkg_app/handler/converter"
	"github.com/kujilabo/cocotola-api/pkg_app/service"
	studentU "github.com/kujilabo/cocotola-api/pkg_app/usecase/student"
	"github.com/kujilabo/cocotola-api/pkg_lib/ginhelper"
	"github.com/kujilabo/cocotola-api/pkg_lib/log"
	userD "github.com/kujilabo/cocotola-api/pkg_user/domain"
	"github.com/kujilabo/cocotola-api/pkg_user/handlerhelper"
)

type AudioHandler interface {
	FindAudioByID(c *gin.Context)
}

type audioHandler struct {
	studentUsecaseAudio studentU.StudentUsecaseAudio
}

func NewAudioHandler(studentUsecaseAudio studentU.StudentUsecaseAudio) AudioHandler {
	return &audioHandler{
		studentUsecaseAudio: studentUsecaseAudio,
	}
}

func (h *audioHandler) FindAudioByID(c *gin.Context) {
	ctx := c.Request.Context()
	logger := log.FromContext(ctx)
	logger.Info("FindAudioByID")

	handlerhelper.HandleSecuredFunction(c, func(organizationID userD.OrganizationID, operatorID userD.AppUserID) error {
		workbookID, err := ginhelper.GetUintFromPath(c, "workbookID")
		if err != nil {
			c.Status(http.StatusBadRequest)
			return nil
		}
		problemID, err := ginhelper.GetUintFromPath(c, "problemID")
		if err != nil {
			c.Status(http.StatusBadRequest)
			return nil
		}
		id := c.Param("audioID")
		audioID, err := strconv.Atoi(id)
		if err != nil {
			c.Status(http.StatusBadRequest)
			return nil
		}

		result, err := h.studentUsecaseAudio.FindAudioByID(ctx, organizationID, operatorID, domain.WorkbookID(workbookID), domain.ProblemID(problemID), domain.AudioID(uint(audioID)))
		if err != nil {
			return err
		}

		response, err := converter.ToAudioResponse(ctx, result)
		if err != nil {
			return err
		}

		c.JSON(http.StatusOK, response)
		return nil
	}, h.errorHandle)
}

func (h *audioHandler) errorHandle(c *gin.Context, err error) bool {
	ctx := c.Request.Context()
	logger := log.FromContext(ctx)
	if errors.Is(err, service.ErrAudioNotFound) {
		logger.Warnf("audioHandler err: %+v", err)
		c.JSON(http.StatusNotFound, gin.H{"message": "Audio not found"})
		return true
	}
	logger.Errorf("error:%v", err)
	return false
}
