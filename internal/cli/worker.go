package main

import (
	"fmt"
	"log"
	db2 "myAwesomeProject/db"
	"myAwesomeProject/internal/repository"
	"myAwesomeProject/internal/usecase"
	"myAwesomeProject/internal/worker"
	"os"
	"os/exec"

	"github.com/kr/beanstalk"
	"github.com/urfave/cli/v2"
)

func startWorker(contactUsecase usecase.ContactUsecase) error {
	cmd := exec.Command("./worker2")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	log.Printf("Запускаем воркер...")
	conn, err := beanstalk.Dial("tcp", "localhost:11300")
	if err != nil {
		log.Fatalf("Ошибка подключения к Beanstalk серверу: %v", err)
	}
	worker2 := worker.NewWorker(conn, "default", contactUsecase)
	worker2.ProcessTask()

	return cmd.Start()
}

func main() {
	db := db2.Connect()
	contactRepo := repository.NewContactStorage(db)
	contactUsecase := usecase.NewContactUsecase(contactRepo)
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
}
