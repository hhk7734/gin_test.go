package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func GinHealthzController(c *gin.Context) {
	c.Status(http.StatusOK)
}
