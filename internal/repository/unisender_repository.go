package repository

import (
	"myAwesomeProject/internal/entities"

	"gorm.io/gorm"
)

type UnisenderRepository interface {
	SaveUnisenderKey(string, string) error
}

type UnisenderStorage struct {
	db *gorm.DB
}

func NewUnisenderStorage(db *gorm.DB) *UnisenderStorage {
	return &UnisenderStorage{
		db: db,
	}
}

func (us *UnisenderStorage) SaveUnisenderKey(key string, accountID string) error {
	newUnisenderKey := entities.NewUnisender(key, accountID)
	result := us.db.Create(&newUnisenderKey)

	return result.Error
}
