package blog

import (
	coredto "github.com/afteracademy/goserve/v2/dto"
	"github.com/afteracademy/goserve/v2/micro"
	"github.com/afteracademy/goserve/v2/network"
	"github.com/gin-gonic/gin"
)

type controller struct {
	micro.Controller
	service Service
}

func NewController(
	authMFunc network.AuthenticationProvider,
	authorizeMFunc network.AuthorizationProvider,
	service Service,
) micro.Controller {
	return &controller{
		Controller: micro.NewController("/", authMFunc, authorizeMFunc),
		service:    service,
	}
}

func (c *controller) MountNats(group micro.NatsGroup) {}

func (c *controller) MountRoutes(group *gin.RouterGroup) {
	group.GET("/id/:id", c.getBlogByIdHandler)
	group.GET("/slug/:slug", c.getBlogBySlugHandler)
}

func (c *controller) getBlogByIdHandler(ctx *gin.Context) {
	mongoId, err := network.ReqParams[coredto.MongoId](ctx)
	if err != nil {
		network.SendBadRequestError(ctx, err.Error(), err)
		return
	}

	blog, err := c.service.GetBlogDtoCacheById(mongoId.ID)
	if err == nil {
		network.SendSuccessDataResponse(ctx, "success", blog)
		return
	}

	blog, err = c.service.GetPublisedBlogById(mongoId.ID)
	if err != nil {
		network.SendMixedError(ctx, err)
		return
	}

	network.SendSuccessDataResponse(ctx, "success", blog)
	c.service.SetBlogDtoCacheById(blog)
}

func (c *controller) getBlogBySlugHandler(ctx *gin.Context) {
	slug, err := network.ReqParams[coredto.Slug](ctx)
	if err != nil {
		network.SendBadRequestError(ctx, err.Error(), err)
		return
	}

	blog, err := c.service.GetBlogDtoCacheBySlug(slug.Slug)
	if err == nil {
		network.SendSuccessDataResponse(ctx, "success", blog)
		return
	}

	blog, err = c.service.GetPublishedBlogBySlug(slug.Slug)
	if err != nil {
		network.SendMixedError(ctx, err)
		return
	}

	network.SendSuccessDataResponse(ctx, "success", blog)
	c.service.SetBlogDtoCacheBySlug(blog)
}
