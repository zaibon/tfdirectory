package server

import (
	"encoding/json"
	"net/http"
)

func Error(w http.ResponseWriter, code int, message string) {
	Json(w, code, map[string]string{"error": message})
}

func Json(w http.ResponseWriter, code int, payload interface{}) {
	response, err := json.Marshal(payload)
	if err != nil {
		panic("fail to marshal payload to json")
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}
