package constants

import (
	"github.com/spf13/viper"
)

var Values ConstantValues

type ConstantValues struct {
	CARS_PER_PAGE uint16    `mapstructure:"CARS_PER_PAGE"`
	NULL_UUID     string    `mapstructure:"NULL_UUID"`
	CAR_TYPES     CAR_TYPES `mapstructure:"CAR_TYPES"`
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

	return nil
}
