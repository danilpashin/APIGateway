package config

import "pkg/env"

type Config struct {
	HTTPPort   int
	DBHost     string
	DBName     string
	DBPort     int
	DBUser     string
	DBPassword string
}

func Load() *Config {
	return &Config{
		HTTPPort:   env.GetEnvAsInt("HTTP_PORT"),
		DBHost:     env.GetEnv("DB_HOST"),
		DBName:     env.GetEnv("DB_NAME"),
		DBPort:     env.GetEnvAsInt("DB_PORT"),
		DBUser:     env.GetEnv("APP_DB_USER"),
		DBPassword: env.GetEnv("APP_DB_PASSWORD"),
	}
}
