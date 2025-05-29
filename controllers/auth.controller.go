package controllers

import (
	//"encoding/json"
	"net/http"
	//"strings"
	//"time"

	//"github.com/12ilya12/go-proj-mng/initializers"
	//"github.com/12ilya12/go-proj-mng/models"
	//u "github.com/12ilya12/go-proj-mng/utils"
	"gorm.io/gorm"
)

type AuthController struct {
	DB *gorm.DB
}

func NewAuthController(DB *gorm.DB) AuthController {
	return AuthController{DB}
}

func (ac *AuthController) Register(w http.ResponseWriter, r *http.Request) {
	/* account := &models.Account{}
	err := json.NewDecoder(r.Body).Decode(account) //decode the request body into struct and failed if any error occur
	if err != nil {
		u.Respond(w, u.Message(false, "Invalid request"))
		return
	}

	resp := account.Create() //Create account
	u.Respond(w, resp) */
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
