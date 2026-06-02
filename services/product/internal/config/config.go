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
		HTTPPort:   env.GetEnvAsInt("PRODUCT_PORT", 8082),
		DBHost:     env.GetEnv("PRODUCT_DB_HOST"),
		DBName:     env.GetEnv("PRODUCT_DB_NAME"),
		DBPort:     env.GetEnvAsInt("PRODUCT_DB_PORT", 5433),
		DBUser:     env.GetEnv("PRODUCT_DB_USER"),
		DBPassword: env.GetEnv("PRODUCT_DB_PASSWORD"),
	}
}
