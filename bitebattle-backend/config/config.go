package config

import "os"

var GooglePlacesAPIKey string

type DbConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Database string
}

var dbConfig DbConfig

func Load() {

	GooglePlacesAPIKey = os.Getenv("GOOGLE_PLACES_API_KEY")
	dbConfig = DbConfig{
		Host:     os.Getenv("DB_HOST"),
		Port:     os.Getenv("DB_PORT"),
		User:     os.Getenv("DB_USER"),
		Password: os.Getenv("DB_PASSWORD"),
		Database: os.Getenv("DB_NAME"),
	}
}
