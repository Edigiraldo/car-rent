package constants

import (
	"os"

	"github.com/spf13/viper"
)

var Values ConstantValues

type ConstantValues struct {
	CARS_PER_PAGE             uint16               `mapstructure:"CARS_PER_PAGE"`
	RESERVATIONS_PER_PAGE     uint16               `mapstructure:"RESERVATIONS_PER_PAGE"`
	MINIMUM_RESERVATION_HOURS uint16               `mapstructure:"MINIMUM_RESERVATION_HOURS"`
	NULL_UUID                 string               `mapstructure:"NULL_UUID"`
	DATETIME_LAYOUT           string               `mapstructure:"DATETIME_LAYOUT"`
	CAR_TYPES                 CAR_TYPES            `mapstructure:"CAR_TYPES"`
	CAR_STATUSES              CAR_STATUSES         `mapstructure:"CAR_STATUSES"`
	USER_TYPES                USER_TYPES           `mapstructure:"USER_TYPES"`
	USER_STATUSES             USER_STATUSES        `mapstructure:"USER_STATUSES"`
	RESERVATION_STATUSES      RESERVATION_STATUSES `mapstructure:"RESERVATION_STATUSES"`
	PAYMENT_STATUSES          PAYMENT_STATUSES     `mapstructure:"PAYMENT_STATUSES"`
}

func InitValues() error {
	viper.SetConfigName("constants")
	viper.SetConfigType("json")
	viper.AddConfigPath("./internal/pkg/constants")

	if err := viper.ReadInConfig(); err != nil {
		return err
	}

	if err := viper.Unmarshal(&Values); err != nil {
		return err
	}

	// initializes lists of values of structs to have iterable objects
	initializeCarTypesValues()
	initializeCarStatusesValues()
	initializeUserTypesValues()
	initializeUserStatusesValues()
	initializeReservationStatusesValues()
	initializePaymentStatusesValues()

	return nil
}

// Loads constants from a directory that is not in the root of the project
func InitValuesFrom(PathToRoot string) error {
	if err := os.Chdir(PathToRoot); err != nil {
		return err
	}

	if err := InitValues(); err != nil {
		return err
	}

	return nil
}
