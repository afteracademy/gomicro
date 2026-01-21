package dto

import (
	"github.com/afteracademy/gomicro/auth-service/api/user/model"
	"github.com/google/uuid"
)

type RoleInfo struct {
	ID   uuid.UUID      `json:"id" binding:"required" validate:"required"`
	Code model.RoleCode `json:"code" binding:"required" validate:"required,uppercase"`
}

func NewRoleInfo(role *model.Role) *RoleInfo {
	return &RoleInfo{
		ID:   role.ID,
		Code: role.Code,
	}
}
