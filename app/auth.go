package app

import (
	"context"
	"net/http"
	"os"
	"strings"

	"github.com/12ilya12/go-proj-mng/models"
	u "github.com/12ilya12/go-proj-mng/utils"
	"github.com/dgrijalva/jwt-go"
)

func JwtAuthentication(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//Список эндпоинтов, не требующих аутентификации
		notAuth := []string{"/api/user/new", "/api/user/login"}
		//Путь запроса
		requestPath := r.URL.Path

		//Пропускаем аутентификацию, если эндпоинт в белом списке
		for _, value := range notAuth {

			if value == requestPath {
				next.ServeHTTP(w, r)
				return
			}
		}

		response := make(map[string]interface{})
		//Извлекаем токен из запроса
		tokenHeader := r.Header.Get("Authorization")

		if tokenHeader == "" {
			//Токена нет в запросе, возвращаем 403
			response = u.Message(false, "Missing auth token")
			w.WriteHeader(http.StatusForbidden)
			w.Header().Add("Content-Type", "application/json")
			u.Respond(w, response)
			return
		}

		//Токен записывается в формате `Bearer {token-body}`.
		splitted := strings.Split(tokenHeader, " ")
		if len(splitted) != 2 {
			response = u.Message(false, "Invalid/Malformed auth token")
			w.WriteHeader(http.StatusForbidden)
			w.Header().Add("Content-Type", "application/json")
			u.Respond(w, response)
			return
		}

		//Берём часть непосредственно с токеном
		tokenPart := splitted[1]
		tk := &models.Token{}

		token, err := jwt.ParseWithClaims(tokenPart, tk, func(token *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("token_password")), nil
		})

		if err != nil {
			//Некорректный токен. Возвращаем 403
			response = u.Message(false, "Malformed authentication token")
			w.WriteHeader(http.StatusForbidden)
			w.Header().Add("Content-Type", "application/json")
			u.Respond(w, response)
			return
		}

		if !token.Valid {
			//Токен не валиден. Возможно пользователь не зарегистрирован. Возвращаем 403
			response = u.Message(false, "Token is not valid.")
			w.WriteHeader(http.StatusForbidden)
			w.Header().Add("Content-Type", "application/json")
			u.Respond(w, response)
			return
		}

		//Аутентификация пройдена. Продолжаем обработку запроса, добавив в контекст информацию об аутентифицированном пользователя
		//fmt.Sprintf("User %", tk.UserId)
		ctx := context.WithValue(r.Context(), "user", tk.UserId)
		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)
	})
}
