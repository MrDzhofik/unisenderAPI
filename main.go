package main

import (
	"log"
	"net"

	grpcserver "myAwesomeProject/internal/grpc"
	"myAwesomeProject/internal/handlers"
	"myAwesomeProject/internal/repository"
	"myAwesomeProject/internal/usecase"
	"myAwesomeProject/migrations"
	pb "myAwesomeProject/proto/accountpb"
	"net/http"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"github.com/go-gormigrate/gormigrate/v2"
	"github.com/gorilla/mux"
)

func main() {
	// Инициализация
	db := Connect()

	log.Println("Подключена база данных!")

	accountRepo := repository.NewAccountStorage(db)
	accountUsecase := usecase.NewAccountUsecase(accountRepo)
	accountHandler := handlers.NewAccountHandler(accountUsecase)

	contactRepo := repository.NewContactStorage(db)
	contactUsecase := usecase.NewContactUsecase(contactRepo)
	contactHandler := handlers.NewContactHandler(contactUsecase)

	uniRepo := repository.NewUnisenderStorage(db)
	uniUsecase := usecase.NewUnisenderUsecase(uniRepo)
	uniHandler := handlers.NewUnisenderHandler(uniUsecase)

	// Настройка миграций
	migrationsList := []*gormigrate.Migration{
		migrations.CreateAccountMigration(),
		migrations.CreateIntegrationMigration(),
		migrations.CreateContactMigration(),
		migrations.CreateUnisenderMigration(),
	}

	m := gormigrate.New(db, gormigrate.DefaultOptions, migrationsList)

	if err := m.Migrate(); err != nil {
		log.Fatalf("Ошибка применения миграции: %v", err)
	}

	log.Println("Миграции успешно применены!")

	// Настройка роутера
	r := mux.NewRouter()

	// Определение маршрутов

	r.HandleFunc("/", uniHandler.SaveUnisenderKey).Methods("GET", "POST")
	// аккаунт
	r.HandleFunc("/accounts", accountHandler.GetAccounts).Methods("GET")
	r.HandleFunc("/account/create", accountHandler.CreateAccount).Methods("POST")

	// интеграции
	r.HandleFunc("/integrations", accountHandler.GetAccountIntegrations).Methods("GET")
	r.HandleFunc("/integration/create", accountHandler.CreateAccountIntegration).Methods("POST")

	// редирект
	r.HandleFunc("/contacts", contactHandler.GetContacts).Methods("GET")

	log.Println("Роутер успешно настроен!")

	go func() {
		grpcPort := ":8081"
		grpcServer := grpc.NewServer()

		accountServer := grpcserver.NewAccountServer(accountUsecase)
		pb.RegisterAccountServiceServer(grpcServer, accountServer)
		grpc.NewServer()

		lis, err := net.Listen("tcp", grpcPort)
		if err != nil {
			log.Fatalf("Ошибка слушателя: %v", err)
		}

		reflection.Register(grpcServer)

		log.Printf("gRPC сервер слушает на порту %s", grpcPort)
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("Ошибка запуска: %v", err)
		}

		log.Println("gRPC запущен на http://localhost:8081")
	}()

	// Запуск сервера
	log.Println("Сервер запущен на http://localhost:8080")
	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatalf("Ошибка запуска сервера: %v", err)
	}
}
