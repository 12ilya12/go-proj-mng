package middlewares

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/12ilya12/go-proj-mng/common"
	"github.com/12ilya12/go-proj-mng/models"
	u "github.com/12ilya12/go-proj-mng/utils"
	"github.com/dgrijalva/jwt-go"
)

func JwtAuthentication(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//Список эндпоинтов, не требующих аутентификации
		notAuth := []string{"/auth/login", "/auth/register", "/alive", "/swagger/index.html"}
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
			response = u.Message("Отсутствует аутентификационный токен")
			u.Respond(w, response, http.StatusForbidden)
			return
		}

		//Токен записывается в формате `Bearer {token-body}`.
		splitted := strings.Split(tokenHeader, " ")
		if len(splitted) != 2 {
			response = u.Message("Невалидный/некорректный токен")
			u.Respond(w, response, http.StatusForbidden)
			return
		}

		//Берём часть непосредственно с токеном
		tokenPart := splitted[1]
		claims := &models.Claims{}

		token, err := jwt.ParseWithClaims(tokenPart, claims, func(token *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("token_password")), nil
		})

		if err != nil {
			//Некорректный токен. Возвращаем 403
			response = u.Message("Malformed authentication token")
			u.Respond(w, response, http.StatusForbidden)
			return
		}

		if !token.Valid {
			//Токен не валиден. Возможно пользователь не зарегистрирован. Возвращаем 403
			response = u.Message("Токен не валиден. Возможно пользователь не зарегистрирован.")
			u.Respond(w, response, http.StatusForbidden)
			return
		}

		//Аутентификация пройдена. Продолжаем обработку запроса, добавив в контекст информацию об аутентифицированном пользователе
		ctx := context.WithValue(r.Context(), common.UserContextKey, claims.UserId)
		ctx = context.WithValue(ctx, common.RoleContextKey, claims.Role)
		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)
	})
}

func RoleBasedAccessControl(allowedRoles ...string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			role := fmt.Sprintf("%v", r.Context().Value(common.RoleContextKey))
			for _, allowedRole := range allowedRoles {
				if role == allowedRole {
					//Доступ разрешён
					next.ServeHTTP(w, r)
					return
				}
			}
			//Доступ запрещён
			http.Error(w, "Пользователю с ролью "+role+" доступ запрещён", http.StatusForbidden)
		})
	}
}
