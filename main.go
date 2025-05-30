package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/12ilya12/go-proj-mng/app"
	"github.com/12ilya12/go-proj-mng/controllers"
	"github.com/12ilya12/go-proj-mng/initializers"
	"github.com/12ilya12/go-proj-mng/repos"
	"github.com/12ilya12/go-proj-mng/routes"
	"github.com/12ilya12/go-proj-mng/services"
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

	//Инициализация репозиториев
	UserRepo := repos.NewUserRepository(initializers.DB)

	//Инициализация сервисов
	UserService := services.NewUserService(UserRepo)
	AuthService := services.NewAuthService(UserService)

	//Инициализация контроллеров
	AuthController := controllers.NewAuthController(AuthService)
	AuthRouteController := routes.NewAuthRouteController(AuthController)

	//Заполнение роутов
	router := mux.NewRouter()
	AuthRouteController.AuthRoute(router)
	//router.HandleFunc("/auth/register", controllers.Register).Methods("POST")

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
