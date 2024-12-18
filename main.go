package main

import (
	"log"
	"myAwesomeProject/internal/worker"
	"net"
	"os"

	grpcserver "myAwesomeProject/internal/grpc"
	"myAwesomeProject/internal/handlers"
	"myAwesomeProject/internal/producer"
	"myAwesomeProject/internal/repository"
	"myAwesomeProject/internal/usecase"
	"myAwesomeProject/migrations"
	pb "myAwesomeProject/proto/accountpb"
	"net/http"

	"github.com/kr/beanstalk"
	"github.com/urfave/cli/v2"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"github.com/go-gormigrate/gormigrate/v2"
	"github.com/gorilla/mux"
)

func main() {
	// Инициализация
	// Базы данных
	db := Connect()

	log.Println("Подключена база данных!")

	// Сервера очереди сообщений
	conn, err := beanstalk.Dial("tcp", "localhost:11300")
	if err != nil {
		log.Fatalf("Ошибка подключения к Beanstalk: %v", err)
	}
	defer conn.Close()

	log.Println("Подключено к серверу очереди сообщений!")

	accountRepo := repository.NewAccountStorage(db)
	accountUsecase := usecase.NewAccountUsecase(accountRepo)
	accountHandler := handlers.NewAccountHandler(accountUsecase)

	contactRepo := repository.NewContactStorage(db)
	contactUsecase := usecase.NewContactUsecase(contactRepo)
	contactHandler := handlers.NewContactHandler(contactUsecase)

	uniProducer := producer.NewContactSyncProducer(conn)
	uniRepo := repository.NewUnisenderStorage(db)
	uniUsecase := usecase.NewUnisenderUsecase(uniRepo)
	uniHandler := handlers.NewUnisenderHandler(uniUsecase, uniProducer)

	webHookHandler := handlers.NewWebHookHandler(uniProducer)

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

	// WebHook
	r.HandleFunc("/webhook", webHookHandler.HookHandler).Methods("POST")

	log.Println("Роутер успешно настроен!")

	// gRPC сервер
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

	// Создание и запуск CLI-приложения
	app := &cli.App{
		Name:  "worker-cli",
		Usage: "CLI для запуска обработчика задач Beanstalk",
		Commands: []*cli.Command{
			{
				Name:   "run-worker",
				Usage:  "Запускает обработку задач для синхронизации контактов",
				Action: worker.RunWorker,
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
