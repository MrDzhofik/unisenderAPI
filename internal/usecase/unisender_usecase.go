package usecase

import (
	"myAwesomeProject/internal/repository"
)

type UnisenderUsecase interface {
	SaveUnisenderKey(string, string) error
}

type unisenderUsecase struct {
	repo repository.UnisenderRepository
}

func NewUnisenderUsecase(repo repository.UnisenderRepository) UnisenderUsecase {
	return &unisenderUsecase{repo: repo}
}

func (cu *unisenderUsecase) SaveUnisenderKey(key string, accountID string) error {
	return cu.repo.SaveUnisenderKey(key, accountID)
}
