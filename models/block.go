package models

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetBlockByID(c *gin.Context) {
	blockID := c.Param("id")
	c.JSON(http.StatusOK, blockID)
}
