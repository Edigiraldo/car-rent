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

func TestCarsRegister(t *testing.T) {
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

func TestCarsGet(t *testing.T) {
	car := domain.Car{
		ID:             uuid.New(),
		Type:           "Sedan",
		Seats:          4,
		HourlyRentCost: 21.1,
		City:           "Los Angeles",
		Status:         "Available",
	}

	type args struct {
		ctx context.Context
		ID  uuid.UUID
	}
	type wants struct {
		car       domain.Car
		withError bool
	}
	tests := []struct {
		name     string
		args     args
		wants    wants
		setMocks func(*carsDependencies)
	}{
		{
			name: "returns nil error when car was found by the given id",
			args: args{
				ctx: context.TODO(),
				ID:  car.ID,
			},
			wants: wants{
				car:       car,
				withError: false,
			},
			setMocks: func(d *carsDependencies) {
				d.carsRepository.EXPECT().Get(gomock.Any(), car.ID).Return(car, nil)
			},
		},
		{
			name: "returns an error when car was not found by the given id",
			args: args{
				ctx: context.TODO(),
				ID:  car.ID,
			},
			wants: wants{
				car:       domain.Car{},
				withError: true,
			},
			setMocks: func(d *carsDependencies) {
				d.carsRepository.EXPECT().Get(gomock.Any(), car.ID).Return(domain.Car{}, errors.New(ErrCarNotFound))
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
			car, err := carsService.Get(test.args.ctx, test.args.ID)

			assert.Equal(t, test.wants.car, car)
			assert.Equal(t, test.wants.withError, err != nil)
		})
	}
}

func TestCarsFullUpdate(t *testing.T) {
	car := domain.Car{
		ID:             uuid.New(),
		Type:           "Luxury",
		Seats:          6,
		HourlyRentCost: 56.5,
		City:           "Austin",
		Status:         "Available",
	}

	type args struct {
		ctx context.Context
		car domain.Car
	}
	type wants struct {
		err error
	}
	tests := []struct {
		name     string
		args     args
		wants    wants
		setMocks func(*carsDependencies)
	}{
		{
			name: "returns nil error when car update was successfully updated",
			args: args{
				ctx: context.TODO(),
				car: car,
			},
			wants: wants{
				err: nil,
			},
			setMocks: func(d *carsDependencies) {
				d.carsRepository.EXPECT().FullUpdate(gomock.Any(), car).Return(nil)
			},
		},
		{
			name: "returns an error when car update fails",
			args: args{
				ctx: context.TODO(),
				car: car,
			},
			wants: wants{
				err: errors.New("failure while updating car"),
			},
			setMocks: func(d *carsDependencies) {
				d.carsRepository.EXPECT().FullUpdate(gomock.Any(), car).Return(errors.New("failure while updating car"))
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
			err := carsService.FullUpdate(test.args.ctx, test.args.car)

			assert.Equal(t, test.wants.err, err)
		})
	}
}

func TestCarsDelete(t *testing.T) {
	type args struct {
		ctx context.Context
		ID  uuid.UUID
	}
	type wants struct {
		err error
	}
	tests := []struct {
		name     string
		args     args
		wants    wants
		setMocks func(*carsDependencies)
	}{
		{
			name: "returns nil error when car deletion was successful",
			args: args{
				ctx: context.TODO(),
				ID:  uuid.New(),
			},
			wants: wants{
				err: nil,
			},
			setMocks: func(d *carsDependencies) {
				d.carsRepository.EXPECT().Delete(gomock.Any(), gomock.Any()).Return(nil)
			},
		},
		{
			name: "returns an error when car deletion fails",
			args: args{
				ctx: context.TODO(),
				ID:  uuid.New(),
			},
			wants: wants{
				err: errors.New("failure while updating car"),
			},
			setMocks: func(d *carsDependencies) {
				d.carsRepository.EXPECT().Delete(gomock.Any(), gomock.Any()).Return(errors.New("failure while updating car"))
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
			err := carsService.Delete(test.args.ctx, test.args.ID)

			assert.Equal(t, test.wants.err, err)
		})
	}
}
