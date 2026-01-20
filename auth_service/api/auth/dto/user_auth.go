package dto

import (
	"github.com/afteracademy/gomicro/auth-service/api/user/dto"
	"github.com/afteracademy/gomicro/auth-service/api/user/model"
	"github.com/afteracademy/goserve/v2/utility"
	"github.com/go-playground/validator/v10"
)

type UserAuth struct {
	User   *dto.InfoPrivateUser `json:"user" validate:"required"`
	Tokens *UserTokens          `json:"tokens" validate:"required"`
}

func NewUserAuth(user *model.User, tokens *UserTokens) *UserAuth {
	return &UserAuth{
		User:   dto.NewInfoPrivateUser(user),
		Tokens: tokens,
	}
}

func (d *UserAuth) GetValue() *UserAuth {
	return d
}

func (d *UserAuth) ValidateErrors(errs validator.ValidationErrors) ([]string, error) {
	return utility.FormatValidationErrors(errs), nil
}
