package entities

type Contact struct {
	Name  string `json:"name"`
	Phone string `json:"phone"`
	Email string `json:"email"`
}

type ResponseContact struct {
	ID                int           `json:"id"`
	Name              string        `json:"name"`
	CustomFieldValues []CustomField `json:"custom_fields_values"`
}

type CustomField struct {
	FieldCode string       `json:"field_code"`
	Values    []FieldValue `json:"values"`
}

type FieldValue struct {
	Value string `json:"value"`
}

func NewContact(rc ResponseContact) Contact {
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
		Name:  rc.Name,
		Phone: phone,
		Email: email,
	}
}
