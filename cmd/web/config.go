package web

import (
	"os"

	"github.com/spf13/viper"
)

func LoadConfig() error {
	fileName := os.Getenv("env")
	viper.SetConfigName(fileName + ".env")
	viper.AddConfigPath("./config")
	err := viper.ReadInConfig()
	if err != nil {
		return err
	}
	return nil
}
