package handlers

import "net/http"

type Health struct {
}

func NewHealth() Health {
	return Health{}
}

func (h Health) Pong(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte("pong"))
}
