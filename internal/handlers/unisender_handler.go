package handlers

import (
	"fmt"
	"myAwesomeProject/internal/entities"
	"myAwesomeProject/internal/producer"
	"myAwesomeProject/internal/usecase"
	"net/http"
)

type UnisenderHandler struct {
	unisenderUsecase usecase.UnisenderUsecase
	producer         *producer.ContactSyncProducer
}

func NewUnisenderHandler(unisenderUsecase usecase.UnisenderUsecase, producer *producer.ContactSyncProducer) *UnisenderHandler {
	return &UnisenderHandler{
		unisenderUsecase: unisenderUsecase,
		producer:         producer,
	}
}

func (uh *UnisenderHandler) SaveUnisenderKey(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, "Ошибка разбора формы", http.StatusBadRequest)
		return
	}

	unisenderKey := r.FormValue("unisender_key")
	accountID := r.FormValue("account_id")

	_, err := uh.producer.AddContactTask(entities.ContactTask{})
	if err != nil {
		http.Error(w, fmt.Sprintf("Ошибка добавления задачи: %v", err), http.StatusInternalServerError)
		return
	}

	err = uh.unisenderUsecase.SaveUnisenderKey(unisenderKey, accountID)

	if err != nil {
		http.Error(w, "Ошибка записи ключа", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf("Задача на синхронизацию контактов пользователя с ID %s добавлена в очередь", accountID)))
}
