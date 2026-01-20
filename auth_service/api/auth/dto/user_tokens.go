package dto

import (
	"github.com/afteracademy/goserve/v2/utility"
	"github.com/go-playground/validator/v10"
)

type UserTokens struct {
	AccessToken  string `json:"accessToken" binding:"required" validate:"required"`
	RefreshToken string `json:"refreshToken" binding:"required" validate:"required"`
}

func NewUserTokens(access string, refresh string) *UserTokens {
	return &UserTokens{
		AccessToken:  access,
		RefreshToken: refresh,
	}
}

func (d *UserTokens) GetValue() *UserTokens {
	return d
}

func (d *UserTokens) ValidateErrors(errs validator.ValidationErrors) ([]string, error) {
	return utility.FormatValidationErrors(errs), nil
}
