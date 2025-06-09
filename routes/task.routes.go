package routes

import (
	"github.com/12ilya12/go-proj-mng/controllers"
	"github.com/12ilya12/go-proj-mng/middlewares"
	"github.com/gorilla/mux"
)

type TaskRouteController struct {
	taskController controllers.TaskController
}

func NewTaskRouteController(taskController controllers.TaskController) TaskRouteController {
	return TaskRouteController{taskController}
}

func (rc *TaskRouteController) TaskRoute(router *mux.Router) {
	taskRouter := router.PathPrefix("/tasks").Subrouter()
	taskRouterAdminOnly := router.PathPrefix("/tasks").Subrouter()
	taskRouter.HandleFunc("/", rc.taskController.GetAll).Methods("GET")
	taskRouter.HandleFunc("/{id:[0-9]+}", rc.taskController.GetById).Methods("GET")
	taskRouterAdminOnly.HandleFunc("/", rc.taskController.Create).Methods("POST")
	taskRouter.HandleFunc("/{id:[0-9]+}", rc.taskController.Update).Methods("PATCH")
	taskRouterAdminOnly.HandleFunc("/{id:[0-9]+}", rc.taskController.Delete).Methods("DELETE")

	taskRouter.Use(middlewares.RoleBasedAccessControl("USER", "ADMIN"))
	taskRouterAdminOnly.Use(middlewares.RoleBasedAccessControl("ADMIN"))
}
