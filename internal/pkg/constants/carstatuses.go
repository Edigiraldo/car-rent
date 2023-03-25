package constants

import "reflect"

var carStatuses []string

// Initializes carTypes slice with values of CAR_STATUSES struct
func initializeCarStatusesValues() {
	valueOf := reflect.ValueOf(Values.CAR_STATUSES)
	for i := 0; i < valueOf.NumField(); i++ {
		field := valueOf.Field(i)
		carStatuses = append(carStatuses, field.String())
	}

}

type CAR_STATUSES struct {
	AVAILABLE   string `mapstructure:"AVAILABLE"`
	UNAVAILABLE string `mapstructure:"UNAVAILABLE"`
}

// Get the values in car statuses
func (cs *CAR_STATUSES) Values() []string {
	return carStatuses
}
