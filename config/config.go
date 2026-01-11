package config

import (
	"log"
	"path/filepath"
	"runtime"

	"github.com/spf13/viper"
)

var config *viper.Viper

// Init initializes the configuration using viper
func Init() {
	config = viper.New()
	config.SetConfigType("yaml")
	config.SetConfigName("config")
	config.AddConfigPath(getConfigDir())
	config.AddConfigPath(".") // Look in current directory as well

	if err := config.ReadInConfig(); err != nil {
		log.Fatalf("Error reading config file: %v", err)
	}
}

func getConfigDir() string {
	_, b, _, _ := runtime.Caller(0)
	return filepath.Dir(b)
}

func GetConfig() *viper.Viper {
	if config == nil {
		Init()
	}
	return config
}

func GetString(key string) string {
	return GetConfig().GetString(key)
}
