package usecase

import (
	"myAwesomeProject/internal/entities"
	"myAwesomeProject/internal/repository"
)

type ContactUsecase interface {
	GetContacts(string, string) ([]entities.Contact, error)
}

type contactUsecase struct {
	repo repository.ContactRepository
}

func NewContactUsecase(repo repository.ContactRepository) ContactUsecase {
	return &contactUsecase{repo: repo}
}

func (cu *contactUsecase) GetContacts(token string, subdomain string) ([]entities.Contact, error) {
	return cu.repo.GetContacts(token, subdomain)
}
