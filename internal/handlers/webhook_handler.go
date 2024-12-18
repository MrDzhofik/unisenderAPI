package handlers

import (
	"encoding/json"
	"io"
	"myAwesomeProject/internal/entities"
	"myAwesomeProject/internal/producer"
	"net/http"
)

type WebHookHandler struct {
	producer *producer.ContactSyncProducer
}

type webHookBody struct {
	Contacts struct {
		Add    []entities.ResponseContact `json:"add"`
		Update []entities.ResponseContact `json:"update"`
		Delete []entities.ResponseContact `json:"delete"`
	} `json:"contacts"`
}

func NewWebHookHandler(producer *producer.ContactSyncProducer) *WebHookHandler {
	return &WebHookHandler{
		producer: producer,
	}
}

func (wh *WebHookHandler) HookHandler(w http.ResponseWriter, r *http.Request) {
	var response webHookBody

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Переданы неправильные или неполные данные", http.StatusBadRequest)
	}
	defer r.Body.Close()

	err = json.Unmarshal(body, &response)
	if err != nil {
		http.Error(w, "Неправильный формат JSON", http.StatusBadRequest)
		return
	}

	ids := make([]uint64, 0)

	for _, contact := range response.Contacts.Add {
		newContact := entities.NewContact(contact, contact.AccountID)
		newTask := entities.NewContactTask("Add", newContact.ID, newContact.Name, newContact.Email)

		id, err := wh.producer.AddContactTask(newTask)

		if err != nil {
			http.Error(w, "Ошибка добавления задачи в очередь", http.StatusInternalServerError)
		}

		ids = append(ids, id)
	}

	for _, contact := range response.Contacts.Update {
		newContact := entities.NewContact(contact, contact.AccountID)
		newTask := entities.NewContactTask("Update", newContact.ID, newContact.Name, newContact.Email)

		id, err := wh.producer.AddContactTask(newTask)

		if err != nil {
			http.Error(w, "Ошибка добавления задачи в очередь", http.StatusInternalServerError)
		}

		ids = append(ids, id)
	}

	for _, contact := range response.Contacts.Delete {
		newContact := entities.NewContact(contact, contact.AccountID)
		newTask := entities.NewContactTask("Delete", newContact.ID, newContact.Name, newContact.Email)

		id, err := wh.producer.AddContactTask(newTask)

		if err != nil {
			http.Error(w, "Ошибка добавления задачи в очередь", http.StatusInternalServerError)
		}

		ids = append(ids, id)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Все задачи успешно выставлены в очередь!"))
}
