package ports

import "net/http"

type HeathController interface {
	Pong(w http.ResponseWriter, r *http.Request)
}

type CarsController interface {
	Register(w http.ResponseWriter, r *http.Request)
}
