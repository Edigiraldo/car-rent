package domain

import (
	"time"

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
