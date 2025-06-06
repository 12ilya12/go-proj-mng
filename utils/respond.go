package utils

import (
	"encoding/json"
	"net/http"
)

func Message(message string) map[string]interface{} {
	return map[string]interface{}{"message": message}
}

func Respond(w http.ResponseWriter, data map[string]interface{}, status ...int) {
	if len(status) > 0 {
		//Лучшего способа необязательного параметра функции
		//(или параметра со значением по умолчанию) на Golang не нашёл
		w.WriteHeader(status[0])
	}
	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}
