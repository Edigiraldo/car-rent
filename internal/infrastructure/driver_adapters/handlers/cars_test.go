package handlers

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Edigiraldo/car-rent/internal/core/domain"
	"github.com/Edigiraldo/car-rent/internal/infrastructure/driver_adapters/handlers/dtos"
	mocks "github.com/Edigiraldo/car-rent/internal/pkg/mocks"
	"github.com/golang/mock/gomock"
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
			name: "returns status code 200 when body is appropriate",
			args: args{
				car: car,
			},
			wants: wants{
				statusCode: http.StatusOK,
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

			body, _ := json.Marshal(test.args.car)
			req, err := http.NewRequest("POST", "/cars", bytes.NewBuffer(body))
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
