package health

import (
	"github.com/afteracademy/goserve/v2/micro"
	"github.com/afteracademy/goserve/v2/network"
	"github.com/gin-gonic/gin"
)

type controller struct {
	micro.Controller
	service Service
}

func NewController(
	service Service,
) micro.Controller {
	return &controller{
		Controller: micro.NewController("/health", nil, nil),
		service:    service,
	}
}

func (c *controller) MountRoutes(group *gin.RouterGroup) {
	group.GET("", c.getHealthHandler)
}

func (c *controller) getHealthHandler(ctx *gin.Context) {
	health, err := c.service.CheckHealth()
	if err != nil {
		network.SendMixedError(ctx, err)
		return
	}

	network.SendSuccessDataResponse(ctx, "success", health)
}
