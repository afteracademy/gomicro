package dto

type TokenRefresh struct {
	RefreshToken string `json:"refreshToken" binding:"required" validate:"required"`
}
