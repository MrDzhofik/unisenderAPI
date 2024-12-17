package producer

import (
	"fmt"
	"log"

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

func (p *ContactSyncProducer) AddSyncTask(accountID string) error {
	task := []byte(fmt.Sprintf("Синхронизируем контакты для пользоателя с ID: %s", accountID))

	_, err := p.beanstalkConn.Put(task, 1, 0, 120)
	if err != nil {
		return fmt.Errorf("ошибка добавления задачи в очередь: %v", err)
	}
	log.Printf("Задача на синхронизацию контактов пользователя с ID %s добавлена в очередь", accountID)
	return nil
}
