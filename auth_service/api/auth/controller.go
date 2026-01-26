package auth

import (
	"github.com/afteracademy/gomicro/auth-service/api/auth/dto"
	"github.com/afteracademy/gomicro/auth-service/api/auth/message"
	"github.com/afteracademy/gomicro/auth-service/api/user"
	"github.com/afteracademy/gomicro/auth-service/common"
	"github.com/afteracademy/gomicro/auth-service/utils"
	"github.com/afteracademy/goserve/v2/micro"
	"github.com/afteracademy/goserve/v2/network"
	"github.com/gin-gonic/gin"
)

type controller struct {
	micro.Controller
	common.ContextPayload
	service     Service
	userService user.Service
}

func NewController(
	authProvider network.AuthenticationProvider,
	authorizeProvider network.AuthorizationProvider,
	service Service,
	userService user.Service,
) micro.Controller {
	return &controller{
		Controller:     micro.NewController("/", authProvider, authorizeProvider),
		ContextPayload: common.NewContextPayload(),
		service:        service,
		userService:    userService,
	}
}

func (c *controller) MountNats(group micro.NatsGroup) {
	group.AddEndpoint("authentication", micro.NatsHandlerFunc(c.authenticationHandler))
	group.AddEndpoint("authorization", micro.NatsHandlerFunc(c.authorizationHandler))
}

func (c *controller) authenticationHandler(req micro.NatsRequest) {
	text, err := micro.JsonToMsg[message.Text](req.Data())
	if err != nil {
		micro.RespondNatsError(req, err)
		return
	}

	user, _, err := c.service.Authenticate(text.Value)
	if err != nil {
		micro.RespondNatsError(req, err)
		return
	}

	micro.RespondNatsMessage(req, message.NewUser(user))
}

func (c *controller) authorizationHandler(req micro.NatsRequest) {
	userRole, err := micro.JsonToMsg[message.UserRole](req.Data())
	if err != nil {
		micro.RespondNatsError(req, err)
		return
	}

	user, err := c.userService.FetchUserById(userRole.User.ID)
	if err != nil {
		micro.RespondNatsError(req, err)
		return
	}

	err = c.service.Authorize(user, userRole.Roles...)
	if err != nil {
		micro.RespondNatsError(req, err)
		return
	}

	micro.RespondNatsMessage(req, message.NewUser(user))
}

func (c *controller) MountRoutes(group *gin.RouterGroup) {
	group.GET("/verify/apikey", c.verifyApikeyHandler)
	group.POST("/signup/basic", c.signUpBasicHandler)
	group.POST("/signin/basic", c.signInBasicHandler)
	group.POST("/token/refresh", c.tokenRefreshHandler)
	group.DELETE("/signout", c.Authentication(), c.signOutBasic)
}

func (c *controller) verifyApikeyHandler(ctx *gin.Context) {
	key := ctx.GetHeader(network.ApiKeyHeader)
	if len(key) == 0 {
		network.SendUnauthorizedError(ctx, "permission denied: missing x-api-key header", nil)
		return
	}

	_, err := c.service.FetchApiKey(key)
	if err != nil {
		network.SendForbiddenError(ctx, "permission denied: invalid x-api-key", err)
		return
	}

	network.SendSuccessMsgResponse(ctx, "success")
}

func (c *controller) signUpBasicHandler(ctx *gin.Context) {
	body, err := network.ReqBody[dto.SignUpBasic](ctx)
	if err != nil {
		network.SendBadRequestError(ctx, err.Error(), err)
		return
	}

	data, err := c.service.SignUpBasic(body)
	if err != nil {
		network.SendMixedError(ctx, err)
		return
	}

	network.SendSuccessDataResponse(ctx, "success", data)
}

func (c *controller) signInBasicHandler(ctx *gin.Context) {
	body, err := network.ReqBody[dto.SignInBasic](ctx)
	if err != nil {
		network.SendBadRequestError(ctx, err.Error(), err)
		return
	}

	dto, err := c.service.SignInBasic(body)
	if err != nil {
		network.SendMixedError(ctx, err)
		return
	}

	network.SendSuccessDataResponse(ctx, "success", dto)
}

func (c *controller) signOutBasic(ctx *gin.Context) {
	keystore := c.MustGetKeystore(ctx)

	err := c.service.SignOut(keystore)
	if err != nil {
		network.SendInternalServerError(ctx, "something went wrong", err)
		return
	}

	network.SendSuccessMsgResponse(ctx, "signout success")
}

func (c *controller) tokenRefreshHandler(ctx *gin.Context) {
	body, err := network.ReqBody[dto.TokenRefresh](ctx)
	if err != nil {
		network.SendBadRequestError(ctx, err.Error(), err)
		return
	}

	authHeader := ctx.GetHeader(network.AuthorizationHeader)
	accessToken := utils.ExtractBearerToken(authHeader)

	dto, err := c.service.RenewToken(body, accessToken)
	if err != nil {
		network.SendMixedError(ctx, err)
		return
	}

	network.SendSuccessDataResponse(ctx, "success", dto)
}
