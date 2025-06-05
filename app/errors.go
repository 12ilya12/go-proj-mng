package app

import (
	"net/http"

	u "github.com/12ilya12/go-proj-mng/utils"
)

func NotFoundHandler(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		u.Respond(w, u.Message("This resources was not found on our server"), http.StatusNotFound)
		next.ServeHTTP(w, r)
	})
}
