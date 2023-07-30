package models

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetTrans(c *gin.Context) {
	txHash := c.Param("txHash")

	c.JSON(http.StatusOK, txHash)
}
