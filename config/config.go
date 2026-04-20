package config

import (
	"log"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/spf13/viper"
)

var config *viper.Viper

// Init initializes the configuration using viper
func Init() {
	config = viper.New()
	config.SetConfigType("yaml")
	config.SetConfigName("config")
	config.AddConfigPath("/etc/secrets")
	config.AddConfigPath(getConfigDir())
	config.AddConfigPath(".") // Look in current directory as well
	config.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	config.AutomaticEnv()
	config.SetDefault("server.port", ":5000")

	bindEnv("server.port", "SERVER_PORT", "PORT")
	bindEnv("pdf_service.base_url", "PDF_SERVICE_BASE_URL")
	bindEnv("mongo.uri", "MONGO_URI")
	bindEnv("mongo.database", "MONGO_DATABASE")
	bindEnv("extraction_service.base_url", "EXTRACTION_SERVICE_BASE_URL")
	bindEnv("jwt.secret", "JWT_SECRET")

	if err := config.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			log.Fatalf("Error reading config file: %v", err)
		}
		log.Printf("Config file not found, falling back to environment variables and defaults")
	}
}

func bindEnv(key string, envNames ...string) {
	args := append([]string{key}, envNames...)
	if err := config.BindEnv(args...); err != nil {
		log.Fatalf("Error binding env for %s: %v", key, err)
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
