package models

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func LatestBlocks(c *gin.Context) {
	n := c.DefaultQuery("limit", "10")

	c.JSON(http.StatusOK, n)
}
