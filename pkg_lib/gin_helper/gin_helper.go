package gin_helper

import (
	"strconv"

	"github.com/gin-gonic/gin"
)

func GetUint(c *gin.Context, param string) (uint, error) {
	idS := c.Param(param)
	id, err := strconv.Atoi(idS)
	if err != nil {
		return 0, err
	}
	return uint(id), nil
}

func GetInt(c *gin.Context, param string) (int, error) {
	idS := c.Param(param)
	id, err := strconv.Atoi(idS)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func GetString(c *gin.Context, param string) string {
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
