package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/Edigiraldo/car-rent/internal/core/domain"
	"github.com/Edigiraldo/car-rent/internal/core/ports"
	"github.com/Edigiraldo/car-rent/internal/core/services"
	"github.com/Edigiraldo/car-rent/internal/infrastructure/driver_adapters/dtos"
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
