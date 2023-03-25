package models

import (
	"github.com/Edigiraldo/car-rent/internal/core/domain"
	"github.com/google/uuid"
)

type Car struct {
	ID             uuid.UUID `json:"id"`
	Type           string    `json:"type"`
	Seats          int16     `json:"seats"`
	HourlyRentCost float64   `json:"hourly_rent_cost"`
	City           string    `json:"city"`
	Status         string    `json:"status"`
}

type CarType string

func (c *Car) ToDomain() domain.Car {
	return domain.Car{
		ID:             c.ID,
		Type:           c.Type,
		Seats:          c.Seats,
		HourlyRentCost: c.HourlyRentCost,
		City:           c.City,
		Status:         c.Status,
	}
}

func LoadCarFromDomain(dc domain.Car) Car {
	return Car{
		ID:             dc.ID,
		Type:           dc.Type,
		Seats:          dc.Seats,
		HourlyRentCost: dc.HourlyRentCost,
		City:           dc.City,
		Status:         dc.Status,
	}

}
