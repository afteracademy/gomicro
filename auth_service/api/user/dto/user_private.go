package dto

import (
	"github.com/afteracademy/gomicro/auth-service/api/user/model"
	"github.com/google/uuid"
)

type UserPrivate struct {
	ID            uuid.UUID   `json:"id" binding:"required" validate:"required"`
	Email         string      `json:"email" binding:"required" validate:"required,email"`
	Name          string      `json:"name" binding:"required" validate:"required"`
	ProfilePicURL *string     `json:"profilePicUrl,omitempty" validate:"omitempty,url"`
	Roles         []*RoleInfo `json:"roles" validate:"required,dive,required"`
}

func NewUserPrivate(user *model.User) *UserPrivate {
	var roles []*RoleInfo
	for _, role := range user.Roles {
		roles = append(roles, NewRoleInfo(role))
	}

	return &UserPrivate{
		ID:            user.ID,
		Email:         user.Email,
		Name:          user.Name,
		ProfilePicURL: user.ProfilePicURL,
		Roles:         roles,
	}
}
