package dto

import (
	"github.com/afteracademy/gomicro/blog-service/api/auth/message"
	"github.com/afteracademy/goserve/v2/utility"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

type InfoAuthor struct {
	ID            uuid.UUID `json:"_id" binding:"required" validate:"required"`
	Name          string    `json:"name" binding:"required" validate:"required"`
	ProfilePicURL *string   `json:"profilePicUrl,omitempty" validate:"omitempty,url"`
}

func NewInfoPrivateUser(user *message.User) *InfoAuthor {
	return &InfoAuthor{
		ID:            user.ID,
		Name:          user.Name,
		ProfilePicURL: user.ProfilePicURL,
	}
}

func (d *InfoAuthor) GetValue() *InfoAuthor {
	return d
}

func (d *InfoAuthor) ValidateErrors(errs validator.ValidationErrors) ([]string, error) {
	return utility.FormatValidationErrors(errs), nil
}
