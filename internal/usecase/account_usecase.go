package usecase

import (
	"myAwesomeProject/internal/entities"
	"myAwesomeProject/internal/repository"
)

type AccountUsecase interface {
	CreateAccount(account entities.Account) error
	GetAccounts() ([]entities.Account, error)
	CreateIntegration(integration entities.AccountIntegration) error
	GetIntegrations() ([]entities.AccountIntegration, error)
}

type accountUsecase struct {
	repo repository.AccountRepository
}

func NewAccountUsecase(repo repository.AccountRepository) AccountUsecase {
	return &accountUsecase{repo: repo}
}

func (u *accountUsecase) CreateAccount(account entities.Account) error {
	return u.repo.AddAccount(account)
}

func (u *accountUsecase) GetAccounts() ([]entities.Account, error) {
	return u.repo.GetAccounts()
}

func (u *accountUsecase) CreateIntegration(integration entities.AccountIntegration) error {
	return u.repo.AddIntegration(integration)
}

func (u *accountUsecase) GetIntegrations() ([]entities.AccountIntegration, error) {
	return u.repo.GetIntegrations()
}
