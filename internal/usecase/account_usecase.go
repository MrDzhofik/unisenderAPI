package usecase

import (
	"myAwesomeProject/internal/entities"
	"myAwesomeProject/internal/repository"
	pb "myAwesomeProject/proto/accountpb"
)

type AccountUsecase interface {
	CreateAccount(account entities.Account) error
	GetAccounts() ([]entities.Account, error)
	CreateIntegration(integration entities.AccountIntegration) error
	GetIntegrations() ([]entities.AccountIntegration, error)
	DeleteAccountByID(string) error
}

type accountUsecase struct {
	repo repository.AccountRepository
	pb.UnimplementedAccountServiceServer
}

func NewAccountUsecase(repo repository.AccountRepository) AccountUsecase {
	return &accountUsecase{repo: repo}
}

func (u *accountUsecase) mustEmbedUnimplementedAccountServiceServer() {
	return
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

func (u *accountUsecase) DeleteAccountByID(accountID string) error {
	return u.repo.DeleteAccountByID(accountID)
}
