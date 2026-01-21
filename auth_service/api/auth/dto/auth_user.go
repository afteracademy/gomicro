package dto

import (
	userdto "github.com/afteracademy/gomicro/auth-service/api/user/dto"
	"github.com/afteracademy/gomicro/auth-service/api/user/model"
)

type UserAuth struct {
	User   *userdto.UserPrivate `json:"user" validate:"required"`
	Tokens *Tokens              `json:"tokens" validate:"required"`
}

func NewUserAuth(user *model.User, tokens *Tokens) *UserAuth {
	return &UserAuth{
		User:   userdto.NewUserPrivate(user),
		Tokens: tokens,
	}
}
