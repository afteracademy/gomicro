package user

import (
	"github.com/afteracademy/gomicro/auth-service/api/auth/message"
	"github.com/afteracademy/gomicro/auth-service/common"
	coredto "github.com/afteracademy/goserve/v2/dto"
	"github.com/afteracademy/goserve/v2/micro"
	"github.com/afteracademy/goserve/v2/network"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type controller struct {
	micro.Controller
	common.ContextPayload
	service Service
}

func NewController(
	authProvider network.AuthenticationProvider,
	authorizeProvider network.AuthorizationProvider,
	service Service,
) micro.Controller {
	return &controller{
		Controller:     micro.NewController("/profile", authProvider, authorizeProvider),
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
		micro.SendNatsError(req, err)
		return
	}

	userId, err := primitive.ObjectIDFromHex(text.Value)
	if err != nil {
		micro.SendNatsError(req, err)
		return
	}

	user, err := c.service.FindUserPublicProfile(userId)
	if err != nil {
		micro.SendNatsError(req, err)
		return
	}

	micro.SendNatsMessage(req, message.NewUser(user))
}

func (c *controller) MountRoutes(group *gin.RouterGroup) {
	group.GET("/id/:id", c.getPublicProfileHandler)
	private := group.Use(c.Authentication())
	private.GET("/mine", c.getPrivateProfileHandler)
}

func (c *controller) getPublicProfileHandler(ctx *gin.Context) {
	mongoId, err := network.ReqParams[coredto.MongoId](ctx)
	if err != nil {
		network.SendBadRequestError(ctx, err.Error(), err)
		return
	}

	data, err := c.service.GetUserPublicProfile(mongoId.ID)
	if err != nil {
		network.SendMixedError(ctx, err)
		return
	}

	network.SendSuccessDataResponse(ctx, "success", data)
}

func (c *controller) getPrivateProfileHandler(ctx *gin.Context) {
	user := c.MustGetUser(ctx)

	data, err := c.service.GetUserPrivateProfile(user)
	if err != nil {
		network.SendMixedError(ctx, err)
		return
	}

	network.SendSuccessDataResponse(ctx, "success", data)
}
