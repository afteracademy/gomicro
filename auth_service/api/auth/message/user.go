package message

import (
	"github.com/afteracademy/gomicro/auth-service/api/user/model"
	"github.com/google/uuid"
)

type User struct {
	ID            uuid.UUID `json:"id" validate:"required,uuid"`
	Name          string    `json:"name" validate:"required"`
	Email         string    `json:"email" validate:"required,email"`
	ProfilePicURL *string   `json:"profilePicUrl,omitempty" validate:"omitempty,url"`
}

func NewUser(user *model.User) *User {
	return &User{
		ID:            user.ID,
		Name:          user.Name,
		Email:         user.Email,
		ProfilePicURL: user.ProfilePicURL,
	}
}
