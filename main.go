package main

import (
	"log"
	"myAwesomeProject/internal/handlers"
	"myAwesomeProject/internal/repository"
	"myAwesomeProject/internal/usecase"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	// Инициализация
	repo := repository.NewInMemoryDB()
	accountUsecase := usecase.NewAccountUsecase(repo)
	accountHandler := handlers.NewAccountHandler(accountUsecase)

	r := mux.NewRouter()

	// Определение маршрутов
	r.HandleFunc("/accounts", accountHandler.GetAccounts).Methods("GET")
	r.HandleFunc("/account/create", accountHandler.CreateAccount).Methods("POST")
	r.HandleFunc("/integrations", accountHandler.GetAccountIntegrations).Methods("GET")
	r.HandleFunc("/integration/create", accountHandler.CreateAccountIntegration).Methods("POST")

	// Запуск сервера
	log.Println("Сервер запущен на http://localhost:8080")
	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatalf("Ошибка запуска сервера: %v", err)
	}
}
