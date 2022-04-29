package service

import "context"

type GoogleAuthClient interface {
	RetrieveAccessToken(ctx context.Context, code string) (*GoogleAuthResponse, error)
	RetrieveUserInfo(ctx context.Context, googleAuthResponse *GoogleAuthResponse) (*GoogleUserInfo, error)
}

type GoogleAuthResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type GoogleUserInfo struct {
	Email string `json:"email"`
	Name  string `json:"name"`
}
