package entities

import "github.com/google/uuid"

type Account struct {
	AccountID    string `json:"account_id" gorm:"primary_key"`
	AccessToken  string `json:"access_token" gorm:"not null"`
	RefreshToken string `json:"refresh_token" gorm:"not null"`
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
