package repository

import (
	"myAwesomeProject/internal/entities"

	"gorm.io/gorm"
)

type AccountRepository interface {
	AddAccount(accounts entities.Account) error
	GetAccounts() ([]entities.Account, error)
	AddIntegration(integration entities.AccountIntegration) error
	GetIntegrations() ([]entities.AccountIntegration, error)
}

type AccountStorage struct {
	db *gorm.DB
}

func NewAccountStorage(db *gorm.DB) *AccountStorage {
	return &AccountStorage{
		db: db,
	}
}

func (as *AccountStorage) AddAccount(account entities.Account) error {
	newAccount := entities.NewAccount(account.AccessToken, account.RefreshToken, account.Expires)
	result := as.db.Create(&newAccount)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (as *AccountStorage) AddIntegration(accountIntegration entities.AccountIntegration) error {
	newIntegration := entities.NewIntegration(accountIntegration.SecretKey, accountIntegration.RedirectURL, accountIntegration.AuthenticationCode)
	result := as.db.Create(&newIntegration)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (as *AccountStorage) GetAccounts() ([]entities.Account, error) {
	accounts := make([]entities.Account, 0)
	result := as.db.Find(&accounts)
	if result.Error != nil {
		return nil, result.Error
	}

	return accounts, nil
}

func (as *AccountStorage) GetIntegrations() ([]entities.AccountIntegration, error) {
	integrations := make([]entities.AccountIntegration, 0)
	result := as.db.Find(&integrations)
	if result.Error != nil {
		return nil, result.Error
	}

	return integrations, nil
}
