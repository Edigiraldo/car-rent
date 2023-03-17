package domain

import "github.com/google/uuid"

type Car struct {
	ID             uuid.UUID `json:"id"`
	Type           CarType   `json:"type"`
	Seats          int16     `json:"seats"`
	HourlyRentCost float64   `json:"hourly_rent_cost"`
	City           string    `json:"city"`
	Status         CarStatus `json:"status"`
}

type CarType string

const (
	Sedan     CarType = "Sedan"
	Luxury    CarType = "Luxury"
	SportsCar CarType = "Sports Car"
	Limousine CarType = "Limousine"
)

type CarStatus string

const (
	Available   CarStatus = "Available"
	Unavailable CarStatus = "Unavailable"
)
