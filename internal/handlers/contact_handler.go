package handlers

import (
	"encoding/json"
	"fmt"
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

	token := "eyJ0eXAiOiJKV1QiLCJhbGciOiJSUzI1NiIsImp0aSI6IjNkZGMxOGQ5OTZlZDcyMWJiZTFkMzBmYTMzZj" +
		"ZjYjcyYTI1ZmNmMjc2ZWMzODg5MGZhZGI2ZTMzZDBkNmMzMjYxYjJiMTM2NDk4MzRjZTUyIn0.eyJhdWQiOiI5Y" +
		"TEwYjBjOC1iODY4LTRjZTQtODczYS0xY2RjZDFmMjAwZmIiLCJqdGkiOiIzZGRjMThkOTk2ZWQ3MjFiYmUxZDMw" +
		"ZmEzM2Y2Y2I3MmEyNWZjZjI3NmVjMzg4OTBmYWRiNmUzM2QwZDZjMzI2MWIyYjEzNjQ5ODM0Y2U1MiIsImlhdCI" +
		"6MTczMzc1MTYwNCwibmJmIjoxNzMzNzUxNjA0LCJleHAiOjE3MzM4MzgwMDQsInN1YiI6IjExODUwMDU4IiwiZ3" +
		"JhbnRfdHlwZSI6IiIsImFjY291bnRfaWQiOjMyMTA1NjAyLCJiYXNlX2RvbWFpbiI6ImFtb2NybS5ydSIsInZlc" +
		"nNpb24iOjIsInNjb3BlcyI6WyJwdXNoX25vdGlmaWNhdGlvbnMiLCJmaWxlcyIsImNybSIsImZpbGVzX2RlbGV0" +
		"ZSIsIm5vdGlmaWNhdGlvbnMiXSwiaGFzaF91dWlkIjoiM2MzZWJmNzUtZjUxNS00YjFmLTk1NGUtM2YwNGMyNmQ" +
		"2NmQzIiwiYXBpX2RvbWFpbiI6ImFwaS1iLmFtb2NybS5ydSJ9.VYrbQ-J-oFeUknIbQQ2rHsvn4M1XN81kVFIfb" +
		"rQN2m2yo2B4pvDc2aaOucru-lTBZWlvaOL259SI96dQKphpYI4n1YtGw92dL5fcnSZWEzn4NucfniitG_0Qkqj0" +
		"e-jDlwyO9lU2iraKinypkhOby71w8KZAW-aoD9OzjDpWYd7SHP83xjo4sDz0oevAJDcnv7kI5IEGyTKPLz3kcPE" +
		"dR8Pqx2JCKLO1U33oTOsmUe9eCwFibOOq83urpfCfpGK6zt_A8h_SbwuT-nXtv1F9MROu6_Q_Xw0HSWDpm-vyyj" +
		"CR410eYjAdLzNpwONZv7a7ygqMUuOs3n9YxwsSW0soUw"
	subdomain := "emdzaharovtest"

	contacts, err := h.contactUsecase.GetContacts(token, subdomain)
	fmt.Println("Handler: ", err)
	if err != nil {
		http.Error(w, "Не удалось получить контакты", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(contacts)
}
