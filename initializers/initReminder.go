package initializers

import (
	"log"

	reminder "github.com/12ilya12/go-proj-mng/reminder-service/gen/reminder"
	"google.golang.org/grpc"
)

func InitReminder() (reminder.ReminderServiceClient, error) {
	conn, err := grpc.Dial("reminder-service:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Нет соединения: %v", err)
	}
	defer conn.Close()

	reminderClient := reminder.NewReminderServiceClient(conn)
	return reminderClient, nil
}
