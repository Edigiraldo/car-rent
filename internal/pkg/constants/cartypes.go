package constants

import "reflect"

var carTypes []string

// Initializes carTypes slice with values of CAR_TYPES struct
func initializeCarTypesValues() {
	valueOf := reflect.ValueOf(Values.CAR_TYPES)
	for i := 0; i < valueOf.NumField(); i++ {
		field := valueOf.Field(i)
		carTypes = append(carTypes, field.String())
	}

}

type CAR_TYPES struct {
	SEDAN      string `mapstructure:"SEDAN"`
	LUXURY     string `mapstructure:"LUXURY"`
	SPORTS_CAR string `mapstructure:"SPORTS CAR"`
	LIMOUSINE  string `mapstructure:"LIMOUSINE"`
}

// Get the values in car types
func (ct *CAR_TYPES) Values() []string {
	return carTypes
}
