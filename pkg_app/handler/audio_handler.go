package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"golang.org/x/xerrors"

	"github.com/kujilabo/cocotola-api/pkg_app/application"
	"github.com/kujilabo/cocotola-api/pkg_app/domain"
	"github.com/kujilabo/cocotola-api/pkg_lib/handlerhelper"
	"github.com/kujilabo/cocotola-api/pkg_lib/log"
	user "github.com/kujilabo/cocotola-api/pkg_user/domain"
)

type AudioHandler interface {
	FindAudioByID(c *gin.Context)
}

type audioHandler struct {
	audioService application.AudioService
}

func NewAudioHandler(audioService application.AudioService) AudioHandler {
	return &audioHandler{
		audioService: audioService,
	}
}

func (h *audioHandler) FindAudioByID(c *gin.Context) {
	ctx := c.Request.Context()
	logger := log.FromContext(ctx)
	logger.Info("FindAudioByID")

	handlerhelper.HandleSecuredFunction(c, func(organizationID user.OrganizationID, operatorID user.AppUserID) error {
		id := c.Param("audioID")
		audioID, err := strconv.Atoi(id)
		if err != nil {
			c.Status(http.StatusBadRequest)
			return nil
		}

		audio, err := h.audioService.FindAudioByID(ctx, domain.AudioID(uint(audioID)))
		if err != nil {
			return err
		}

		response := map[string]string{
			"audioContent": audio.GetAudioContent(),
		}

		c.JSON(http.StatusOK, response)
		return nil
	}, h.errorHandle)

}

func (h *audioHandler) errorHandle(c *gin.Context, err error) bool {
	ctx := c.Request.Context()
	logger := log.FromContext(ctx)
	if xerrors.Is(err, domain.ErrAudioNotFound) {
		logger.Warnf("audioHandler err: %+v", err)
		c.JSON(http.StatusNotFound, gin.H{"message": "Audio not found"})
		return true
	}
	logger.Errorf("error:%v", err)
	return false
}
