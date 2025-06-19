package controllers

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/12ilya12/go-proj-mng/models"
	"github.com/12ilya12/go-proj-mng/services"
	u "github.com/12ilya12/go-proj-mng/utils"
	"github.com/go-playground/validator"
)

type AuthController struct {
	authService services.AuthService
	vld         *validator.Validate
}

func NewAuthController(authService services.AuthService) AuthController {
	vld := validator.New()
	return AuthController{authService, vld}
}

func (ac *AuthController) Register(w http.ResponseWriter, r *http.Request) {
	userDto := models.User{}
	//Декодируем тело запроса в структуру dto
	err := json.NewDecoder(r.Body).Decode(&userDto)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = ac.vld.Struct(userDto)
	if err != nil {
		var validationErrors validator.ValidationErrors
		if errors.As(err, &validationErrors) {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	}

	err = ac.authService.Register(&userDto)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadGateway)
		return
	}
	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(userDto)
}

func (ac *AuthController) Login(w http.ResponseWriter, r *http.Request) {
	userDto := models.User{}
	err := json.NewDecoder(r.Body).Decode(&userDto)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = ac.vld.Struct(userDto)
	if err != nil {
		var validationErrors validator.ValidationErrors
		if errors.As(err, &validationErrors) {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	}

	token, err := ac.authService.Login(userDto.Login, userDto.Password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}
	u.Respond(w, map[string]interface{}{"access_token": token})
}
