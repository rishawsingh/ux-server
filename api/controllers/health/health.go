package health

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/remotestate/golang/internal"
)

type Controller struct {
	logger *internal.Logger
	env    *internal.EnvConfig
}

func NewController(logger *internal.Logger, env *internal.EnvConfig) *Controller {
	return &Controller{logger: logger, env: env}
}

func (h *Controller) Hello(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status": "server is running",
		"build":  h.env.BuildNumber(),
	})
}
