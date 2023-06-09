package handlers

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/Edigiraldo/car-rent/internal/core/domain"
	"github.com/Edigiraldo/car-rent/internal/core/ports"
	"github.com/Edigiraldo/car-rent/internal/core/services"
	"github.com/Edigiraldo/car-rent/internal/infrastructure/driver_adapters/dtos"
	"github.com/Edigiraldo/car-rent/internal/pkg/constants"
	"github.com/Edigiraldo/car-rent/pkg/httphandler"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

var (
	ErrorInvalidTimeFrame string = "invalid time frame"
)

type Reservations struct {
	ReservationsService ports.ReservationsService
}

func NewReservations(rs ports.ReservationsService) Reservations {
	return Reservations{
		ReservationsService: rs,
	}
}

// @Summary Create a reservation
// @Description Create a reservation with the provided information
// @ID create-reservation
// @Accept json
// @Produce json
// @Param reservation body docs.ReservationRequest true "Reservation information (allowed statuses: Reserved, Canceled, Completed; allowed payment statuses: Paid, Pending, Canceled)"
// @Success 201 {object} docs.ReservationResponse "Created reservation"
// @Failure 400 {object} docs.ErrorMinimumReservationHours "Bad Request"
// @Failure 500 {object} docs.ErrorInternalServer "Internal Server Error"
// @Tags Reservations
// @Router /reservations [post]
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
			err.Error() == services.ErrCarNotAvailable ||
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

// @Summary Get a reservation
// @Description Get a reservation by UUID
// @ID get-reservation
// @Produce json
// @Param id path string true "Reservation UUID" format(uuid)
// @Success 200 {object} docs.ReservationResponse "Obtained reservation"
// @Failure 400 {object} docs.ErrorInvalidReservationStatus "Bad Request"
// @Failure 404 {object} docs.ErrorReservationNotFound "Not Found"
// @Failure 500 {object} docs.ErrorInternalServer "Internal Server Error"
// @Tags Reservations
// @Router /reservations/{id} [get]
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

// @Summary Update a reservation
// @Description Update a reservation by UUID
// @ID update-reservation
// @Accept json
// @Produce json
// @Param id path string true "Reservation UUID" format(uuid)
// @Param reservation body docs.ReservationRequest true "Reservation information (allowed statuses: Reserved, Canceled, Completed; allowed payment statuses: Paid, Pending, Canceled)"
// @Success 200 {object} docs.ReservationResponse "Updated reservation"
// @Failure 400 {object} docs.ErrorInvalidReservationTimeFrame "Bad Request"
// @Failure 404 {object} docs.ErrorReservationNotFound "Not Found"
// @Failure 500 {object} docs.ErrorInternalServer "Internal Server Error"
// @Tags Reservations
// @Router /reservations/{id} [put]
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

// @Summary Delete a reservation
// @Description Delete a reservation by UUID
// @ID delete-reservation
// @Produce json
// @Param id path string true "Reservation UUID" format(uuid)
// @Success 204 "No Content"
// @Failure 400 {object} docs.ErrorinvalidUUID "Bad Request"
// @Failure 404 {object} docs.ErrorReservationNotFound "Not Found"
// @Failure 500 {object} docs.ErrorInternalServer "Internal Server Error"
// @Tags Reservations
// @Router /reservations/{id} [delete]
func (rh Reservations) Delete(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id := params["id"]
	ID, err := uuid.Parse(id)
	if err != nil {
		httphandler.WriteErrorResponse(w, http.StatusBadRequest, ErrInvalidID)
		return
	}

	if err = rh.ReservationsService.Delete(r.Context(), ID); err != nil {
		if err.Error() == services.ErrReservationNotFound {
			httphandler.WriteErrorResponse(w, http.StatusNotFound, err.Error())
		} else {
			httphandler.WriteErrorResponse(w, http.StatusInternalServerError, ErrInternalServerError)
			log.Println(err)
		}

		return
	}

	httphandler.WriteSuccessResponse(w, http.StatusNoContent, nil)
}

// @Summary Get reservations
// @Description Get reservations
// @ID get-reservations
// @Produce json
// @Param from_reservation_id query string false "Last seen reservation" format(uuid)
// @Param start_date query string false "Star date"
// @Param end_date query string false "End date"
// @Success 200 {object} docs.Reservations "Obtained reservations"
// @Failure 400 {object} docs.ErrorInvalidTimeFrame "Bad Request"
// @Failure 500 {object} docs.ErrorInternalServer "Internal Server Error"
// @Tags Reservations
// @Router /reservations/ [get]
func (rh Reservations) List(w http.ResponseWriter, r *http.Request) {
	var startDate, endDate time.Time
	var err error

	fromReservationID := r.URL.Query().Get("from_reservation_id")
	if _, err := uuid.Parse(fromReservationID); fromReservationID != "" && err != nil {
		httphandler.WriteErrorResponse(w, http.StatusBadRequest, fmt.Sprintf("from_reservation_id: %s", err.Error()))
		return
	}
	if sDate := r.URL.Query().Get("start_date"); sDate != "" {
		if startDate, err = time.Parse(constants.Values.DATETIME_LAYOUT, sDate); err != nil {
			httphandler.WriteErrorResponse(w, http.StatusBadRequest, fmt.Sprintf("start_date: %s", err.Error()))
			return
		}
	}
	if eDate := r.URL.Query().Get("end_date"); eDate != "" {
		if endDate, err = time.Parse(constants.Values.DATETIME_LAYOUT, eDate); err != nil {
			httphandler.WriteErrorResponse(w, http.StatusBadRequest, fmt.Sprintf("end_date: %s", err.Error()))
			return
		}
	}
	if !startDate.IsZero() && !endDate.IsZero() {
		if endDate.Before(startDate) || startDate.Equal(endDate) {
			httphandler.WriteErrorResponse(w, http.StatusBadRequest, ErrorInvalidTimeFrame)
			return
		}
	}

	cars, err := rh.ReservationsService.List(r.Context(), fromReservationID, startDate, endDate)
	if err != nil {
		httphandler.WriteErrorResponse(w, http.StatusInternalServerError, ErrInternalServerError)
		log.Println(err)

		return
	}

	listReservationsResponse := getReservationsResponse(cars)

	httphandler.WriteSuccessResponse(w, http.StatusOK, listReservationsResponse)
}

// @Summary Get reservations by Car id
// @Description Get reservations by Car id
// @ID get-reservation-by-car
// @Produce json
// @Param car_id path string true "Car id" format(uuid)
// @Success 200 {object} docs.Reservations "Obtained reservations"
// @Failure 400 {object} docs.ErrorinvalidUUID "Bad Request"
// @Failure 500 {object} docs.ErrorInternalServer "Internal Server Error"
// @Tags Reservations
// @Router /cars/{car_id}/reservations [get]
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

// @Summary Get reservations by User id
// @Description Get reservations by User id
// @ID get-reservation-by-user
// @Produce json
// @Param user_id path string true "User id" format(uuid)
// @Success 200 {object} docs.Reservations "Obtained reservations"
// @Failure 400 {object} docs.ErrorinvalidUUID "Bad Request"
// @Failure 500 {object} docs.ErrorInternalServer "Internal Server Error"
// @Tags Reservations
// @Router /users/{user_id}/reservations [get]
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
