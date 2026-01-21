package dto

import (
	"github.com/afteracademy/gomicro/auth-service/api/user/model"
	"github.com/google/uuid"
)

type UserPublic struct {
	ID            uuid.UUID `json:"id" binding:"required" validate:"required"`
	Name          string    `json:"name" binding:"required" validate:"required"`
	ProfilePicURL *string   `json:"profilePicUrl,omitempty" validate:"omitempty,url"`
}

func NewUserPublic(user *model.User) *UserPublic {
	return &UserPublic{
		ID:            user.ID,
		Name:          user.Name,
		ProfilePicURL: user.ProfilePicURL,
	}
}
