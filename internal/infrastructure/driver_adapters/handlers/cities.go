package handlers

import (
	"log"
	"net/http"

	"github.com/Edigiraldo/car-rent/internal/core/ports"
	"github.com/Edigiraldo/car-rent/internal/infrastructure/driver_adapters/dtos"
	"github.com/Edigiraldo/car-rent/pkg/httphandler"
)

type Cities struct {
	CitiesService ports.CitiesService
}

func NewCities(cs ports.CitiesService) Cities {
	return Cities{
		CitiesService: cs,
	}
}

func (ch Cities) ListNames(w http.ResponseWriter, r *http.Request) {
	citiesName, err := ch.CitiesService.ListNames(r.Context())
	if err != nil {
		httphandler.WriteErrorResponse(w, http.StatusInternalServerError, ErrInternalServerError)
		log.Println(err)

		return
	}

	response := dtos.ListCitiesNameResponse{
		CitiesName: citiesName,
	}

	httphandler.WriteSuccessResponse(w, http.StatusOK, response)
}
