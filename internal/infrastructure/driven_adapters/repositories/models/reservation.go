package models

import (
	"time"

	"github.com/Edigiraldo/car-rent/internal/core/domain"
	"github.com/google/uuid"
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

func (r Reservation) ToDomain(cityName string) domain.Reservation {
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

func LoadReservationFromDomain(dr domain.Reservation) Reservation {
	return Reservation{
		ID:            dr.ID,
		UserID:        dr.UserID,
		CarID:         dr.CarID,
		Status:        dr.Status,
		PaymentStatus: dr.PaymentStatus,
		StartDate:     dr.StartDate,
		EndDate:       dr.EndDate,
	}

}
