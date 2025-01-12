package usecase

import (
	"myAwesomeProject/internal/entities"
	"myAwesomeProject/internal/repository"
)

type ContactUsecase interface {
	GetContacts(string, string) ([]entities.Contact, error)
	AddContact(entities.Contact) error
	UpdateContact(entities.Contact) error
	DeleteContact(string) error
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

func (cu *contactUsecase) AddContact(contact entities.Contact) error {
	return cu.repo.AddContact(contact)
}

func (cu *contactUsecase) UpdateContact(contact entities.Contact) error {
	return cu.repo.UpdateContact(contact)
}

func (cu *contactUsecase) DeleteContact(contactID string) error {
	return cu.repo.DeleteContact(contactID)
}
