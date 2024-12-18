package worker

import (
	"encoding/json"
	"fmt"
	"log"
	"myAwesomeProject/internal/entities"
	"myAwesomeProject/internal/usecase"

	"github.com/kr/beanstalk"
	"github.com/urfave/cli/v2"
)

type Worker struct {
	tube           *beanstalk.Tube
	contactUsecase usecase.ContactUsecase
}

func NewWorker(conn *beanstalk.Conn, tubeName string) *Worker {
	tube := &beanstalk.Tube{
		Conn: conn,
		Name: tubeName,
	}
	return &Worker{tube: tube}
}

func (w *Worker) ProcessTask() {
	for {
		id, body, err := w.tube.PeekReady()
		if err != nil {
			log.Printf("Ошибка при резервировании задачи: %v", err)
			continue
		}

		err = w.handleTask(body)
		if err != nil {
			log.Printf("Ошибка при обработке задачи: %v", err)
		}

		err = w.tube.Conn.Delete(id)
		if err != nil {
			log.Printf("Ошибка при удалении задачи: %v", err)
		} else {
			log.Printf("Задача с ID %d успешно обработана и удалена\n", id)
		}
	}
}

func (w *Worker) handleTask(body []byte) error {
	fmt.Printf("Получена задача: %s\n", body)

	var task entities.ContactTask
	err := json.Unmarshal(body, &task)
	if err != nil {
		return err
	}
	switch task.Action {
	case "Add":
		contact := entities.Contact{
			ClientID: task.ContactID,
			Name:     task.ContactName,
			Email:    task.ContactEmail,
		}
		err = w.contactUsecase.AddContact(contact)
	case "Update":
		contact := entities.Contact{
			ClientID: task.ContactID,
			Name:     task.ContactName,
			Email:    task.ContactEmail,
		}
		err = w.contactUsecase.UpdateContact(contact)
	case "Delete":
		err = w.contactUsecase.DeleteContact(task.ContactID)
	}

	return err
}

func RunWorker(c *cli.Context) error {
	conn, err := beanstalk.Dial("tcp", "localhost:11300")
	if err != nil {
		log.Fatalf("Ошибка подключения к Beanstalk серверу: %v", err)
	}
	defer conn.Close()

	worker := NewWorker(conn, "contact_sync_queue")

	worker.ProcessTask()

	return nil
}
