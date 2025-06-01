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
	userCreateDto := models.UserCreate{}
	//Декодируем тело запроса в структуру dto
	dec := json.NewDecoder(r.Body)
	err := dec.Decode(&userCreateDto)
	//b, err := json.Marshal(r.Body)
	//fmt.Println(b)
	if err != nil {
		u.Respond(w, u.Message(false, "Некорректный запрос"))
		return
	}
	newUser, err := ac.authService.Register(&userCreateDto)
	if err != nil {
		u.Respond(w, u.Message(false, "Ошибка при регистрации пользователя"))
		return
	}

	//u.Respond(w, newUser)
	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(newUser)
}

func (ac *AuthController) Login(w http.ResponseWriter, r *http.Request) {
	/* account := &models.Account{}
	err := json.NewDecoder(r.Body).Decode(account) //decode the request body into struct and failed if any error occur
	if err != nil {
		u.Respond(w, u.Message(false, "Invalid request"))
		return
	}

	resp := models.Login(account.Login, account.Password)
	u.Respond(w, resp) */
}
