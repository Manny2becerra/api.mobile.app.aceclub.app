package config

import (
	"os"

	_ "github.com/joho/godotenv/autoload"
)

type Config struct {
	Port                            string
	DB                              DatabaseConfig
	JWTAccessSecret                 string
	JWTRefreshSecret                string
	EmailHost                       string
	EmailHostUser                   string
	EmailHostPassword               string
	EmailHostPort                   string
	EmailHostFrom                   string
	EstablishmentProfileImageBucket string
	AWSRegion                       string
	AWSAccessKey                    string
	AWSSecretKey                    string
	AWSSessionToken                 string
}

type DatabaseConfig struct {
	ConnectionString string
	DbName           string
}

func LoadConfig() (*Config, error) {
	ServerConfigurations := Config{
		Port:                            os.Getenv("PORT"),
		JWTAccessSecret:                 os.Getenv("JWT_ACCESS_SECRET"),
		JWTRefreshSecret:                os.Getenv("JWT_REFRESH_SECRET"),
		EmailHost:                       os.Getenv("EMAIL_HOST"),
		EmailHostUser:                   os.Getenv("EMAIL_HOST_USER"),
		EmailHostPassword:               os.Getenv("EMAIL_HOST_PASSWORD"),
		EmailHostPort:                   os.Getenv("EMAIL_HOST_PORT"),
		EmailHostFrom:                   os.Getenv("EMAIL_HOST_FROM"),
		EstablishmentProfileImageBucket: os.Getenv("ESTABLISHMENT_PROFILE_IMAGE_BUCKET"),
		AWSRegion:                       os.Getenv("AWS_REGION"),
		AWSAccessKey:                    os.Getenv("AWS_ACCESS_KEY"),
		AWSSecretKey:                    os.Getenv("AWS_SECRET_KEY"),
		DB: DatabaseConfig{
			ConnectionString: os.Getenv("DATABASE_URL"),
			DbName:           os.Getenv("DATABASE_NAME"),
		},
	}
	return &ServerConfigurations, nil
}
