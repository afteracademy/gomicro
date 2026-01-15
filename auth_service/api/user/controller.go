package user

import (
	"github.com/gin-gonic/gin"
	"github.com/afteracademy/gomicro/auth-service/api/auth/message"
	"github.com/afteracademy/gomicro/auth-service/common"
	coredto "github.com/afteracademy/goserve/arch/dto"
	"github.com/afteracademy/goserve/arch/micro"
	"github.com/afteracademy/goserve/arch/network"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type controller struct {
	micro.BaseController
	common.ContextPayload
	service Service
}

func NewController(
	authProvider network.AuthenticationProvider,
	authorizeProvider network.AuthorizationProvider,
	service Service,
) micro.Controller {
	return &controller{
		BaseController: micro.NewBaseController("/profile", authProvider, authorizeProvider),
		ContextPayload: common.NewContextPayload(),
		service:        service,
	}
}

func (c *controller) MountNats(group micro.NatsGroup) {
	group.AddEndpoint("user", micro.NatsHandlerFunc(c.userHandler))
}

func (c *controller) userHandler(req micro.NatsRequest) {
	text, err := micro.ParseMsg[message.Text](req.Data())
	if err != nil {
		c.SendNats(req).Error(err)
		return
	}

	userId, err := primitive.ObjectIDFromHex(text.Value)
	if err != nil {
		c.SendNats(req).Error(err)
		return
	}

	user, err := c.service.FindUserPublicProfile(userId)
	if err != nil {
		c.SendNats(req).Error(err)
		return
	}

	c.SendNats(req).Message(message.NewUser(user))
}

func (c *controller) MountRoutes(group *gin.RouterGroup) {
	group.GET("/id/:id", c.getPublicProfileHandler)
	private := group.Use(c.Authentication())
	private.GET("/mine", c.getPrivateProfileHandler)
}

func (c *controller) getPublicProfileHandler(ctx *gin.Context) {
	mongoId, err := network.ReqParams(ctx, coredto.EmptyMongoId())
	if err != nil {
		c.Send(ctx).BadRequestError(err.Error(), err)
		return
	}

	data, err := c.service.GetUserPublicProfile(mongoId.ID)
	if err != nil {
		c.Send(ctx).MixedError(err)
		return
	}

	c.Send(ctx).SuccessDataResponse("success", data)
}

func (c *controller) getPrivateProfileHandler(ctx *gin.Context) {
	user := c.MustGetUser(ctx)

	data, err := c.service.GetUserPrivateProfile(user)
	if err != nil {
		c.Send(ctx).MixedError(err)
		return
	}

	c.Send(ctx).SuccessDataResponse("success", data)
}
