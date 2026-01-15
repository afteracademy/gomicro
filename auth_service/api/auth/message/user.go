package message

import (
	"github.com/afteracademy/gomicro/auth-service/api/user/model"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID            primitive.ObjectID `json:"_id"`
	Name          string             `json:"name"`
	Email         string             `json:"email"`
	ProfilePicURL *string            `json:"profilePicUrl,omitempty"`
}

func NewUser(user *model.User) *User {
	return &User{
		ID:            user.ID,
		Name:          user.Name,
		Email:         user.Email,
		ProfilePicURL: user.ProfilePicURL,
	}
}
