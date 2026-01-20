package dto

import (
	"github.com/afteracademy/gomicro/auth-service/api/user/model"
	"github.com/afteracademy/goserve/v2/utility"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type InfoPrivateUser struct {
	ID            primitive.ObjectID `json:"_id" binding:"required" validate:"required"`
	Email         string             `json:"email" binding:"required" validate:"required,email"`
	Name          string             `json:"name" binding:"required" validate:"required"`
	ProfilePicURL *string            `json:"profilePicUrl,omitempty" validate:"omitempty,url"`
	Roles         []*InfoRole        `json:"roles" validate:"required,dive,required"`
}

func NewInfoPrivateUser(user *model.User) *InfoPrivateUser {
	roles := make([]*InfoRole, len(user.Roles))
	for i, role := range user.RoleDocs {
		roles[i] = NewInfoRole(role)
	}

	return &InfoPrivateUser{
		ID:            user.ID,
		Email:         user.Email,
		Name:          user.Name,
		ProfilePicURL: user.ProfilePicURL,
		Roles:         roles,
	}
}

func (d *InfoPrivateUser) GetValue() *InfoPrivateUser {
	return d
}

func (d *InfoPrivateUser) ValidateErrors(errs validator.ValidationErrors) ([]string, error) {
	return utility.FormatValidationErrors(errs), nil
}
