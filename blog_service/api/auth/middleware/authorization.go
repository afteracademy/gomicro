package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/afteracademy/gomicro/blog-service/api/auth"
	"github.com/afteracademy/gomicro/blog-service/common"
	"github.com/afteracademy/goserve/v2/network"
)

type authorizationProvider struct {
	network.ResponseSender
	common.ContextPayload
	authService auth.Service
}

func NewAuthorizationProvider(authService auth.Service) network.AuthorizationProvider {
	return &authorizationProvider{
		ResponseSender: network.NewResponseSender(),
		ContextPayload: common.NewContextPayload(),
		authService:    authService,
	}
}

func (m *authorizationProvider) Middleware(roleNames ...string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		user := m.MustGetUser(ctx)

		err := m.authService.Authorize(user, roleNames...)
		if err != nil {
			m.Send(ctx).ForbiddenError(err.Error(), err)
			return
		}

		ctx.Next()
	}
}
