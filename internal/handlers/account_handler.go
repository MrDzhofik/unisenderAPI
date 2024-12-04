package handlers

import (
	"encoding/json"
	"myAwesomeProject/internal/entities"
	"myAwesomeProject/internal/usecase"
	"net/http"
)

type AccountHandler struct {
	accountUsecase usecase.AccountUsecase
}

func NewAccountHandler(accountUsecase usecase.AccountUsecase) *AccountHandler {
	return &AccountHandler{accountUsecase: accountUsecase}
}

func (h *AccountHandler) CreateAccount(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
		return
	}

	var account entities.Account
	if err := json.NewDecoder(r.Body).Decode(&account); err != nil {
		http.Error(w, "Неверные данные", http.StatusBadRequest)
		return
	}

	err := h.accountUsecase.CreateAccount(account)
	if err != nil {
		http.Error(w, "Не удалось добавить аккаунт", http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("Аккаунт создан"))
}

func (h *AccountHandler) GetAccounts(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
		return
	}

	accounts, err := h.accountUsecase.GetAccounts()
	if err != nil {
		http.Error(w, "Не удалось прочитать аккаунты", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(accounts)
}

func (h *AccountHandler) CreateAccountIntegration(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
		return
	}

	var account entities.AccountIntegration
	if err := json.NewDecoder(r.Body).Decode(&account); err != nil {
		http.Error(w, "Неверные данные", http.StatusBadRequest)
		return
	}

	err := h.accountUsecase.CreateIntegration(account)
	if err != nil {
		http.Error(w, "Не удалось добавить интеграцию", http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("Интеграция создана"))
}

func (h *AccountHandler) GetAccountIntegrations(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
		return
	}

	accounts, err := h.accountUsecase.GetIntegrations()
	if err != nil {
		http.Error(w, "Не удалось прочитать интеграции", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(accounts)
}
