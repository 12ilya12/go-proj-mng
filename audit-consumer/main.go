package main

import (
	"encoding/json"
	"log"
	"strings"

	"github.com/IBM/sarama"
	elasticsearch "github.com/elastic/go-elasticsearch"
)

// Протоколируемое событие аудита
type AuditEvent struct {
	UserId    string `json:"user_id"`   //Идентификатор пользователя, отправляющего запрос
	UserRole  string `json:"user_role"` //Роль пользователя, отправляющего запрос
	ReqType   string `json:"req_type"`  //Тип запроса (GET, POST...)
	URL       string `json:"url"`       //(URL запроса)
	Timestamp string `json:"timestamp"` //Время запроса
	Metadata  string `json:"metadata"`  //Дополнительные данные (например в формате JSON)
}

func main() {
	//Подключаемся к Kafka
	consumer, err := sarama.NewConsumer([]string{"kafka:9092"}, nil)
	if err != nil {
		log.Fatalf("Ошибка при создании потребителя Kafka: %v", err)
	}
	//По окончании работы сервиса, закрываем потребителя
	defer consumer.Close()

	//Подписываемся на топик
	partitionConsumer, err := consumer.ConsumePartition("user_audit_log", 0, sarama.OffsetOldest)
	if err != nil {
		log.Fatalf("Не удалось подписаться на топик user_audit_log: %v", err)
	}
	//По окончании работы сервиса, закрываем потребителя топика
	defer partitionConsumer.Close()

	//Подключаемся к ElasticSearch
	es, err := elasticsearch.NewClient(elasticsearch.Config{
		Addresses: []string{"http://elasticsearch:9200"},
	})
	if err != nil {
		log.Fatalf("Не удалось соединиться с Elasticsearch: %v", err)
	}

	//Обработка сообщений
	for msg := range partitionConsumer.Messages() {
		var event AuditEvent
		if err := json.Unmarshal(msg.Value, &event); err != nil {
			log.Printf("Произошла ошибка при парсинге события аудита: %v", err)
			continue
		}
		//Сохранение в ElasticSearch
		doc, _ := json.Marshal(event)
		res, err := es.Index(
			"audit_log",                    //Индекс
			strings.NewReader(string(doc)), //Тело документа
		)
		if err != nil {
			log.Printf("Ошибка сохранения в ElasticSearch: %v", err)
		} else {
			log.Printf("Сохранение выполнено. Статус: %s", res.Status())
		}
	}
}
