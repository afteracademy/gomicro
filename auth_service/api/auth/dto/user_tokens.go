package dto

type UserTokens struct {
	AccessToken  string `json:"accessToken" binding:"required" validate:"required"`
	RefreshToken string `json:"refreshToken" binding:"required" validate:"required"`
}

func NewUserTokens(access string, refresh string) *UserTokens {
	return &UserTokens{
		AccessToken:  access,
		RefreshToken: refresh,
	}
}
