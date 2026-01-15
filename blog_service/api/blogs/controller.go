package blogs

import (
	"github.com/gin-gonic/gin"
	"github.com/afteracademy/gomicro/blog-service/api/blogs/dto"
	coredto "github.com/afteracademy/goserve/arch/dto"
	"github.com/afteracademy/goserve/arch/micro"
	"github.com/afteracademy/goserve/arch/network"
)

type controller struct {
	micro.BaseController
	service Service
}

func NewController(
	authMFunc network.AuthenticationProvider,
	authorizeMFunc network.AuthorizationProvider,
	service Service,
) micro.Controller {
	return &controller{
		BaseController: micro.NewBaseController("/list", authMFunc, authorizeMFunc),
		service:        service,
	}
}

func (c *controller) MountNats(group micro.NatsGroup) {}

func (c *controller) MountRoutes(group *gin.RouterGroup) {
	group.GET("/latest", c.getLatestBlogsHandler)
	group.GET("/tag/:tag", c.getTaggedBlogsHandler)
	group.GET("/similar/id/:id", c.getSimilarBlogsHandler)
}

func (c *controller) getLatestBlogsHandler(ctx *gin.Context) {
	pagination, err := network.ReqQuery(ctx, coredto.EmptyPagination())
	if err != nil {
		c.Send(ctx).BadRequestError(err.Error(), err)
		return
	}

	blogs, err := c.service.GetPaginatedLatestBlogs(pagination)
	if err != nil {
		c.Send(ctx).MixedError(err)
		return
	}

	c.Send(ctx).SuccessDataResponse("success", blogs)
}

func (c *controller) getTaggedBlogsHandler(ctx *gin.Context) {
	tag, err := network.ReqParams(ctx, dto.EmptyTag())
	if err != nil {
		c.Send(ctx).BadRequestError(err.Error(), err)
		return
	}

	pagination, err := network.ReqQuery(ctx, coredto.EmptyPagination())
	if err != nil {
		c.Send(ctx).BadRequestError(err.Error(), err)
		return
	}

	blogs, err := c.service.GetPaginatedTaggedBlogs(tag.Tag, pagination)
	if err != nil {
		c.Send(ctx).MixedError(err)
		return
	}

	c.Send(ctx).SuccessDataResponse("success", blogs)
}

func (c *controller) getSimilarBlogsHandler(ctx *gin.Context) {
	mongoId, err := network.ReqParams(ctx, coredto.EmptyMongoId())
	if err != nil {
		c.Send(ctx).BadRequestError(err.Error(), err)
		return
	}

	blogs, err := c.service.GetSimilarBlogsDtoCache(mongoId.ID)
	if err == nil {
		c.Send(ctx).SuccessDataResponse("success", blogs)
		return
	}

	blogs, err = c.service.GetSimilarBlogs(mongoId.ID)
	if err != nil {
		c.Send(ctx).MixedError(err)
		return
	}

	c.Send(ctx).SuccessDataResponse("success", blogs)
	c.service.SetSimilarBlogsDtoCache(mongoId.ID, blogs)
}
