package httphandler

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type errorResponse struct {
	Title  string `json:"title"`
	Status int    `json:"status"`
	Detail string `json:"detail"`
}

func newErrorResponse(status int, message string) errorResponse {
	return errorResponse{
		Title:  http.StatusText(status),
		Status: status,
		Detail: message,
	}
}

func (e errorResponse) String() string {
	b, err := json.MarshalIndent(e, "", "  ")
	if err != nil {
		return fmt.Sprintf("Error marshaling error response: %v", err)
	}
	return string(b)
}
