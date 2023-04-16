package services

import (
	"context"
	"errors"
	"testing"

	mocks "github.com/Edigiraldo/car-rent/internal/pkg/mocks"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

type citiesDependencies struct {
	citiesRepository *mocks.MockCitiesRepo
}

func NewCitiesDependencies(citiesRepo *mocks.MockCitiesRepo) *citiesDependencies {
	return &citiesDependencies{
		citiesRepository: citiesRepo,
	}
}

func TestCitiesListNames(t *testing.T) {
	type args struct {
		ctx context.Context
	}
	type wants struct {
		cityNames []string
		err       error
	}
	tests := []struct {
		name     string
		args     args
		wants    wants
		setMocks func(*citiesDependencies)
	}{
		{
			name: "returns nil error when the cities name list was successful retrieved",
			args: args{
				ctx: context.TODO(),
			},
			wants: wants{
				cityNames: []string{"CityName A", "CityName B", "CityName C"},
				err:       nil,
			},
			setMocks: func(d *citiesDependencies) {
				d.citiesRepository.EXPECT().ListNames(gomock.Any()).Return([]string{"CityName A", "CityName B", "CityName C"}, nil)
			},
		},
		{
			name: "returns an error when server fails retrieving cities name",
			args: args{
				ctx: context.TODO(),
			},
			wants: wants{
				err: errors.New("failure listing cities name"),
			},
			setMocks: func(d *citiesDependencies) {
				d.citiesRepository.EXPECT().ListNames(gomock.Any()).Return(nil, errors.New("failure listing cities name"))
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			mockCtlr := gomock.NewController(t)
			citiesRepo := mocks.NewMockCitiesRepo(mockCtlr)
			d := NewCitiesDependencies(citiesRepo)
			test.setMocks(d)

			citiesService := NewCities(citiesRepo)
			cityNames, err := citiesService.ListNames(test.args.ctx)

			assert.Equal(t, test.wants.cityNames, cityNames)
			assert.Equal(t, test.wants.err, err)
		})
	}
}
