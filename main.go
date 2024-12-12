package main

import (
	"log"
	"myAwesomeProject/internal/handlers"
	"myAwesomeProject/internal/repository"
	"myAwesomeProject/internal/usecase"
	"myAwesomeProject/migrations"
	"net/http"

	"github.com/go-gormigrate/gormigrate/v2"
	"github.com/gorilla/mux"
)

func main() {
	// Инициализация
	db := Connect()

	accountRepo := repository.NewAccountStorage(db)
	accountUsecase := usecase.NewAccountUsecase(accountRepo)
	accountHandler := handlers.NewAccountHandler(accountUsecase)

	contactRepo := repository.NewContactStorage(db)
	contactUsecase := usecase.NewContactUsecase(contactRepo)
	contactHandler := handlers.NewContactHandler(contactUsecase)

	// Настройка миграций
	migrationsList := []*gormigrate.Migration{
		migrations.CreateAccountMigration(),
		migrations.CreateIntegrationMigration(),
		migrations.CreateContactMigration(),
	}

	m := gormigrate.New(db, gormigrate.DefaultOptions, migrationsList)

	if err := m.Migrate(); err != nil {
		log.Fatalf("Ошибка применения миграции: %v", err)
	}

	log.Println("Миграции успешно применены!")

	// Настройка роутера
	r := mux.NewRouter()

	// Определение маршрутов

	// аккаунт
	r.HandleFunc("/accounts", accountHandler.GetAccounts).Methods("GET")
	r.HandleFunc("/account/create", accountHandler.CreateAccount).Methods("POST")

	// интеграции
	r.HandleFunc("/integrations", accountHandler.GetAccountIntegrations).Methods("GET")
	r.HandleFunc("/integration/create", accountHandler.CreateAccountIntegration).Methods("POST")

	// редирект
	r.HandleFunc("/redirect", handlers.RedirectHandler).Methods("GET")
	r.HandleFunc("/contacts", contactHandler.GetContacts).Methods("GET")

	// Запуск сервера
	log.Println("Сервер запущен на http://localhost:8080")
	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatalf("Ошибка запуска сервера: %v", err)
	}
}
