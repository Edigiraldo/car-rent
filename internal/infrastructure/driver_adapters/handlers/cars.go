package handlers

import (
	"log"
	"net/http"

	"github.com/Edigiraldo/car-rent/internal/core/domain"
	"github.com/Edigiraldo/car-rent/internal/core/ports"
	"github.com/Edigiraldo/car-rent/internal/core/services"
	"github.com/Edigiraldo/car-rent/internal/infrastructure/driver_adapters/dtos"
	"github.com/Edigiraldo/car-rent/pkg/httphandler"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

var (
	ErrInvalidID           = "id could not be converted to uuid"
	ErrInternalServerError = "internal server error"
	ErrCityQueryParamEmpty = "city query param can not be empty"
)

type Cars struct {
	CarsService ports.CarsService
}

func NewCars(cs ports.CarsService) Cars {
	return Cars{
		CarsService: cs,
	}
}

// @Summary Register a new car
// @Description Register a new car with the provided information
// @ID register-car
// @Accept json
// @Produce json
// @Param car body docs.CarRequest true "Car information (allowed types: Sedan, Luxury, Sports Car, Limousine; allowed statuses: Available, Unavailable)"
// @Success 201 {object} docs.CarResponse "Created car"
// @Failure 400 {object} docs.ErrorInvalidCityName "Bad Request"
// @Failure 500 {object} docs.ErrorInternalServer "Internal Server Error"
// @Tags Cars
// @Router /cars [post]
func (ch Cars) Register(w http.ResponseWriter, r *http.Request) {
	var newCar domain.Car
	car, err := dtos.CarFromBody(r.Body)
	if err != nil {
		httphandler.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	if newCar, err = ch.CarsService.Register(r.Context(), car.ToDomain()); err != nil {
		if err.Error() == services.ErrInvalidCityName {
			httphandler.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		log.Println(err)
		httphandler.WriteErrorResponse(w, http.StatusInternalServerError, ErrInternalServerError)

		return
	}

	car.FromDomain(newCar)

	httphandler.WriteSuccessResponse(w, http.StatusCreated, car)
}

// @Summary Get a car
// @Description Get a car by UUID
// @ID get-car
// @Produce json
// @Param id path string true "Car UUID" format(uuid)
// @Success 200 {object} docs.CarResponse "Obtained car"
// @Failure 400 {object} docs.ErrorinvalidUUID "Bad Request"
// @Failure 404 {object} docs.ErrorCarNotFound "Not Found"
// @Failure 500 {object} docs.ErrorInternalServer "Internal Server Error"
// @Tags Cars
// @Router /cars/{id} [get]
func (ch Cars) Get(w http.ResponseWriter, r *http.Request) {
	var car dtos.Car

	params := mux.Vars(r)
	id := params["id"]
	ID, err := uuid.Parse(id)
	if err != nil {
		httphandler.WriteErrorResponse(w, http.StatusBadRequest, ErrInvalidID)
		return
	}

	dc, err := ch.CarsService.Get(r.Context(), ID)
	if err != nil {
		if err.Error() == services.ErrCarNotFound {
			httphandler.WriteErrorResponse(w, http.StatusNotFound, err.Error())
		} else {
			httphandler.WriteErrorResponse(w, http.StatusInternalServerError, ErrInternalServerError)
			log.Println(err)
		}

		return
	}

	car.FromDomain(dc)

	httphandler.WriteSuccessResponse(w, http.StatusOK, car)
}

// @Summary Update a car
// @Description Update a car by UUID
// @ID update-car
// @Accept json
// @Produce json
// @Param id path string true "Car UUID" format(uuid)
// @Param car body docs.CarRequest true "Car information (allowed types: Sedan, Luxury, Sports Car, Limousine; allowed statuses: Available, Unavailable)"
// @Success 200 {object} docs.CarResponse "Updated car"
// @Failure 400 {object} docs.ErrorInvalidCarStatus "Bad Request"
// @Failure 404 {object} docs.ErrorCarNotFound "Not Found"
// @Failure 500 {object} docs.ErrorInternalServer "Internal Server Error"
// @Tags Cars
// @Router /cars/{id} [put]
func (ch Cars) FullUpdate(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id := params["id"]
	ID, err := uuid.Parse(id)
	if err != nil {
		httphandler.WriteErrorResponse(w, http.StatusBadRequest, ErrInvalidID)
		return
	}

	car, err := dtos.CarFromBody(r.Body)
	if err != nil {
		httphandler.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	// Get the ID from path param
	car.ID = ID

	if err = ch.CarsService.FullUpdate(r.Context(), car.ToDomain()); err != nil {
		if err.Error() == services.ErrCarNotFound {
			httphandler.WriteErrorResponse(w, http.StatusNotFound, err.Error())
		} else if err.Error() == services.ErrInvalidCityName {
			httphandler.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
		} else {
			httphandler.WriteErrorResponse(w, http.StatusInternalServerError, ErrInternalServerError)
			log.Println(err)
		}

		return
	}

	httphandler.WriteSuccessResponse(w, http.StatusOK, car)
}

// @Summary Delete a car
// @Description Delete a car by UUID
// @ID delete-car
// @Produce json
// @Param id path string true "Car UUID" format(uuid)
// @Success 204 "No Content"
// @Failure 400 {object} docs.ErrorinvalidUUID "Bad Request"
// @Failure 404 {object} docs.ErrorCarNotFound "Not Found"
// @Failure 500 {object} docs.ErrorInternalServer "Internal Server Error"
// @Tags Cars
// @Router /cars/{id} [delete]
func (ch Cars) Delete(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id := params["id"]
	ID, err := uuid.Parse(id)
	if err != nil {
		httphandler.WriteErrorResponse(w, http.StatusBadRequest, ErrInvalidID)
		return
	}

	if err = ch.CarsService.Delete(r.Context(), ID); err != nil {
		if err.Error() == services.ErrCarNotFound {
			httphandler.WriteErrorResponse(w, http.StatusNotFound, err.Error())
		} else {
			httphandler.WriteErrorResponse(w, http.StatusInternalServerError, ErrInternalServerError)
			log.Println(err)
		}

		return
	}
	httphandler.WriteSuccessResponse(w, http.StatusNoContent, nil)
}

// @Summary List cars
// @Description Lists cars from a city in pages of 20 elements. from_car_id parameter
// @Description is taken as the last seen car in a previous page.
// @ID list-cars
// @Produce json
// @Param city query string true "City name"
// @Param from_car_id query string false "Last seen car ID" format(uuid)
// @Success 200 {object} docs.ListCarsResponse "Obtained car"
// @Failure 400 {object} docs.ErrorCityQueryParamEmpty "Bad Request"
// @Failure 500 {object} docs.ErrorInternalServer "Internal Server Error"
// @Tags Cars
// @Router /cars/ [get]
func (ch Cars) List(w http.ResponseWriter, r *http.Request) {
	city := r.URL.Query().Get("city")
	if city == "" {
		httphandler.WriteErrorResponse(w, http.StatusBadRequest, ErrCityQueryParamEmpty)
		return
	}
	from_car_id := r.URL.Query().Get("from_car_id")
	if _, err := uuid.Parse(from_car_id); err != nil && from_car_id != "" {
		httphandler.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	cars, err := ch.CarsService.List(r.Context(), city, from_car_id)
	if err != nil && err.Error() != services.ErrInvalidCityName {
		httphandler.WriteErrorResponse(w, http.StatusInternalServerError, ErrInternalServerError)
		log.Println(err)

		return
	}

	listCarsResponse := getListCarsResponse(cars)

	httphandler.WriteSuccessResponse(w, http.StatusOK, listCarsResponse)
}

// Gets a list of cars and builds the user response
func getListCarsResponse(cars []domain.Car) (listCarsResponse dtos.ListCarsResponse) {
	listCarsResponse.Cars = make([]dtos.Car, 0)
	for _, domainCar := range cars {
		car := dtos.Car{}
		car.FromDomain(domainCar)

		listCarsResponse.Cars = append(listCarsResponse.Cars, car)
	}

	return listCarsResponse
}
