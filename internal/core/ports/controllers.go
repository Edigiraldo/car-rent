package ports

import "net/http"

// mockgen -source=internal/core/ports/controllers.go -destination=internal/pkg/mocks/controllers.go

type HeathController interface {
	Pong(w http.ResponseWriter, r *http.Request)
}

type CarsController interface {
	Register(w http.ResponseWriter, r *http.Request)
	Get(w http.ResponseWriter, r *http.Request)
}
