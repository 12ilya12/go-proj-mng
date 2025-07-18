package routes

import (
	"github.com/12ilya12/go-proj-mng/controllers"
	"github.com/12ilya12/go-proj-mng/middlewares"
	"github.com/gorilla/mux"
)

type StatusRouteController struct {
	statusController controllers.StatusController
}

func NewStatusRouteController(statusController controllers.StatusController) StatusRouteController {
	return StatusRouteController{statusController}
}

func (rc *StatusRouteController) StatusRoute(router *mux.Router) {
	statusRouter := router.PathPrefix("/statuses").Subrouter()
	statusRouterAdminOnly := router.PathPrefix("/statuses").Subrouter()
	statusRouter.HandleFunc("", rc.statusController.GetAll).Methods("GET")
	statusRouter.HandleFunc("/{id:[0-9]+}", rc.statusController.GetById).Methods("GET")
	statusRouterAdminOnly.HandleFunc("", rc.statusController.Create).Methods("POST")
	statusRouterAdminOnly.HandleFunc("/{id:[0-9]+}", rc.statusController.Update).Methods("PATCH")
	statusRouterAdminOnly.HandleFunc("/{id:[0-9]+}", rc.statusController.Delete).Methods("DELETE")

	statusRouter.Use(middlewares.RoleBasedAccessControl("USER", "ADMIN"))
	statusRouterAdminOnly.Use(middlewares.RoleBasedAccessControl("ADMIN"))
}
