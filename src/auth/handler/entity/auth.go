package entity

type AuthResponse struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}

type RefreshTokenParameter struct {
	RefreshToken string `json:"refreshToken"`
}
