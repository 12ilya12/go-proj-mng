package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/12ilya12/go-proj-mng/app"
	"github.com/12ilya12/go-proj-mng/controllers"
	u "github.com/12ilya12/go-proj-mng/utils"
	"github.com/gorilla/mux"
)

func main() {
	router := mux.NewRouter()

	router.HandleFunc("/auth/register", controllers.Register).Methods("POST")

	router.HandleFunc("/alive", func(w http.ResponseWriter, r *http.Request) {
		u.Respond(w, u.Message(true, "Жив, цел, Орёл!"))
	}).Methods("GET")

	//Подключаем мидлвар для аутентификации по JWT
	router.Use(app.JwtAuthentication)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}

	fmt.Println(port)

	err := http.ListenAndServe(":"+port, router)
	if err != nil {
		fmt.Print(err)
	}
}
