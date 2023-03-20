package services

import (
	"context"
	"errors"
	"testing"

	"github.com/Edigiraldo/car-rent/internal/core/domain"
	mocks "github.com/Edigiraldo/car-rent/internal/pkg/mocks"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

type carsDependencies struct {
	carsRepository *mocks.MockCarsRepo
}

func NewCarsDependencies(carsRepo *mocks.MockCarsRepo) *carsDependencies {
	return &carsDependencies{
		carsRepository: carsRepo,
	}
}

func TestRegister(t *testing.T) {
	type args struct {
		ctx context.Context
		car domain.Car
	}
	type wants struct {
		withError bool
	}
	tests := []struct {
		name     string
		args     args
		wants    wants
		setMocks func(*carsDependencies)
	}{
		{
			name: "returns nil error when car struct is populated appropriately",
			args: args{
				ctx: context.TODO(),
				car: domain.Car{
					ID:             uuid.UUID([]byte("2e7f6919-9f1e-4286-bbe1-7e75412e2912")),
					Type:           "Luxury",
					Seats:          6,
					HourlyRentCost: 56.5,
					City:           "Austin",
					Status:         "Available",
				},
			},
			wants: wants{
				withError: false,
			},
			setMocks: func(d *carsDependencies) {
				d.carsRepository.EXPECT().Insert(gomock.Any(), gomock.Any()).Return(nil)
			},
		},
		{
			name: "returns an error when car type is not an allowed value",
			args: args{
				ctx: context.TODO(),
				car: domain.Car{
					ID:             uuid.UUID([]byte("2e7f6919-9f1e-4286-bbe1-7e75412e2912")),
					Type:           "Luxu",
					Seats:          6,
					HourlyRentCost: 56.5,
					City:           "Austin",
					Status:         "Available",
				},
			},
			wants: wants{
				withError: true,
			},
			setMocks: func(d *carsDependencies) {
				d.carsRepository.EXPECT().Insert(gomock.Any(), gomock.Any()).Return(errors.New("type Luxu is not allowed"))
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			mockCtlr := gomock.NewController(t)
			carsRepo := mocks.NewMockCarsRepo(mockCtlr)
			d := NewCarsDependencies(carsRepo)
			test.setMocks(d)

			carsService := NewCars(carsRepo)
			_, err := carsService.Register(test.args.ctx, test.args.car)

			assert.Equal(t, test.wants.withError, err != nil)
		})
	}
}
