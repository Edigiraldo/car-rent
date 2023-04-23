package constants

import "reflect"

var reservationStatuses []string

// Initializes reservationTypes slice with values of RESERVATION_STATUSES struct
func initializeReservationStatusesValues() {
	valueOf := reflect.ValueOf(Values.RESERVATION_STATUSES)
	for i := 0; i < valueOf.NumField(); i++ {
		field := valueOf.Field(i)
		reservationStatuses = append(reservationStatuses, field.String())
	}
}

type RESERVATION_STATUSES struct {
	RESERVED  string `mapstructure:"RESERVED"`
	CANCELED  string `mapstructure:"CANCELED"`
	COMPLETED string `mapstructure:"COMPLETED"`
}

// Get the values in reservation statuses
func (cs *RESERVATION_STATUSES) Values() []string {
	return reservationStatuses
}
