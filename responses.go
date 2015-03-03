package gojison

import (
	"encoding/json"
	"net/http"
)

// Error recieves http.ResponseWriter, error interface and
// response code. Sets the code to the ResponseWriter header
// and then writes out the error as JSON.
//
// If the code parameter is 0, http.StatusBadRequest will be used.
func Error(w http.ResponseWriter, err error, code int) {
	if code == 0 {
		code = http.StatusBadRequest
	}
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
}

// Success recieves http.ResponseWriter, success message and a
// response code. Sets the code to the ResponseWriter header
// and then writes out the message as JSON.
//
// If the code parameter is 0, http.StatusOK will be used.
func Success(w http.ResponseWriter, message string, code int) {
	if code == 0 {
		code = http.StatusOK
	}
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(map[string]string{"success": message})
}
