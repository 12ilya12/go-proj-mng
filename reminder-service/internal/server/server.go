package server

import (
	"context"
	"time"

	"github.com/12ilya12/go-proj-mng/reminder-service/gen/reminder"
	"github.com/12ilya12/go-proj-mng/reminder-service/internal/repos"
	"github.com/google/uuid"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type Server struct {
	repo *repos.RedisRepo

	reminder.UnimplementedReminderServiceServer
}

func NewServer(repo *repos.RedisRepo) *Server {
	return &Server{repo: repo}
}

func (s *Server) CreateReminder(ctx context.Context,
	req *reminder.CreateReminderRequest) (*reminder.Reminder, error) {
	remindAt := req.Deadline.AsTime().Add(-time.Hour * 24 * time.Duration(req.DaysBefore))
	rem := repos.Reminder{
		Id:        uuid.NewString(),
		TaskId:    req.TaskId,
		Message:   req.Message,
		RemindAt:  remindAt,
		Triggered: false,
	}
	err := s.repo.AddReminder(ctx, rem)
	if err != nil {
		return nil, err
	}
	return &reminder.Reminder{
		Id:        rem.Id,
		TaskId:    rem.TaskId,
		Message:   rem.Message,
		RemindAt:  timestamppb.New(rem.RemindAt),
		Triggered: rem.Triggered,
	}, nil
}

//func (s *Server) GetReminders
