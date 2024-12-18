package worker

import (
	"encoding/json"
	"log"
	"myAwesomeProject/internal/entities"
	"myAwesomeProject/internal/usecase"

	"github.com/kr/beanstalk"
)

type Worker struct {
	tube           *beanstalk.Tube
	contactUsecase usecase.ContactUsecase
}

func NewWorker(conn *beanstalk.Conn, tubeName string, contactUsecase usecase.ContactUsecase) *Worker {
	tube := &beanstalk.Tube{
		Conn: conn,
		Name: tubeName,
	}
	return &Worker{
		tube:           tube,
		contactUsecase: contactUsecase,
	}
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
