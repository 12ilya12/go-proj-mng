package initializers

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	ServerPort     string
	DBHost         string
	DBUserName     string
	DBUserPassword string
	DBName         string
	DBPort         string
}

func LoadConfig() (config Config, err error) {
	err = godotenv.Load() //Загрузить файл .env
	if err != nil {
		fmt.Print(err)
		return
	}

	config.DBUserName = os.Getenv("db_user")
	config.DBUserPassword = os.Getenv("db_pass")
	config.DBName = os.Getenv("db_name")
	config.DBHost = os.Getenv("db_host")
	config.DBPort = os.Getenv("db_port")
	config.ServerPort = os.Getenv("server_port")

	return
}
