package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/12ilya12/go-proj-mng/app"
	"github.com/12ilya12/go-proj-mng/controllers"
	"github.com/12ilya12/go-proj-mng/initializers"
	u "github.com/12ilya12/go-proj-mng/utils"
	"github.com/gorilla/mux"
)

func main() {
	//Инициализируем конфиг
	config, err := initializers.LoadConfig()
	if err != nil {
		log.Fatal("Ошибка при загрузке переменных среды", err)
	}

	//Соединение с БД
	initializers.ConnectDB(&config)

	router := mux.NewRouter()

	router.HandleFunc("/auth/register", controllers.Register).Methods("POST")

	router.HandleFunc("/alive", func(w http.ResponseWriter, r *http.Request) {
		u.Respond(w, u.Message(true, "Жив, цел, Орёл!"))
	}).Methods("GET")

	//Подключаем мидлвар для аутентификации по JWT
	router.Use(app.JwtAuthentication)

	if config.ServerPort == "" {
		config.ServerPort = "8000"
	}

	fmt.Println(config.ServerPort)

	err = http.ListenAndServe(":"+config.ServerPort, router)
	if err != nil {
		fmt.Print(err)
	}
}
