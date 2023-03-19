package dtos

import (
	"bytes"
	"encoding/json"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCarFromBody(t *testing.T) {
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
