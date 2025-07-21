package initializers

import (
	"log"

	reminder "github.com/12ilya12/go-proj-mng/reminder-service/gen/reminder"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func InitReminder() reminder.ReminderServiceClient {
	conn, err := grpc.NewClient("reminder-service:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Нет соединения: %v", err)
	}
	defer conn.Close()

	reminderClient := reminder.NewReminderServiceClient(conn)
	return reminderClient
}
