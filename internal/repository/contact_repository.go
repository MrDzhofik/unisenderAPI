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
	Page     int `json:"_page"`
	Embedded struct {
		Contacts []entities.Contact `json:"contacts"`
	} `json:"_embedded"`
}

type CustomField struct {
	FieldCode string       `json:"field_code"`
	Values    []FieldValue `json:"values"`
}

type FieldValue struct {
	Value string `json:"value"`
}

func NewContacts() *Contacts {
	return &Contacts{
		contacts: make([]entities.Contact, 0),
	}
}

func (c *Contacts) GetContacts(token string, subdomain string) ([]entities.Contact, error) {
	// Формируем URL API
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

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Ошибка чтения тела ответа:", err)
		return nil, err
	}

	// Вывод тела ответа как строки
	fmt.Println("Тело ответа:")
	fmt.Println(string(body))

	var apiResponse ApiResponse
	err = json.Unmarshal(body, &apiResponse)
	if err != nil {
		log.Fatalf("Ошибка парсинга JSON: %v", err)
	}

	return apiResponse.Embedded.Contacts, nil
}
