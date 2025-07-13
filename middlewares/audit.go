package middlewares

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/12ilya12/go-proj-mng/common"
	"github.com/IBM/sarama"
)

// Протоколируемое событие аудита
type AuditEvent struct {
	UserId    string    `json:"user_id"`   //Идентификатор пользователя, отправляющего запрос
	UserRole  string    `json:"user_role"` //Роль пользователя, отправляющего запрос
	ReqType   string    `json:"req_type"`  //Тип запроса (GET, POST...)
	URL       string    `json:"url"`       //(URL запроса)
	Timestamp time.Time `json:"timestamp"` //Время запроса
	Metadata  string    `json:"metadata"`  //Дополнительные данные (например в формате JSON)
}

func AuditMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//Передаём запрос дальше
		next.ServeHTTP(w, r)

		//Формируем событие аудита
		event := AuditEvent{
			UserId:    fmt.Sprintf("%v", r.Context().Value(common.UserContextKey)),
			UserRole:  fmt.Sprintf("%v", r.Context().Value(common.RoleContextKey)),
			ReqType:   r.Method,
			URL:       r.URL.Path,
			Timestamp: time.Now().UTC(),
		}

		//Отправляем событие в Kafka
		sendAuditEvent(event)
	})
}

func sendAuditEvent(event AuditEvent) {
	//Создание продюсера
	producer, err := sarama.NewSyncProducer([]string{"kafka:9092"}, nil)
	if err != nil {
		log.Printf("Ошибка при создании продюсера Kafka: %v", err)
		return
	}
	//Закрываем продюсер по окончанию отправки сообщения
	defer producer.Close()

	eventJSON, err := json.Marshal(event)
	if err != nil {
		log.Printf("Ошибка при маршалинге структуры AuditEvent в JSON: %v", err)
		return
	}
	//Создаём сообщение
	msg := &sarama.ProducerMessage{
		Topic: "user_audit_log",
		Value: sarama.StringEncoder(eventJSON),
	}
	//И отправляем его
	_, _, err = producer.SendMessage(msg)
	if err != nil {
		log.Printf("Ошибка при отправке сообщения аудита: %v", err)
		return
	}
}
