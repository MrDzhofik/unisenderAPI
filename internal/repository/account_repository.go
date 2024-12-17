package repository

import (
	"errors"
	"myAwesomeProject/internal/entities"

	"gorm.io/gorm"
)

type AccountRepository interface {
	AddAccount(accounts entities.Account) error
	GetAccounts() ([]entities.Account, error)
	AddIntegration(integration entities.AccountIntegration) error
	GetIntegrations() ([]entities.AccountIntegration, error)
	DeleteAccountByID(string) error
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

	return result.Error
}

func (as *AccountStorage) AddIntegration(accountIntegration entities.AccountIntegration) error {
	newIntegration := entities.NewIntegration(accountIntegration.SecretKey, accountIntegration.RedirectURL, accountIntegration.AuthenticationCode)
	result := as.db.Create(&newIntegration)

	return result.Error
}

func (as *AccountStorage) GetAccounts() ([]entities.Account, error) {
	accounts := make([]entities.Account, 0)
	result := as.db.Find(&accounts)

	return accounts, result.Error
}

func (as *AccountStorage) GetIntegrations() ([]entities.AccountIntegration, error) {
	integrations := make([]entities.AccountIntegration, 0)
	result := as.db.Find(&integrations)

	return integrations, result.Error
}

func (as *AccountStorage) DeleteAccountByID(accountID string) error {
	result := as.db.Delete(&entities.Account{}, "account_id = ?", accountID)
	if result.RowsAffected == 0 {
		return errors.New("аккаунт с таким ID не найден")
	}
	return result.Error
}
