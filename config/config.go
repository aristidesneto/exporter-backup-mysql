package config

import (
	"log"

	"github.com/spf13/viper"
)

func Configuration(path string) {
    viper.SetConfigName("config")
    viper.SetConfigType("yaml")
    viper.AddConfigPath(path)
	
	err := viper.ReadInConfig()
    if err != nil {
        log.Fatalf("Error to open config file: %s", err)
    }
}