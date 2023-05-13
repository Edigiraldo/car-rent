package docs

import (
	"time"

	"github.com/google/uuid"
)

type Reservations struct {
	Reservations []ReservationResponse `json:"reservations"`
}

type ReservationRequest struct {
	UserID        uuid.UUID `json:"user_id" example:"a29b1af4-9650-4379-8a8b-7f6c4d374e7f"`
	CarID         uuid.UUID `json:"car_id" example:"0ddac1d8-c7f2-44a6-8c7e-3d06410f7be1"`
	Status        string    `json:"status" example:"Reserved"`
	PaymentStatus string    `json:"payment_status" example:"Paid"`
	StartDate     time.Time `json:"start_date" example:"2023-05-15T10:00:00Z"`
	EndDate       time.Time `json:"end_date" example:"2023-05-16T18:00:00Z"`
}

type ReservationResponse struct {
	ID            uuid.UUID `json:"id,omitempty" example:"882dfcf8-98c9-4a25-9637-ae4564928b10"`
	UserID        uuid.UUID `json:"user_id" example:"a29b1af4-9650-4379-8a8b-7f6c4d374e7f"`
	CarID         uuid.UUID `json:"car_id" example:"0ddac1d8-c7f2-44a6-8c7e-3d06410f7be1"`
	Status        string    `json:"status" example:"Reserved"`
	PaymentStatus string    `json:"payment_status" example:"Paid"`
	StartDate     time.Time `json:"start_date" example:"2027-05-15T10:00:00Z"`
	EndDate       time.Time `json:"end_date" example:"2027-05-22T18:00:00Z"`
}
