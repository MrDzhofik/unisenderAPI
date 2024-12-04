package entities

import "github.com/google/uuid"

type Account struct {
	AccountID    string `json:"account_id"`
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	Expires      int    `json:"expires"`
}

func NewAccount(accessToken, refreshToken string, expires int) Account {
	return Account{
		AccountID:    uuid.NewString(),
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		Expires:      expires,
	}
}
