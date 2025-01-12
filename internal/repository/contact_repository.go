package repository

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"myAwesomeProject/internal/entities"
	"net/http"

	"gorm.io/gorm"
)

type ContactRepository interface {
	GetContacts(string, string) ([]entities.Contact, error)
	AddContact(entities.Contact) error
	UpdateContact(entities.Contact) error
	DeleteContact(string) error
}

type ContactStorage struct {
	db *gorm.DB
}

type ApiResponse struct {
	Embedded struct {
		Contacts []entities.ResponseContact `json:"contacts"`
	} `json:"_embedded"`
}

func NewContactStorage(db *gorm.DB) *ContactStorage {
	return &ContactStorage{
		db: db,
	}
}

func (cs *ContactStorage) GetContacts(token string, subdomain string) ([]entities.Contact, error) {
	url := fmt.Sprintf("https://%s.amocrm.ru/api/v4/contacts", subdomain)

	// Создаём HTTP-запрос
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	// Добавляем заголовок авторизации
	req.Header.Set("Authorization", "Bearer "+token)

	// Отправляем запрос
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Проверяем статус-код ответа
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("ошибка чтения контактов: %s", resp.Status)
	}

	// Читаем тело ответа
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var apiResponse ApiResponse
	err = json.Unmarshal(body, &apiResponse)

	if err != nil {
		log.Fatalf("Ошибка парсинга JSON: %v", err)
	}

	var account entities.Account
	result := cs.db.Where("access+token = ?", token).First(account)
	if result.Error != nil {
		return nil, result.Error
	}

	response := make([]entities.Contact, 0)

	for _, contact := range apiResponse.Embedded.Contacts {
		newContact := entities.NewContact(contact, account.AccountID)
		response = append(response, newContact)
	}

	result = cs.db.Create(&response)
	if result.Error != nil {
		return nil, result.Error
	}

	return response, nil
}

func (cs *ContactStorage) AddContact(contact entities.Contact) error {
	result := cs.db.Create(&contact)

	return result.Error
}

func (cs *ContactStorage) UpdateContact(contact entities.Contact) error {
	result := cs.db.Model(&entities.Contact{}).Where("id = ?", contact.ID).Updates(contact)

	return result.Error
}

func (cs *ContactStorage) DeleteContact(contactID string) error {
	contact := entities.Contact{
		ID: contactID,
	}
	result := cs.db.Where("id = ?", contactID).Delete(&contact)

	return result.Error
}
