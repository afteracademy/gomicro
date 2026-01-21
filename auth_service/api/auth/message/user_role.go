package message

type UserRole struct {
	User  *User    `json:"user" validate:"required"`
	Roles []string `json:"roles" validate:"required,dive"`
}

func NewUserRole(user *User, roles ...string) *UserRole {
	return &UserRole{
		User:  user,
		Roles: roles,
	}
}
