package dto

import (
	"github.com/afteracademy/gomicro/auth-service/api/user/model"
	"github.com/afteracademy/goserve/v2/utility"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type InfoRole struct {
	ID   primitive.ObjectID `json:"_id" binding:"required" validate:"required"`
	Code model.RoleCode     `json:"code" binding:"required" validate:"required,rolecode"`
}

func NewInfoRole(role *model.Role) *InfoRole {
	return &InfoRole{
		ID:   role.ID,
		Code: role.Code,
	}
}

func EmptyInfoRole() *InfoRole {
	return &InfoRole{}
}

func (d *InfoRole) GetValue() *InfoRole {
	return d
}

func (d *InfoRole) ValidateErrors(errs validator.ValidationErrors) ([]string, error) {
	return utility.FormatValidationErrors(errs), nil
}
