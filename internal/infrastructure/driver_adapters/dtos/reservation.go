package dtos

import (
	"encoding/json"
	"errors"
	"io"
	"time"

	"github.com/Edigiraldo/car-rent/internal/core/domain"
	"github.com/Edigiraldo/car-rent/internal/pkg/constants"
	"github.com/Edigiraldo/car-rent/pkg/utils"
	"github.com/google/uuid"
)

var (
	ErrInvalidReservationStatus = "invalid reservation status"
	ErrInvalidPaymentStatus     = "invalid payment status"
)

type Reservation struct {
	ID            uuid.UUID `json:"id,omitempty"`
	UserID        uuid.UUID `json:"user_id"`
	CarID         uuid.UUID `json:"car_id"`
	Status        string    `json:"status"`
	PaymentStatus string    `json:"payment_status"`
	StartDate     time.Time `json:"start_date"`
	EndDate       time.Time `json:"end_date"`
}

func (r Reservation) ToDomain() domain.Reservation {
	return domain.Reservation{
		ID:            r.ID,
		UserID:        r.UserID,
		CarID:         r.CarID,
		Status:        r.Status,
		PaymentStatus: r.PaymentStatus,
		StartDate:     r.StartDate,
		EndDate:       r.EndDate,
	}
}

func (r *Reservation) FromDomain(dr domain.Reservation) {
	r.ID = dr.ID
	r.UserID = dr.UserID
	r.CarID = dr.CarID
	r.Status = dr.Status
	r.PaymentStatus = dr.PaymentStatus
	r.StartDate = dr.StartDate
	r.EndDate = dr.EndDate
}

func ReservationFromBody(body io.Reader) (Reservation, error) {
	var reservation Reservation
	err := json.NewDecoder(body).Decode(&reservation)
	if err != nil {
		return Reservation{}, err
	}

	if !isValidReservationStatus(reservation.Status) {
		return Reservation{}, errors.New(ErrInvalidReservationStatus)
	}

	if !isValidPaymentStatus(reservation.PaymentStatus) {
		return Reservation{}, errors.New(ErrInvalidPaymentStatus)
	}

	return reservation, nil
}

func isValidReservationStatus(status string) bool {
	reservationStatuses := constants.Values.RESERVATION_STATUSES.Values()

	return utils.IsInSlice(reservationStatuses, status)
}

func isValidPaymentStatus(status string) bool {
	paymentStatuses := constants.Values.PAYMENT_STATUSES.Values()

	return utils.IsInSlice(paymentStatuses, status)
}
