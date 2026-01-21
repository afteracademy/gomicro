package dto

type Tokens struct {
	AccessToken  string `json:"accessToken" binding:"required" validate:"required"`
	RefreshToken string `json:"refreshToken" binding:"required" validate:"required"`
}

func NewTokens(access string, refresh string) *Tokens {
	return &Tokens{
		AccessToken:  access,
		RefreshToken: refresh,
	}
}
