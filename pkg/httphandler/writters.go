package httphandler

import (
	"encoding/json"
	"net/http"
)

func WriteErrorResponse(w http.ResponseWriter, status int, message string) {
	errorResponse := newErrorResponse(status, message)
	body, err := json.Marshal(errorResponse)
	if err != nil {
		http.Error(w, "error generating error response", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(body)
}

func WriteSuccessResponse(w http.ResponseWriter, status int, dto interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if dto != nil {
		json.NewEncoder(w).Encode(dto)
	}
}
