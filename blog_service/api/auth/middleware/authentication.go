package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/afteracademy/gomicro/blog-service/api/auth"
	"github.com/afteracademy/gomicro/blog-service/common"
	"github.com/afteracademy/goserve/arch/network"
)

type authenticationProvider struct {
	network.ResponseSender
	common.ContextPayload
	authService auth.Service
}

func NewAuthenticationProvider(authService auth.Service) network.AuthenticationProvider {
	return &authenticationProvider{
		ResponseSender: network.NewResponseSender(),
		ContextPayload: common.NewContextPayload(),
		authService:    authService,
	}
}

func (m *authenticationProvider) Middleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authHeader := ctx.GetHeader(network.AuthorizationHeader)

		user, err := m.authService.Authenticate(authHeader)
		if err != nil {
			m.Send(ctx).UnauthorizedError(err.Error(), err)
			return
		}

		m.SetUser(ctx, user)
		ctx.Next()
	}
}
