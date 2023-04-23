package dtos

import (
	"bytes"
	"encoding/json"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestReservationFromBody(t *testing.T) {
	initConstantsFromDtos(t)

	type args struct {
		reservation Reservation
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
			name: "returns nil error when body structure is as expected",
			args: args{
				reservation: Reservation{
					UserID:        uuid.New(),
					CarID:         uuid.New(),
					Status:        "Reserved",
					PaymentStatus: "Pending",
					StartDate:     time.Now(),
					EndDate:       time.Now().AddDate(0, 0, 7),
				},
			},
			wants: wants{
				err: nil,
			},
		},
		{
			name: "returns invalid reservation status when received one is not one of the expected values",
			args: args{
				reservation: Reservation{
					UserID:        uuid.New(),
					CarID:         uuid.New(),
					Status:        "reserved",
					PaymentStatus: "Pending",
					StartDate:     time.Now(),
					EndDate:       time.Now().AddDate(0, 0, 7),
				},
			},
			wants: wants{
				err: errors.New(ErrInvalidReservationStatus),
			},
		},
		{
			name: "returns invalid payment status when received one is not one of the expected values",
			args: args{
				reservation: Reservation{
					UserID:        uuid.New(),
					CarID:         uuid.New(),
					Status:        "Reserved",
					PaymentStatus: "pending",
					StartDate:     time.Now(),
					EndDate:       time.Now().AddDate(0, 0, 7),
				},
			},
			wants: wants{
				err: errors.New(ErrInvalidPaymentStatus),
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			reservationJSON, err := json.Marshal(test.args.reservation)
			if err != nil {
				t.Fatal(err)
			}

			reservationReader := bytes.NewBuffer(reservationJSON)
			_, err = ReservationFromBody(reservationReader)

			assert.Equal(t, test.wants.err, err)
		})
	}
}
