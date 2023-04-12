package constants

import "reflect"

var userTypes []string

// Initializes userTypes slice with values of USER_TYPES struct
func initializeUserTypesValues() {
	valueOf := reflect.ValueOf(Values.USER_TYPES)
	for i := 0; i < valueOf.NumField(); i++ {
		field := valueOf.Field(i)
		userTypes = append(userTypes, field.String())
	}

}

type USER_TYPES struct {
	CUSTOMER string `mapstructure:"CUSTOMER"`
	ADMIN    string `mapstructure:"ADMIN"`
}

// Get the values in car types
func (ct *USER_TYPES) Values() []string {
	return userTypes
}
