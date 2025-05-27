package initializers

import {
	"github.com/joho/godotenv"
}

type Config struct {
	ServerPort string
	DBHost	string
	DBUserName string
	DBUserPassword string
	DBName string
	DBPort string
}

func LoadConfig() (config Config, err error) {
	e := godotenv.Load() //Загрузить файл .env
	if e != nil {
		fmt.Print(e)
	}

	username := os.Getenv("db_user")
	password := os.Getenv("db_pass")
	dbName := os.Getenv("db_name")
	dbHost := os.Getenv("db_host")
}