package handlers

import (
	"encoding/json"
	"myAwesomeProject/internal/usecase"
	"net/http"
)

type ContactHandler struct {
	contactUsecase usecase.ContactUsecase
}

func NewContactHandler(contactUsecase usecase.ContactUsecase) *ContactHandler {
	return &ContactHandler{contactUsecase: contactUsecase}
}

func (h *ContactHandler) GetContacts(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
		return
	}

	token := "eyJhdWQiOiI5YTEwYjBjOC1iODY4LTRjZTQtODczYS0xY2RjZDFmMjAwZmIiLCJqdGkiOiJiM2FiOTQyZjdlODY4MGVkNzZlZGQzNDA3NDc2MTdmZDIxOTgwYjc1ZGFhNDJlZTA3M2UwNmYxMjAzMGRhZWQ1ZDUxMjA1MzUxNDNlM2E0ZiIsImlhdCI6MTczNDAwMDkyMiwibmJmIjoxNzM0MDAwOTIyLCJleHAiOjE3MzQwODczMjIsInN1YiI6IjExODUwMDU4IiwiZ3JhbnRfdHlwZSI6IiIsImFjY291bnRfaWQiOjMyMTA1NjAyLCJiYXNlX2RvbWFpbiI6ImFtb2NybS5ydSIsInZlcnNpb24iOjIsInNjb3BlcyI6WyJwdXNoX25vdGlmaWNhdGlvbnMiLCJmaWxlcyIsImNybSIsImZpbGVzX2RlbGV0ZSIsIm5vdGlmaWNhdGlvbnMiXSwiaGFzaF91dWlkIjoiMjI1NTZjYzMtZTIxYy00N2E0LTgxOTAtOTk5NDRkZjIzOTAzIiwiYXBpX2RvbWFpbiI6ImFwaS1iLmFtb2NybS5ydSJ9"
	subdomain := "emdzaharovtest"

	contacts, err := h.contactUsecase.GetContacts(token, subdomain)
	if err != nil {
		http.Error(w, "Не удалось получить контакты", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(contacts)
}
