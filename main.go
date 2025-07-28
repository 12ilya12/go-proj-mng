package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/12ilya12/go-proj-mng/controllers"
	_ "github.com/12ilya12/go-proj-mng/docs"
	"github.com/12ilya12/go-proj-mng/initializers"
	"github.com/12ilya12/go-proj-mng/middlewares"
	"github.com/12ilya12/go-proj-mng/repos"
	"github.com/12ilya12/go-proj-mng/routes"
	"github.com/12ilya12/go-proj-mng/services"
	u "github.com/12ilya12/go-proj-mng/utils"
	"github.com/gorilla/mux"
	httpSwagger "github.com/swaggo/http-swagger"
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
	StatusRepo := repos.NewStatusRepositoryImpl(initializers.DB)
	CategoryRepo := repos.NewCategoryRepositoryImpl(initializers.DB)
	TaskRepo := repos.NewTaskRepository(initializers.DB)
	DependencyRepo := repos.NewDependencyRepository(initializers.DB)

	//Инициализация клиента напоминаний
	ReminderClient := initializers.InitReminder()

	//Инициализация сервисов
	UserService := services.NewUserService(UserRepo)
	AuthService := services.NewAuthService(UserService)
	StatusService := services.NewStatusServiceImpl(StatusRepo)
	CategoryService := services.NewCategoryServiceImpl(CategoryRepo)
	TaskService := services.NewTaskService(TaskRepo, StatusRepo, CategoryRepo, UserRepo, ReminderClient)
	DependencyService := services.NewDependencyService(DependencyRepo, TaskRepo)

	//Инициализация контроллеров
	AuthController := controllers.NewAuthController(AuthService)
	AuthRouteController := routes.NewAuthRouteController(AuthController)
	StatusController := controllers.NewStatusController(StatusService)
	StatusRouteController := routes.NewStatusRouteController(StatusController)
	CategoryController := controllers.NewCategoryController(CategoryService)
	CategoryRouteController := routes.NewCategoryRouteController(CategoryController)
	TaskController := controllers.NewTaskController(TaskService)
	TaskRouteController := routes.NewTaskRouteController(TaskController)
	DependencyController := controllers.NewDependencyController(DependencyService)
	DependencyRouteController := routes.NewDependencyRouteController(DependencyController)

	//Заполнение роутов
	router := mux.NewRouter()
	AuthRouteController.AuthRoute(router)
	StatusRouteController.StatusRoute(router)
	CategoryRouteController.CategoryRoute(router)
	TaskRouteController.TaskRoute(router)
	DependencyRouteController.DependencyRoute(router)

	//Тестовый роут
	router.HandleFunc("/alive", func(w http.ResponseWriter, r *http.Request) {
		u.Respond(w, u.Message("Жив, цел, Орёл!"))
	}).Methods("GET")

	router.PathPrefix("/doc/").Handler(httpSwagger.WrapHandler)

	//Подключаем мидлвар для аутентификации по JWT
	router.Use(middlewares.JwtAuthentication)

	//Подключаем мидлвар для аудита действий пользователя
	//router.Use(middlewares.AuditMiddleware)

	if config.ServerPort == "" {
		config.ServerPort = "8000"
	}

	fmt.Println(config.ServerPort)

	err = http.ListenAndServe(":"+config.ServerPort, router)
	if err != nil {
		fmt.Print(err)
	}
}
