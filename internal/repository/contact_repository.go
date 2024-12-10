package repository

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"myAwesomeProject/internal/entities"
	"net/http"
	"sync"
)

type ContactRepository interface {
	GetContacts(string, string) ([]entities.Contact, error)
}

type Contacts struct {
	mu       sync.RWMutex
	contacts []entities.Contact
}

type ApiResponse struct {
	Embedded struct {
		Contacts []entities.ResponseContact `json:"contacts"`
	} `json:"_embedded"`
}

func NewContacts() *Contacts {
	return &Contacts{
		contacts: make([]entities.Contact, 0),
	}
}

func (c *Contacts) GetContacts(token string, subdomain string) ([]entities.Contact, error) {
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
		return nil, fmt.Errorf("Ошибка чтения контактов: %s", resp.Status)
	}

	// Читаем тело ответа
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Ошибка чтения тела ответа:", err)
		return nil, err
	}

	var apiResponse ApiResponse
	err = json.Unmarshal(body, &apiResponse)

	if err != nil {
		log.Fatalf("Ошибка парсинга JSON: %v", err)
	}

	response := make([]entities.Contact, 0)

	for _, contact := range apiResponse.Embedded.Contacts {
		newContact := entities.NewContact(contact)
		response = append(response, newContact)
	}

	return response, nil
}
