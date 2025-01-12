package entities

import "github.com/google/uuid"

type AccountIntegration struct {
	ID                 string  `json:"-" gorm:"primary_key"`
	ClientID           string  `json:"client_id" gorm:"not null;size:191"`
	Account            Account `gorm:"foreignkey:ClientID;constraint:OnDelete:CASCADE;" json:"-"'` // Внешний ключ
	SecretKey          string  `json:"secret_key" gorm:"not null"`
	RedirectURL        string  `json:"redirect_url" gorm:"not null"`
	AuthenticationCode string  `json:"authentication_code"`
}

func NewIntegration(secretKey, redirectURL, authenticationCode string) AccountIntegration {
	return AccountIntegration{
		ID:                 uuid.NewString(),
		ClientID:           "a6a4efc1-156d-4923-a13e-c64d94288235",
		SecretKey:          secretKey,
		RedirectURL:        redirectURL,
		AuthenticationCode: authenticationCode,
	}
}
