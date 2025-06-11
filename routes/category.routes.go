package routes

import (
	"github.com/12ilya12/go-proj-mng/controllers"
	"github.com/12ilya12/go-proj-mng/middlewares"
	"github.com/gorilla/mux"
)

type CategoryRouteController struct {
	categoryController controllers.CategoryController
}

func NewCategoryRouteController(categoryController controllers.CategoryController) CategoryRouteController {
	return CategoryRouteController{categoryController}
}

func (rc *CategoryRouteController) CategoryRoute(router *mux.Router) {
	categoryRouter := router.PathPrefix("/categories").Subrouter()
	categoryRouterAdminOnly := router.PathPrefix("/categories").Subrouter()
	categoryRouter.HandleFunc("", rc.categoryController.GetAll).Methods("GET")
	categoryRouter.HandleFunc("/{id:[0-9]+}", rc.categoryController.GetById).Methods("GET")
	categoryRouterAdminOnly.HandleFunc("", rc.categoryController.Create).Methods("POST")
	categoryRouterAdminOnly.HandleFunc("/{id:[0-9]+}", rc.categoryController.Update).Methods("PATCH")
	categoryRouter.HandleFunc("/{id:[0-9]+}/has-tasks", rc.categoryController.HasTasks).Methods("GET")
	categoryRouterAdminOnly.HandleFunc("/{id:[0-9]+}", rc.categoryController.Delete).Methods("DELETE")
	categoryRouterAdminOnly.HandleFunc("/{id:[0-9]+}/force", rc.categoryController.DeleteForce).Methods("DELETE")

	categoryRouter.Use(middlewares.RoleBasedAccessControl("USER", "ADMIN"))
	categoryRouterAdminOnly.Use(middlewares.RoleBasedAccessControl("ADMIN"))
}
