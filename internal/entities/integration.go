package entities

import "github.com/google/uuid"

type AccountIntegration struct {
	ClientID           string `json:"client_id"`
	SecretKey          string `json:"secret_key"`
	RedirectURL        string `json:"redirect_url"`
	AuthenticationCode string `json:"authentication_code"`
}

func NewIntegration(secretKey, redirectURL, authenticationCode string) AccountIntegration {
	return AccountIntegration{
		ClientID:           uuid.NewString(),
		SecretKey:          secretKey,
		RedirectURL:        redirectURL,
		AuthenticationCode: authenticationCode,
	}
}
