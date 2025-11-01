package handlers

import (
	"strconv"

	"github.com/gin-gonic/gin"
)

func getIntFromPath(c *gin.Context, variableName string) (int32, error) {
	id := c.Param(variableName)
	idInt32, err := strconv.ParseInt(id, 10, 32)
	return int32(idInt32), err
}
