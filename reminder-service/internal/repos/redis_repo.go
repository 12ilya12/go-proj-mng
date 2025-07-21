package repos

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"
)

type Reminder struct {
	Id        string
	TaskId    string
	Message   string
	RemindAt  time.Time
	Triggered bool
}

type RedisRepo struct {
	client *redis.Client
}

func NewRedisRepo(addr string) *RedisRepo {
	rdb := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: "",
		DB:       0,
	})
	return &RedisRepo{client: rdb}
}

func (r *RedisRepo) AddReminder(ctx context.Context, rem Reminder) error {
	err := r.client.ZAdd(ctx, "reminders", redis.Z{
		Score:  float64(rem.RemindAt.Unix()), //Сортируем по времени напоминания
		Member: fmt.Sprintf("%s|%s|%s", rem.Id, rem.TaskId, rem.Message),
	}).Err()
	return err
}

// Получение выборки напоминаний, время которых меньше или равно текущему времени
func (r *RedisRepo) GetDueReminders(ctx context.Context) ([]Reminder, error) {
	now := time.Now().Unix()

	res, err := r.client.ZRangeByScore(ctx, "reminders", &redis.ZRangeBy{
		Min: "0",
		Max: fmt.Sprintf("%d", now),
	}).Result()

	if err != nil {
		return nil, err
	}

	var reminders []Reminder
	for _, item := range res {
		parts := strings.Split(item, "|")
		if len(parts) < 3 {
			continue
		}
		reminders = append(reminders, Reminder{
			Id:      parts[0],
			TaskId:  parts[1],
			Message: parts[2],
		})
	}
	return reminders, nil
}

func (r *RedisRepo) DeleteReminder(ctx context.Context, id string) error {
	return r.client.ZRem(ctx, "reminders", id).Err()
}
