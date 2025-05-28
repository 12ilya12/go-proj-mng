package main

import (
	"fmt"
	"log"

	"github.com/12ilya12/go-proj-mng/initializers"
	"github.com/12ilya12/go-proj-mng/models"
)

func init() {
	//Инициализируем конфиг
	config, err := initializers.LoadConfig()
	if err != nil {
		log.Fatal("Ошибка при загрузке переменных среды", err)
	}
	//Соединение с БД
	initializers.ConnectDB(&config)
}

func main() {
	//Категории
	categories := []models.Category{
		{Name: "Bug"},
		{Name: "Epic"},
		{Name: "Feature"},
		{Name: "Issue"},
		{Name: "Task"},
		{Name: "Test Case"},
		{Name: "User Story"},
	}

	//Статусы
	statuses := []models.Status{
		{Name: "New"},
		{Name: "Active"},
		{Name: "Resolved"},
		{Name: "Closed"},
	}

	//Пользователи
	users := []models.User{
		{
			Login:    "admin",
			Password: "11111111",
			FullName: "Админов Админ Админович",
			Email:    "admin@projmng.ru",
			Role:     "ADMIN",
		},
	}

	fmt.Println("Заполнение базы данных начальными значениями...")
	for _, category := range categories {
		err := initializers.DB.Save(&category).Error
		if err != nil {
			fmt.Printf("Ошибка при добавлении категории %s/n", category.Name)
		}
	}
	for _, status := range statuses {
		err := initializers.DB.Save(&status).Error
		if err != nil {
			fmt.Printf("Ошибка при добавлении статуса %s/n", status.Name)
		}
	}
	for _, user := range users {
		err := initializers.DB.Save(&user).Error
		if err != nil {
			fmt.Printf("Ошибка при добавлении пользователя %s/n", user.Login)
		}
	}
	fmt.Println("Успешное завершение заполнения базы данных")
}
