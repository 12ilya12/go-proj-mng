package main

import (
	"fmt"
	"log"

	"github.com/12ilya12/go-proj-mng/initializers"
	"github.com/12ilya12/go-proj-mng/models"
)

func init() {
	config, err := initializers.LoadConfig()
	if err != nil {
		log.Fatal("Ошибка при загрузке переменных среды", err)
	}
	initializers.ConnectDB(&config)
}

func main() {
	initializers.DB.AutoMigrate(
		&models.User{},
		&models.Category{},
		&models.Dependency{},
		&models.Status{},
		&models.Task{})
	fmt.Println("Миграция завершена")
}
