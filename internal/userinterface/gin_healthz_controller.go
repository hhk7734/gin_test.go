package userinterface

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type GinHealthzController struct {
}

func (h *GinHealthzController) Healthz(c *gin.Context) {
	c.Status(http.StatusOK)
}
