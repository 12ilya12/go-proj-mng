package worker

import (
	"context"
	"log"
	"time"

	"github.com/12ilya12/go-proj-mng/reminder-service/internal/repos"
)

type Worker struct {
	repo          *repos.RedisRepo
	checkInterval time.Duration
}

func NewWorker(repo *repos.RedisRepo, interval time.Duration) *Worker {
	return &Worker{
		repo:          repo,
		checkInterval: interval,
	}
}

func (w *Worker) Start(ctx context.Context) {
	ticker := time.NewTicker(w.checkInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			w.checkReminders(ctx)
		}
	}
}

func (w *Worker) checkReminders(ctx context.Context) {
	reminders, err := w.repo.GetDueReminders(ctx)
	if err != nil {
		log.Printf("Ошибка при получении напоминаний: %v", err)
		return
	}

	for _, reminder := range reminders {
		log.Printf("НАПОМИНАНИЕ! Задача %s - %s\n", reminder.TaskId, reminder.Message)
		_ = w.repo.DeleteReminder(ctx, reminder.Id)
	}
}
