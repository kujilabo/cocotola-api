package ginhelper

import (
	"strconv"

	"github.com/gin-gonic/gin"
)

func GetUintFromPath(c *gin.Context, param string) (uint, error) {
	idS := c.Param(param)
	id, err := strconv.Atoi(idS)
	if err != nil {
		return 0, err
	}

	return uint(id), nil
}

func GetIntFromPath(c *gin.Context, param string) (int, error) {
	idS := c.Param(param)
	id, err := strconv.Atoi(idS)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func GetStringFromPath(c *gin.Context, param string) string {
	return c.Param(param)
}

func GetIntFromQuery(c *gin.Context, param string) (int, error) {
	idS := c.Query(param)
	id, err := strconv.Atoi(idS)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func GetStringFromQuery(c *gin.Context, param string) string {
	return c.Query(param)
}

func GetIntFromForm(c *gin.Context, param string) (int, error) {
	idS := c.PostForm(param)
	id, err := strconv.Atoi(idS)
	if err != nil {
		return 0, err
	}

	return id, nil
}
