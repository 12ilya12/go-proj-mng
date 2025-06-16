package routes

import (
	"github.com/12ilya12/go-proj-mng/controllers"
	"github.com/12ilya12/go-proj-mng/middlewares"
	"github.com/gorilla/mux"
)

type DependencyRouteController struct {
	dependencyController controllers.DependencyController
}

func NewDependencyRouteController(dependencyController controllers.DependencyController) DependencyRouteController {
	return DependencyRouteController{dependencyController}
}

func (rc *DependencyRouteController) DependencyRoute(router *mux.Router) {
	dependencyRouter := router.PathPrefix("/tasks/{taskId:[0-9]+}/dependencies").Subrouter()
	dependencyRouterAdminOnly := router.PathPrefix("/tasks/{taskId:[0-9]+}/dependencies").Subrouter()
	dependencyRouter.HandleFunc("", rc.dependencyController.Get).Methods("GET")
	dependencyRouter.HandleFunc("", rc.dependencyController.Create).Methods("POST")
	dependencyRouter.HandleFunc("/{dependencyId:[0-9]+}", rc.dependencyController.Delete).Methods("DELETE")

	dependencyRouter.Use(middlewares.RoleBasedAccessControl("USER", "ADMIN"))
	dependencyRouterAdminOnly.Use(middlewares.RoleBasedAccessControl("ADMIN"))
}
