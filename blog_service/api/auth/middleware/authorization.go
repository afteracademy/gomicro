package middleware

import (
	"github.com/afteracademy/gomicro/blog-service/api/auth"
	"github.com/afteracademy/gomicro/blog-service/common"
	"github.com/afteracademy/goserve/v2/network"
	"github.com/gin-gonic/gin"
)

type authorizationProvider struct {
	common.ContextPayload
	authService auth.Service
}

func NewAuthorizationProvider(authService auth.Service) network.AuthorizationProvider {
	return &authorizationProvider{
		ContextPayload: common.NewContextPayload(),
		authService:    authService,
	}
}

func (m *authorizationProvider) Middleware(roleNames ...string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		user := m.MustGetUser(ctx)

		err := m.authService.Authorize(user, roleNames...)
		if err != nil {
			network.SendForbiddenError(ctx, err.Error(), err)
			return
		}

		ctx.Next()
	}
}
