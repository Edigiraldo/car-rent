package ports

import "net/http"

// mockgen -source=internal/core/ports/controllers.go -destination=internal/pkg/mocks/controllers.go

type HeathController interface {
	Pong(w http.ResponseWriter, r *http.Request)
}

type CarsController interface {
	Register(w http.ResponseWriter, r *http.Request)
	Get(w http.ResponseWriter, r *http.Request)
	FullUpdate(w http.ResponseWriter, r *http.Request)
	Delete(w http.ResponseWriter, r *http.Request)
	List(w http.ResponseWriter, r *http.Request)
	GetReservations(w http.ResponseWriter, r *http.Request)
}

type UsersController interface {
	SignUp(w http.ResponseWriter, r *http.Request)
	Get(w http.ResponseWriter, r *http.Request)
	FullUpdate(w http.ResponseWriter, r *http.Request)
	Delete(w http.ResponseWriter, r *http.Request)
	GetReservations(w http.ResponseWriter, r *http.Request)
}

type CitiesController interface {
	ListNames(w http.ResponseWriter, r *http.Request)
}

type ReservationsController interface {
	Book(w http.ResponseWriter, r *http.Request)
	Get(w http.ResponseWriter, r *http.Request)
	FullUpdate(w http.ResponseWriter, r *http.Request)
	Delete(w http.ResponseWriter, r *http.Request)
}
