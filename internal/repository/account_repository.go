package repository

import (
	"myAwesomeProject/internal/entities"
	"sync"
)

type AccountRepository interface {
	AddAccount(accounts entities.Account) error
	GetAccounts() ([]entities.Account, error)
	AddIntegration(integration entities.AccountIntegration) error
	GetIntegrations() ([]entities.AccountIntegration, error)
}

type InMemoryDB struct {
	mu                 sync.RWMutex
	accounts           map[string]entities.Account
	accountIntegration map[string]entities.AccountIntegration
}

func NewInMemoryDB() *InMemoryDB {
	return &InMemoryDB{
		accounts:           make(map[string]entities.Account),
		accountIntegration: make(map[string]entities.AccountIntegration),
	}
}

func (db *InMemoryDB) AddAccount(account entities.Account) error {
	db.mu.Lock()
	defer db.mu.Unlock()
	newAccount := entities.NewAccount(account.AccessToken, account.RefreshToken, account.Expires)
	db.accounts[account.AccountID] = newAccount

	return nil
}

func (db *InMemoryDB) AddIntegration(accountIntegration entities.AccountIntegration) error {
	db.mu.Lock()
	defer db.mu.Unlock()
	newIntegration := entities.NewIntegration(accountIntegration.SecretKey, accountIntegration.RedirectURL, accountIntegration.AuthenticationCode)
	db.accountIntegration[accountIntegration.ClientID] = newIntegration

	return nil
}

func (db *InMemoryDB) GetAccounts() ([]entities.Account, error) {
	db.mu.RLock()
	defer db.mu.RUnlock()

	accounts := make([]entities.Account, 0)
	for _, account := range db.accounts {
		accounts = append(accounts, account)
	}
	return accounts, nil
}

func (db *InMemoryDB) GetIntegrations() ([]entities.AccountIntegration, error) {
	db.mu.RLock()
	defer db.mu.RUnlock()

	accountIntegrations := make([]entities.AccountIntegration, 0)
	for _, accountIntegration := range db.accountIntegration {
		accountIntegrations = append(accountIntegrations, accountIntegration)
	}

	return accountIntegrations, nil
}
