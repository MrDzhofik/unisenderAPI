package entities

import "github.com/google/uuid"

type Contact struct {
	ID       string  `json:"-" gorm:"primary_key"`
	ClientID string  `json:"-" gorm:"not null;size:191"`
	Account  Account `gorm:"foreignkey:ClientID;constraint:OnDelete:CASCADE;" json:"-"` // Внешний ключ
	Name     string  `json:"name" gorm:"not null"`
	Phone    string  `json:"phone"`
	Email    string  `json:"email" gorm:"not null; unique"`
}

type ResponseContact struct {
	ID                string        `json:"id"`
	Name              string        `json:"name"`
	AccountID         string        `json:"account_id"`
	CustomFieldValues []CustomField `json:"custom_fields_values"`
}

type CustomField struct {
	FieldCode string       `json:"field_code"`
	Values    []FieldValue `json:"values"`
}

type FieldValue struct {
	Value string `json:"value"`
}

func NewContact(rc ResponseContact, accountID string) Contact {
	phone := ""
	email := ""

	for _, field := range rc.CustomFieldValues {
		if field.FieldCode == "PHONE" {
			phone = field.Values[0].Value
		} else if field.FieldCode == "EMAIL" {
			email = field.Values[0].Value
		}
	}

	return Contact{
		ID:       uuid.NewString(),
		ClientID: accountID,
		Name:     rc.Name,
		Phone:    phone,
		Email:    email,
	}
}
