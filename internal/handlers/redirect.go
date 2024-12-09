package handlers

import (
	"net/http"
)

func RedirectHandler(w http.ResponseWriter, r *http.Request) {
	// Обрабатываем данные из redirect
	code := r.URL.Query().Get("code")
	if code == "" {
		http.Error(w, "Не переданы данные", http.StatusBadRequest)
		return
	}

	// Логика обработки
	w.Write([]byte("Редирект"))
}
