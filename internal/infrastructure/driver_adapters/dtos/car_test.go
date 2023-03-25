package dtos

import (
	"bytes"
	"encoding/json"
	"errors"
	"testing"

	"github.com/Edigiraldo/car-rent/internal/pkg/constants"
	"github.com/stretchr/testify/assert"
)

var pathToRoot = "./../../../.."

func initConstantsFromDtos(t *testing.T) {
	if err := constants.InitValuesFrom(pathToRoot); err != nil {
		t.Fatal(err)
	}
}

func TestCarFromBody(t *testing.T) {
	initConstantsFromDtos(t)

	type args struct {
		car Car
	}
	type wants struct {
		err error
	}
	tests := []struct {
		name  string
		args  args
		wants wants
	}{
		{
			name: "returns invalid car type when received one is not one of the expected values",
			args: args{
				car: Car{
					Type:           "Seda",
					Seats:          4,
					HourlyRentCost: 21.1,
					City:           "Los Angeles",
					Status:         "Available",
				},
			},
			wants: wants{
				err: errors.New(ErrInvalidCarType),
			},
		},
		{
			name: "returns car structure when body structure is as expected",
			args: args{
				car: Car{
					Type:           "Sedan",
					Seats:          4,
					HourlyRentCost: 21.1,
					City:           "Los Angeles",
					Status:         "Available",
				},
			},
			wants: wants{
				err: nil,
			},
		},
		{
			name: "returns invalid seats number when received one is not greater than zero",
			args: args{
				car: Car{
					Type:           "Sedan",
					Seats:          -4,
					HourlyRentCost: 21.1,
					City:           "Los Angeles",
					Status:         "Available",
				},
			},
			wants: wants{
				err: errors.New(ErrInvalidSeatsNumber),
			},
		},
		{
			name: "returns invalid hourly rent cost when received one is not greater than zero",
			args: args{
				car: Car{
					Type:           "Sedan",
					Seats:          4,
					HourlyRentCost: -21.1,
					City:           "Los Angeles",
					Status:         "Available",
				},
			},
			wants: wants{
				err: errors.New(ErrInvalidHourlyRentCost),
			},
		},
		{
			name: "returns invalid empty city error when it is not provided",
			args: args{
				car: Car{
					Type:           "Sedan",
					Seats:          4,
					HourlyRentCost: 21.1,
					City:           "",
					Status:         "Available",
				},
			},
			wants: wants{
				err: errors.New(ErrEmptyCity),
			},
		},
		{
			name: "returns invalid car status when received one is not one of the expected values",
			args: args{
				car: Car{
					Type:           "Sedan",
					Seats:          4,
					HourlyRentCost: 21.1,
					City:           "Los Angeles",
					Status:         "available",
				},
			},
			wants: wants{
				err: errors.New(ErrInvalidCarStatus),
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			carJSON, err := json.Marshal(test.args.car)
			if err != nil {
				t.Fatal(err)
			}

			carReader := bytes.NewBuffer(carJSON)
			_, err = CarFromBody(carReader)

			assert.Equal(t, test.wants.err, err)
		})
	}
}
