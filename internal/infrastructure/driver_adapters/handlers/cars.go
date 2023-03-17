package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Edigiraldo/car-rent/internal/core/ports"
	"github.com/Edigiraldo/car-rent/internal/infrastructure/driver_adapters/handlers/dtos"
)

type Cars struct {
	CarsService ports.CarsService
}

func NewCars(cs ports.CarsService) *Cars {
	return &Cars{
		CarsService: cs,
	}
}

func (ch *Cars) Register(w http.ResponseWriter, r *http.Request) {
	car, err := dtos.CarFromBody(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	if err = ch.CarsService.Register(r.Context(), car.ToDomain()); err != nil {
		fmt.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(car)
}
