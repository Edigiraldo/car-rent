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
	CityID         uuid.UUID `json:"city_id"`
	Status         string    `json:"status"`
}

type CarType string

func (c *Car) ToDomain(cityName string) domain.Car {
	return domain.Car{
		ID:             c.ID,
		Type:           c.Type,
		Seats:          c.Seats,
		HourlyRentCost: c.HourlyRentCost,
		CityName:       cityName,
		Status:         c.Status,
	}
}

func LoadCarFromDomain(dc domain.Car) Car {
	return Car{
		ID:             dc.ID,
		Type:           dc.Type,
		Seats:          dc.Seats,
		HourlyRentCost: dc.HourlyRentCost,
		Status:         dc.Status,
	}

}
