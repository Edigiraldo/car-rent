package handlers

import (
	"log"
	"net/http"
	"strings"

	"github.com/Edigiraldo/car-rent/internal/core/domain"
	"github.com/Edigiraldo/car-rent/internal/core/ports"
	"github.com/Edigiraldo/car-rent/internal/core/services"
	"github.com/Edigiraldo/car-rent/internal/infrastructure/driver_adapters/dtos"
	"github.com/Edigiraldo/car-rent/pkg/httphandler"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

type Reservations struct {
	ReservationsService ports.ReservationsService
}

func NewReservations(rs ports.ReservationsService) Reservations {
	return Reservations{
		ReservationsService: rs,
	}
}

func (rh Reservations) Book(w http.ResponseWriter, r *http.Request) {
	var newReservation domain.Reservation
	reservation, err := dtos.ReservationFromBody(r.Body)
	if err != nil {
		httphandler.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	if newReservation, err = rh.ReservationsService.Book(r.Context(), reservation.ToDomain()); err != nil {
		if err.Error() == services.ErrUserNotFound ||
			err.Error() == services.ErrCarNotFound ||
			err.Error() == services.ErrInvalidReservationTimeFrame ||
			strings.HasPrefix(err.Error(), services.ErrMinimumReservationHours) {
			httphandler.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}
		httphandler.WriteErrorResponse(w, http.StatusInternalServerError, ErrInternalServerError)
		log.Println(err)

		return
	}

	reservation.FromDomain(newReservation)
	httphandler.WriteSuccessResponse(w, http.StatusCreated, reservation)
}

func (rh Reservations) Get(w http.ResponseWriter, r *http.Request) {
	var reservation dtos.Reservation

	params := mux.Vars(r)
	id := params["id"]
	ID, err := uuid.Parse(id)
	if err != nil {
		httphandler.WriteErrorResponse(w, http.StatusBadRequest, ErrInvalidID)
		return
	}

	dc, err := rh.ReservationsService.Get(r.Context(), ID)
	if err != nil {
		if err.Error() == services.ErrReservationNotFound {
			httphandler.WriteErrorResponse(w, http.StatusNotFound, err.Error())
		} else {
			httphandler.WriteErrorResponse(w, http.StatusInternalServerError, ErrInternalServerError)
			log.Println(err)
		}

		return
	}

	reservation.FromDomain(dc)
	httphandler.WriteSuccessResponse(w, http.StatusOK, reservation)

}

func (rh Reservations) FullUpdate(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id := params["id"]
	ID, err := uuid.Parse(id)
	if err != nil {
		httphandler.WriteErrorResponse(w, http.StatusBadRequest, ErrInvalidID)
		return
	}

	reservation, err := dtos.ReservationFromBody(r.Body)
	if err != nil {
		httphandler.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	// Get the ID from path param
	reservation.ID = ID

	if err = rh.ReservationsService.FullUpdate(r.Context(), reservation.ToDomain()); err != nil {
		if err.Error() == services.ErrReservationNotFound {
			httphandler.WriteErrorResponse(w, http.StatusNotFound, err.Error())
			return
		} else if err.Error() == services.ErrUserNotFound ||
			err.Error() == services.ErrCarNotFound ||
			err.Error() == services.ErrInvalidReservationTimeFrame ||
			strings.HasPrefix(err.Error(), services.ErrMinimumReservationHours) ||
			err.Error() == services.ErrCarNotAvailable {
			httphandler.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		} else {
			httphandler.WriteErrorResponse(w, http.StatusInternalServerError, ErrInternalServerError)
			log.Println(err)
		}

		return
	}

	httphandler.WriteSuccessResponse(w, http.StatusOK, reservation)
}

func (rh Reservations) Delete(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id := params["id"]
	ID, err := uuid.Parse(id)
	if err != nil {
		httphandler.WriteErrorResponse(w, http.StatusBadRequest, ErrInvalidID)
		return
	}

	err = rh.ReservationsService.Delete(r.Context(), ID)
	if err != nil {
		httphandler.WriteErrorResponse(w, http.StatusInternalServerError, ErrInternalServerError)
		log.Println(err)

		return
	}

	httphandler.WriteSuccessResponse(w, http.StatusNoContent, nil)
}

func (rh Reservations) GetByCarID(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	cID := params["id"]
	carID, err := uuid.Parse(cID)
	if err != nil {
		httphandler.WriteErrorResponse(w, http.StatusBadRequest, ErrInvalidID)
		return
	}

	drs, err := rh.ReservationsService.GetByCarID(r.Context(), carID)
	if err != nil {
		httphandler.WriteErrorResponse(w, http.StatusInternalServerError, ErrInternalServerError)
		log.Println(err)

		return
	}

	reservations := getReservationsResponse(drs)
	httphandler.WriteSuccessResponse(w, http.StatusOK, reservations)
}

func (rh Reservations) GetByUserID(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	uID := params["id"]
	userID, err := uuid.Parse(uID)
	if err != nil {
		httphandler.WriteErrorResponse(w, http.StatusBadRequest, ErrInvalidID)
		return
	}

	drs, err := rh.ReservationsService.GetByUserID(r.Context(), userID)
	if err != nil {
		httphandler.WriteErrorResponse(w, http.StatusInternalServerError, ErrInternalServerError)
		log.Println(err)

		return
	}

	reservations := getReservationsResponse(drs)
	httphandler.WriteSuccessResponse(w, http.StatusOK, reservations)
}

func getReservationsResponse(domainReservations []domain.Reservation) (reservations dtos.Reservations) {
	reservations.Reservations = make([]dtos.Reservation, 0)
	for _, domainReservation := range domainReservations {
		car := dtos.Reservation{}
		car.FromDomain(domainReservation)

		reservations.Reservations = append(reservations.Reservations, car)
	}

	return reservations
}
