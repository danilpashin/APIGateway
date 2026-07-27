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
		HTTPPort:   env.GetEnvAsInt("ORDER_PORT", 8081),
		DBHost:     env.GetEnv("ORDER_DB_HOST"),
		DBName:     env.GetEnv("ORDER_DB_NAME"),
		DBPort:     env.GetEnvAsInt("ORDER_DB_PORT", 5432),
		DBUser:     env.GetEnv("ORDER_DB_USER"),
		DBPassword: env.GetEnv("ORDER_DB_PASSWORD"),
	}
}
