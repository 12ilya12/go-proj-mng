package initializers

import {
	"fmt"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
}

var DB *gorm.DB

func ConnectDB(config *Config) {
	var err error
	
	//Параметры подключения
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable")

	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("При соединении с базой произошла ошибка")
	}
	fmt.Println("Соединение с базой произошло успешно")
}