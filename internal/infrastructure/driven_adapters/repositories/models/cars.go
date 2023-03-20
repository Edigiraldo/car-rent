package models

import (
	"github.com/Edigiraldo/car-rent/internal/core/domain"
	"github.com/google/uuid"
)

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

func (c *Car) ToDomain() domain.Car {
	return domain.Car{
		ID:             c.ID,
		Type:           domain.CarType(c.Type),
		Seats:          c.Seats,
		HourlyRentCost: c.HourlyRentCost,
		City:           c.City,
		Status:         domain.CarStatus(c.Status),
	}
}

func LoadCarFromDomain(dc domain.Car) Car {
	return Car{
		ID:             dc.ID,
		Type:           CarType(dc.Type),
		Seats:          dc.Seats,
		HourlyRentCost: dc.HourlyRentCost,
		City:           dc.City,
		Status:         CarStatus(dc.Status),
	}

}
