package middleware

import (
	"github.com/afteracademy/gomicro/blog-service/api/auth"
	"github.com/afteracademy/gomicro/blog-service/common"
	"github.com/afteracademy/goserve/v2/network"
	"github.com/gin-gonic/gin"
)

type authenticationProvider struct {
	common.ContextPayload
	authService auth.Service
}

func NewAuthenticationProvider(authService auth.Service) network.AuthenticationProvider {
	return &authenticationProvider{
		ContextPayload: common.NewContextPayload(),
		authService:    authService,
	}
}

func (m *authenticationProvider) Middleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authHeader := ctx.GetHeader(network.AuthorizationHeader)

		user, err := m.authService.Authenticate(authHeader)
		if err != nil {
			network.SendUnauthorizedError(ctx, err.Error(), err)
			return
		}

		m.SetUser(ctx, user)
		ctx.Next()
	}
}
