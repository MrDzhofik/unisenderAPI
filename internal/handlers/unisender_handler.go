package handlers

import (
	"myAwesomeProject/internal/usecase"
	"net/http"
)

type UnisenderHandler struct {
	unisenderUsecase usecase.UnisenderUsecase
}

func NewUnisenderHandler(unisenderUsecase usecase.UnisenderUsecase) *UnisenderHandler {
	return &UnisenderHandler{unisenderUsecase: unisenderUsecase}
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

	err := uh.unisenderUsecase.SaveUnisenderKey(unisenderKey, accountID)

	if err != nil {
		http.Error(w, "Ошибка записи ключа", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

}
