package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/12ilya12/go-proj-mng/models"
	"github.com/12ilya12/go-proj-mng/services"
	u "github.com/12ilya12/go-proj-mng/utils"
)

type AuthController struct {
	authService services.AuthService
}

func NewAuthController(authService services.AuthService) AuthController {
	return AuthController{authService}
}

func (ac *AuthController) Register(w http.ResponseWriter, r *http.Request) {
	userDto := models.User{}
	//Декодируем тело запроса в структуру dto
	err := json.NewDecoder(r.Body).Decode(&userDto)
	if err != nil {
		u.Respond(w, u.Message("Некорректный запрос"), http.StatusBadRequest)
		return
	}
	//TODO: Валидация данных пользователя
	err = ac.authService.Register(&userDto)
	if err != nil {
		//TODO: Статус ответа должен определяться в зависимости от ошибки.
		//Не факт, что проблема в запросе. Например, могут быть проблемы с БД.
		u.Respond(w, u.Message(err.Error()), http.StatusBadRequest)
		return
	}

	//u.Respond(w, newUser)
	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(userDto)
}

func (ac *AuthController) Login(w http.ResponseWriter, r *http.Request) {
	userDto := models.User{}
	err := json.NewDecoder(r.Body).Decode(&userDto)
	if err != nil {
		u.Respond(w, u.Message("Некорректный запрос. Ошибка: "+err.Error()), http.StatusBadRequest)
		return
	}

	token, err := ac.authService.Login(userDto.Login, userDto.Password)
	if err != nil {
		//w.WriteHeader(http.StatusUnauthorized)
		u.Respond(w, map[string]interface{}{"error": err.Error()}, http.StatusUnauthorized)
		return
	}
	u.Respond(w, map[string]interface{}{"access_token": token})
}
