package message

import (
	"github.com/google/uuid"
)

type RoleCode string

const (
	RoleCodeLearner RoleCode = "LEARNER"
	RoleCodeAdmin   RoleCode = "ADMIN"
	RoleCodeAuthor  RoleCode = "AUTHOR"
	RoleCodeEditor  RoleCode = "EDITOR"
)

type User struct {
	ID            uuid.UUID `json:"id" validate:"required,uuid"`
	Name          string    `json:"name"`
	Email         string    `json:"email"`
	ProfilePicURL *string   `json:"profilePicUrl,omitempty"`
}
