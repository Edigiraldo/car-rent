package handlers

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/Edigiraldo/car-rent/internal/core/domain"
	"github.com/Edigiraldo/car-rent/internal/core/services"
	"github.com/Edigiraldo/car-rent/internal/infrastructure/driver_adapters/handlers/dtos"
	mocks "github.com/Edigiraldo/car-rent/internal/pkg/mocks"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

type carsDependencies struct {
	carsService *mocks.MockCarsService
}

func NewCarsDependencies(carsSrv *mocks.MockCarsService) *carsDependencies {
	return &carsDependencies{
		carsService: carsSrv,
	}
}

func TestCarsRegister(t *testing.T) {
	car := dtos.Car{
		Type:           "Sedan",
		Seats:          4,
		HourlyRentCost: 21.1,
		City:           "Los Angeles",
		Status:         "Available",
	}

	type args struct {
		car dtos.Car
	}
	type wants struct {
		statusCode int
	}
	tests := []struct {
		name     string
		args     args
		wants    wants
		setMocks func(*carsDependencies)
	}{
		{
			name: "returns status code 201 when body is appropriate",
			args: args{
				car: car,
			},
			wants: wants{
				statusCode: http.StatusCreated,
			},
			setMocks: func(d *carsDependencies) {
				d.carsService.EXPECT().Register(gomock.Any(), car.ToDomain()).Return(domain.Car{}, nil)
			},
		},
		{
			name: "returns 400 status code when car type in body is not allowed",
			args: args{
				car: dtos.Car{
					Type: "Seda",
				},
			},
			wants: wants{
				statusCode: http.StatusBadRequest,
			},
			setMocks: func(d *carsDependencies) {
			},
		},
		{
			name: "returns 500 status code when car service fails to register car",
			args: args{
				car: car,
			},
			wants: wants{
				statusCode: http.StatusInternalServerError,
			},
			setMocks: func(d *carsDependencies) {
				d.carsService.EXPECT().Register(gomock.Any(), car.ToDomain()).Return(domain.Car{}, errors.New("error registering car"))
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			mockCtlr := gomock.NewController(t)
			carsSrv := mocks.NewMockCarsService(mockCtlr)
			d := NewCarsDependencies(carsSrv)
			test.setMocks(d)

			baseURL := "/api/v1/"
			body, _ := json.Marshal(test.args.car)
			URL := baseURL + "cars/"
			req, err := http.NewRequest(http.MethodGet, URL, bytes.NewBuffer(body))
			if err != nil {
				t.Fatal(err)
			}

			rr := httptest.NewRecorder()

			carsHandler := NewCars(carsSrv)
			carsHandler.Register(rr, req)

			assert.Equal(t, test.wants.statusCode, rr.Code)
		})
	}
}

func TestCarsGet(t *testing.T) {
	car := dtos.Car{
		ID:             uuid.New(),
		Type:           "Sedan",
		Seats:          4,
		HourlyRentCost: 21.1,
		City:           "Los Angeles",
		Status:         "Available",
	}

	type args struct {
		requestID string
	}
	type wants struct {
		statusCode int
	}
	tests := []struct {
		name     string
		args     args
		wants    wants
		setMocks func(*carsDependencies)
	}{
		{
			name: "returns status code 200 when id is properly set and car was found",
			args: args{
				requestID: car.ID.String(),
			},
			wants: wants{
				statusCode: http.StatusOK,
			},
			setMocks: func(d *carsDependencies) {
				d.carsService.EXPECT().Get(gomock.Any(), car.ID).Return(car.ToDomain(), nil)
			},
		},
		{
			name: "returns 400 status code when path param id is not an uuid",
			args: args{
				requestID: "this-is-an-not-uuid",
			},
			wants: wants{
				statusCode: http.StatusBadRequest,
			},
			setMocks: func(d *carsDependencies) {
			},
		},
		{
			name: "returns 404 status code when the car was not found",
			args: args{
				requestID: car.ID.String(),
			},
			wants: wants{
				statusCode: http.StatusNotFound,
			},
			setMocks: func(d *carsDependencies) {
				d.carsService.EXPECT().Get(gomock.Any(), car.ID).Return(domain.Car{}, errors.New(services.ErrCarNotFound))
			},
		},
		{
			name: "returns 500 status code when there is a server error",
			args: args{
				requestID: car.ID.String(),
			},
			wants: wants{
				statusCode: http.StatusInternalServerError,
			},
			setMocks: func(d *carsDependencies) {
				d.carsService.EXPECT().Get(gomock.Any(), car.ID).Return(domain.Car{}, errors.New("error getting car"))
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			mockCtlr := gomock.NewController(t)
			carsSrv := mocks.NewMockCarsService(mockCtlr)
			d := NewCarsDependencies(carsSrv)
			test.setMocks(d)

			baseURL := "/api/v1/"
			values := url.Values{}
			values.Set("id", test.args.requestID)
			urlObj, _ := url.Parse(baseURL + "cars/" + values.Encode())
			URL := urlObj.String()

			req, err := http.NewRequest(http.MethodGet, URL, nil)
			if err != nil {
				t.Fatal(err)
			}

			// Include request vars for gorilla mux to interpret path params
			vars := map[string]string{
				"id": test.args.requestID,
			}
			req = mux.SetURLVars(req, vars)

			rr := httptest.NewRecorder()

			carsHandler := NewCars(carsSrv)
			carsHandler.Get(rr, req)

			assert.Equal(t, test.wants.statusCode, rr.Code)
		})
	}
}

func TestCarsFullUpdate(t *testing.T) {
	car := dtos.Car{
		ID:             uuid.New(),
		Type:           "Sedan",
		Seats:          4,
		HourlyRentCost: 21.1,
		City:           "Los Angeles",
		Status:         "Available",
	}

	type args struct {
		requestID string
		car       dtos.Car
	}
	type wants struct {
		statusCode int
	}
	tests := []struct {
		name     string
		args     args
		wants    wants
		setMocks func(*carsDependencies)
	}{
		{
			name: "returns status code 200 when car was successfully updated",
			args: args{
				requestID: car.ID.String(),
				car:       car,
			},
			wants: wants{
				statusCode: http.StatusOK,
			},
			setMocks: func(d *carsDependencies) {
				d.carsService.EXPECT().FullUpdate(gomock.Any(), car.ToDomain()).Return(nil)
			},
		},
		{
			name: "returns 400 status code when path param id is not an uuid",
			args: args{
				requestID: "this-is-an-not-uuid",
			},
			wants: wants{
				statusCode: http.StatusBadRequest,
			},
			setMocks: func(d *carsDependencies) {
			},
		},
		{
			name: "returns 400 status code when car status in body is not allowed",
			args: args{
				requestID: car.ID.String(),
				car: dtos.Car{
					Status: "available",
				},
			},
			wants: wants{
				statusCode: http.StatusBadRequest,
			},
			setMocks: func(d *carsDependencies) {
			},
		},
		{
			name: "returns 404 status code when car was not found by id",
			args: args{
				requestID: car.ID.String(),
				car:       car,
			},
			wants: wants{
				statusCode: http.StatusNotFound,
			},
			setMocks: func(d *carsDependencies) {
				d.carsService.EXPECT().FullUpdate(gomock.Any(), car.ToDomain()).Return(errors.New(services.ErrCarNotFound))
			},
		},
		{
			name: "returns 500 status code when car service fails to update the car",
			args: args{
				requestID: car.ID.String(),
				car:       car,
			},
			wants: wants{
				statusCode: http.StatusInternalServerError,
			},
			setMocks: func(d *carsDependencies) {
				d.carsService.EXPECT().FullUpdate(gomock.Any(), car.ToDomain()).Return(errors.New("error registering car"))
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			mockCtlr := gomock.NewController(t)
			carsSrv := mocks.NewMockCarsService(mockCtlr)
			d := NewCarsDependencies(carsSrv)
			test.setMocks(d)

			baseURL := "/api/v1/"
			values := url.Values{}
			values.Set("id", test.args.requestID)
			urlObj, _ := url.Parse(baseURL + "cars/" + values.Encode())
			URL := urlObj.String()

			body, _ := json.Marshal(test.args.car)
			req, err := http.NewRequest(http.MethodPut, URL, bytes.NewBuffer(body))
			if err != nil {
				t.Fatal(err)
			}

			// Include request vars for gorilla mux to interpret path params
			vars := map[string]string{
				"id": test.args.requestID,
			}
			req = mux.SetURLVars(req, vars)

			rr := httptest.NewRecorder()

			carsHandler := NewCars(carsSrv)
			carsHandler.FullUpdate(rr, req)

			assert.Equal(t, test.wants.statusCode, rr.Code)
		})
	}
}
