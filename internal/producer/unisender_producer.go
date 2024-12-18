package producer

import (
	"encoding/json"
	"fmt"
	"log"
	"myAwesomeProject/internal/entities"

	"github.com/kr/beanstalk"
)

type ContactSyncProducer struct {
	beanstalkConn *beanstalk.Conn
}

func NewContactSyncProducer(conn *beanstalk.Conn) *ContactSyncProducer {
	return &ContactSyncProducer{
		beanstalkConn: conn,
	}
}

func (p *ContactSyncProducer) AddContactTask(task entities.ContactTask) (uint64, error) {
	taskBytes, err := json.Marshal(task)
	if err != nil {
		return 0, fmt.Errorf("ошибка преобразования задачи: %w", err)
	}

	id, err := p.beanstalkConn.Put(taskBytes, 1, 0, 120)
	if err != nil {
		return 0, fmt.Errorf("ошибка добавления задачи в очередь: %v", err)
	}
	log.Printf("Задача на %s контактов добавлена в очередь", task.Action)
	return id, nil
}
