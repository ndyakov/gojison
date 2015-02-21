package gojison

import (
	"encoding/json"
	"net/http"
)

func Error(w http.ResponseWriter, err error, code int) {
	if code == 0 {
		code = 400
	}
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
}

func Success(w http.ResponseWriter, message string, code int) {
	if code == 0 {
		code = 200
	}
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(map[string]string{"success": message})
}
