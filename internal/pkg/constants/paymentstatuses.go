package constants

import "reflect"

var paymentStatuses []string

// Initializes paymentTypes slice with values of PAYMENT_STATUSES struct
func initializePaymentStatusesValues() {
	valueOf := reflect.ValueOf(Values.PAYMENT_STATUSES)
	for i := 0; i < valueOf.NumField(); i++ {
		field := valueOf.Field(i)
		paymentStatuses = append(paymentStatuses, field.String())
	}
}

type PAYMENT_STATUSES struct {
	PAID     string `mapstructure:"PAID"`
	PENDING  string `mapstructure:"PENDING"`
	CANCELED string `mapstructure:"CANCELED"`
}

// Get the values in payment statuses
func (cs *PAYMENT_STATUSES) Values() []string {
	return paymentStatuses
}
