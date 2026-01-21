package editor

import (
	"github.com/afteracademy/gomicro/blog-service/api/auth/message"
	"github.com/afteracademy/gomicro/blog-service/common"
	coredto "github.com/afteracademy/goserve/v2/dto"
	"github.com/afteracademy/goserve/v2/micro"
	"github.com/afteracademy/goserve/v2/network"
	"github.com/gin-gonic/gin"
)

type controller struct {
	micro.Controller
	common.ContextPayload
	service Service
}

func NewController(
	authMFunc network.AuthenticationProvider,
	authorizeMFunc network.AuthorizationProvider,
	service Service,
) micro.Controller {
	return &controller{
		Controller:     micro.NewController("/editor", authMFunc, authorizeMFunc),
		ContextPayload: common.NewContextPayload(),
		service:        service,
	}
}

func (c *controller) MountNats(group micro.NatsGroup) {}

func (c *controller) MountRoutes(group *gin.RouterGroup) {
	group.Use(c.Authentication(), c.Authorization(string(message.RoleCodeEditor)))
	group.GET("/id/:id", c.getBlogHandler)
	group.PUT("/publish/id/:id", c.publishBlogHandler)
	group.PUT("/unpublish/id/:id", c.unpublishBlogHandler)
	group.GET("/submitted", c.getSubmittedBlogsHandler)
	group.GET("/published", c.getPublishedBlogsHandler)
}

func (c *controller) getBlogHandler(ctx *gin.Context) {
	mongoId, err := network.ReqParams[coredto.MongoId](ctx)
	if err != nil {
		network.SendBadRequestError(ctx, err.Error(), err)
		return
	}

	blog, err := c.service.GetBlogById(mongoId.ID)
	if err != nil {
		network.SendNotFoundError(ctx, mongoId.ID.Hex()+" not found", err)
		return
	}

	network.SendSuccessDataResponse(ctx, "success", blog)
}

func (c *controller) publishBlogHandler(ctx *gin.Context) {
	mongoId, err := network.ReqParams[coredto.MongoId](ctx)
	if err != nil {
		network.SendBadRequestError(ctx, err.Error(), err)
		return
	}

	user := c.MustGetUser(ctx)

	err = c.service.BlogPublication(mongoId.ID, user, true)
	if err != nil {
		network.SendMixedError(ctx, err)
		return
	}

	network.SendSuccessMsgResponse(ctx, "blog published successfully")
}

func (c *controller) unpublishBlogHandler(ctx *gin.Context) {
	mongoId, err := network.ReqParams[coredto.MongoId](ctx)
	if err != nil {
		network.SendBadRequestError(ctx, err.Error(), err)
		return
	}

	user := c.MustGetUser(ctx)

	err = c.service.BlogPublication(mongoId.ID, user, false)
	if err != nil {
		network.SendMixedError(ctx, err)
		return
	}

	network.SendSuccessMsgResponse(ctx, "blog unpublished successfully")
}

func (c *controller) getSubmittedBlogsHandler(ctx *gin.Context) {
	pagination, err := network.ReqQuery[coredto.Pagination](ctx)
	if err != nil {
		network.SendBadRequestError(ctx, err.Error(), err)
		return
	}

	blogs, err := c.service.GetPaginatedSubmitted(pagination)
	if err != nil {
		network.SendMixedError(ctx, err)
		return
	}

	network.SendSuccessDataResponse(ctx, "success", &blogs)
}

func (c *controller) getPublishedBlogsHandler(ctx *gin.Context) {
	pagination, err := network.ReqQuery[coredto.Pagination](ctx)
	if err != nil {
		network.SendBadRequestError(ctx, err.Error(), err)
		return
	}

	blogs, err := c.service.GetPaginatedPublished(pagination)
	if err != nil {
		network.SendMixedError(ctx, err)
		return
	}

	network.SendSuccessDataResponse(ctx, "success", &blogs)
}
