package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/Edigiraldo/car-rent/internal/core/domain"
	"github.com/Edigiraldo/car-rent/internal/core/ports"
	"github.com/Edigiraldo/car-rent/internal/core/services"
	"github.com/Edigiraldo/car-rent/internal/infrastructure/driver_adapters/dtos"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

type Reservations struct {
	ReservationsService ports.ReservationsService
}

func NewReservations(rs ports.ReservationsService) *Reservations {
	return &Reservations{
		ReservationsService: rs,
	}
}

func (rh *Reservations) Book(w http.ResponseWriter, r *http.Request) {
	var newReservation domain.Reservation
	reservation, err := dtos.ReservationFromBody(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)

		return
	}

	if newReservation, err = rh.ReservationsService.Book(r.Context(), reservation.ToDomain()); err != nil {
		if err.Error() == services.ErrUserNotFound {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		} else if err.Error() == services.ErrCarNotFound {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		log.Println(err)
		http.Error(w, ErrInternalServerError, http.StatusInternalServerError)

		return
	}

	reservation.FromDomain(newReservation)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(reservation)
}

func (rh *Reservations) Get(w http.ResponseWriter, r *http.Request) {
	var reservation dtos.Reservation

	params := mux.Vars(r)
	id := params["id"]
	ID, err := uuid.Parse(id)
	if err != nil {
		http.Error(w, ErrInvalidID, http.StatusBadRequest)

		return
	}

	dc, err := rh.ReservationsService.Get(r.Context(), ID)
	if err != nil {
		if err.Error() == services.ErrReservationNotFound {
			http.Error(w, err.Error(), http.StatusNotFound)
		} else {
			http.Error(w, ErrInternalServerError, http.StatusInternalServerError)
			log.Println(err)
		}

		return
	}

	reservation.FromDomain(dc)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(reservation)
}

func (rh *Reservations) FullUpdate(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id := params["id"]
	ID, err := uuid.Parse(id)
	if err != nil {
		http.Error(w, ErrInvalidID, http.StatusBadRequest)

		return
	}

	reservation, err := dtos.ReservationFromBody(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)

		return
	}

	// Get the ID from path param
	reservation.ID = ID

	if err = rh.ReservationsService.FullUpdate(r.Context(), reservation.ToDomain()); err != nil {
		if err.Error() == services.ErrReservationNotFound ||
			err.Error() == services.ErrUserNotFound ||
			err.Error() == services.ErrCarNotFound {
			http.Error(w, err.Error(), http.StatusNotFound)
			return

		} else {
			http.Error(w, ErrInternalServerError, http.StatusInternalServerError)
			log.Println(err)
		}

		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(reservation)
}
