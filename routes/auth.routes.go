package routes

import (
	"github.com/12ilya12/go-proj-mng/controllers"
	"github.com/gorilla/mux"
)

type AuthRouteController struct {
	authController controllers.AuthController
}

func NewAuthRouteController(authController controllers.AuthController) AuthRouteController {
	return AuthRouteController{authController}
}

func (rc *AuthRouteController) AuthRoute(router *mux.Router) {
	authRouter := router.PathPrefix("/auth").Subrouter()
	authRouter.HandleFunc("/register", rc.authController.Register).Methods("POST")
	authRouter.HandleFunc("/login", rc.authController.Login).Methods("POST")
}
