package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/afteracademy/gomicro/auth-service/api/auth"
	"github.com/afteracademy/gomicro/auth-service/api/user"
	"github.com/afteracademy/gomicro/auth-service/common"
	"github.com/afteracademy/goserve/v2/network"
)

type authenticationProvider struct {
	common.ContextPayload
	authService auth.Service
	userService user.Service
}

func NewAuthenticationProvider(authService auth.Service, userService user.Service) network.AuthenticationProvider {
	return &authenticationProvider{
		ContextPayload: common.NewContextPayload(),
		authService:    authService,
		userService:    userService,
	}
}

func (m *authenticationProvider) Middleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authHeader := ctx.GetHeader(network.AuthorizationHeader)

		user, keystore, err := m.authService.Authenticate(authHeader)
		if err != nil {
			network.SendMixedError(ctx, err)
			return
		}

		m.SetUser(ctx, user)
		m.SetKeystore(ctx, keystore)

		ctx.Next()
	}
}
