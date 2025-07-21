package main

import (
	"context"
	"log"
	"net"
	"time"

	"github.com/12ilya12/go-proj-mng/reminder-service/gen/reminder"
	"github.com/12ilya12/go-proj-mng/reminder-service/internal/repos"
	"github.com/12ilya12/go-proj-mng/reminder-service/internal/server"
	"github.com/12ilya12/go-proj-mng/reminder-service/internal/worker"
	"google.golang.org/grpc"
)

func main() {
	//Инициализация Redis
	repo := repos.NewRedisRepo("localhost:6379")

	//Запуск gRPC сервера
	listener, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Ошибка при прослушивании TCP: %v", err)
	}

	srv := server.NewServer(repo)
	grpcServer := grpc.NewServer()

	reminder.RegisterReminderServiceServer(grpcServer, srv)

	//Запуск worker
	worker := worker.NewWorker(repo, 10*time.Second)
	go worker.Start(context.Background())

	log.Println("Сервер gRPC стартует на :50051")
	err = grpcServer.Serve(listener)
	if err != nil {
		log.Fatalf("Ошибка: %v", err)
	}
}
