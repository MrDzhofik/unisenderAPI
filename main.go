package main

import (
	"fmt"
	"log"
	db2 "myAwesomeProject/db"
	"myAwesomeProject/internal/handlers"
	"myAwesomeProject/internal/producer"
	"myAwesomeProject/migrations"
	pb "myAwesomeProject/proto/accountpb"
	"net"
	"net/http"

	"myAwesomeProject/internal/repository"
	"myAwesomeProject/internal/usecase"
	"os"
	"os/exec"

	grpcserver "myAwesomeProject/internal/grpc"

	"github.com/go-gormigrate/gormigrate/v2"
	"github.com/gorilla/mux"
	"github.com/kr/beanstalk"
	"github.com/urfave/cli/v2"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	// Инициализация
	// Базы данных
	db := db2.Connect()

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

	// Создание и запуск CLI-приложения
	app := &cli.App{
		Name:  "worker-manager",
		Usage: "Управление воркерами для обработки задач Beanstalk",
		Commands: []*cli.Command{
			{
				Name:  "start-workers",
				Usage: "Запускает указанное количество воркеров",
				Flags: []cli.Flag{
					&cli.IntFlag{
						Name:     "count",
						Usage:    "Количество воркеров для запуска",
						Required: true,
					},
				},
				Action: func(c *cli.Context) error {
					count := c.Int("count")
					if count <= 0 {
						return fmt.Errorf("количество воркеров должно быть больше 0")
					}

					for i := 0; i < count; i++ {
						if err := startWorker(contactUsecase); err != nil {
							log.Printf("Ошибка при запуске воркера %d: %v", i, err)
						}
					}

					return nil
				},
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}

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
}

func startWorker(contactUsecase usecase.ContactUsecase) error {
	cmd := exec.Command("./worker")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	log.Printf("Запускаем воркер...")
	conn, err := beanstalk.Dial("tcp", "localhost:11300")
	if err != nil {
		log.Fatalf("Ошибка подключения к Beanstalk серверу: %v", err)
	}
	worker := NewWorker(conn, "default", contactUsecase)
	worker.ProcessTask()

	return cmd.Start()
}
