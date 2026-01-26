package common

import (
	"errors"

	"github.com/afteracademy/gomicro/blog-service/api/auth/message"
	"github.com/gin-gonic/gin"
)

const (
	payloadUser string = "user"
)

type ContextPayload interface {
	SetUser(ctx *gin.Context, value *message.User)
	MustGetUser(ctx *gin.Context) *message.User
}

type payload struct{}

func NewContextPayload() ContextPayload {
	return &payload{}
}

func (u *payload) SetUser(ctx *gin.Context, value *message.User) {
	ctx.Set(payloadUser, value)
}

func (u *payload) MustGetUser(ctx *gin.Context) *message.User {
	value, ok := ctx.MustGet(payloadUser).(*message.User)
	if !ok {
		panic(errors.New(payloadUser + " missing for context"))
	}
	return value
}
