package constants

import "reflect"

var userStatuses []string

// Initializes userTypes slice with values of USER_STATUSES struct
func initializeUserStatusesValues() {
	valueOf := reflect.ValueOf(Values.USER_STATUSES)
	for i := 0; i < valueOf.NumField(); i++ {
		field := valueOf.Field(i)
		userStatuses = append(userStatuses, field.String())
	}

}

type USER_STATUSES struct {
	ACTIVE   string `mapstructure:"ACTIVE"`
	INACTIVE string `mapstructure:"INACTIVE"`
}

// Get the values in user statuses
func (cs *USER_STATUSES) Values() []string {
	return userStatuses
}
