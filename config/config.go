package config

import (
	"log"

	"github.com/spf13/viper"
)

type Database struct {
	Host     string `mapstructure:"host"`
	Port     string `mapstructure:"port"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	Name     string `mapstructure:"name"`
}

type Server struct {
	Port string `mapstructure:"port"`
}

type Config struct {
	Database Database `mapstructure:"database"`
	Server   Server   `mapstructure:"server"`
}

var AppConfig Config

func Init() {
	viper.SetConfigName("config.local")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./config/")

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error reading config file: %v", err)
	}

	viper.AutomaticEnv()

	if err := viper.Unmarshal(&AppConfig); err != nil {
		log.Fatalf("Error unmarshalling config: %v", err)
	}

	log.Println("Config initialized successfully")
}

func GetConfig() Config {
	return AppConfig
}
