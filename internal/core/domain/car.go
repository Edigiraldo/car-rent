package domain

import "github.com/google/uuid"

type Car struct {
	ID             uuid.UUID `json:"id"`
	Type           string    `json:"type"`
	Seats          int16     `json:"seats"`
	HourlyRentCost float64   `json:"hourly_rent_cost"`
	City           string    `json:"city"`
	Status         string    `json:"status"`
}
