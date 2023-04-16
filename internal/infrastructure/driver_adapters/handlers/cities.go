package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/Edigiraldo/car-rent/internal/core/ports"
	"github.com/Edigiraldo/car-rent/internal/infrastructure/driver_adapters/dtos"
)

type Cities struct {
	CitiesService ports.CitiesService
}

func NewCities(cs ports.CitiesService) *Cities {
	return &Cities{
		CitiesService: cs,
	}
}

func (ch *Cities) ListNames(w http.ResponseWriter, r *http.Request) {
	citiesName, err := ch.CitiesService.ListNames(r.Context())
	if err != nil {
		http.Error(w, ErrInternalServerError, http.StatusInternalServerError)
		log.Println(err)

		return
	}

	response := dtos.ListCitiesNameResponse{
		CitiesName: citiesName,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
