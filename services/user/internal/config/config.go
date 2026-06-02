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
		HTTPPort:   env.GetEnvAsInt("USER_PORT", 8083),
		DBHost:     env.GetEnv("USER_DB_HOST"),
		DBName:     env.GetEnv("USER_DB_NAME"),
		DBPort:     env.GetEnvAsInt("USER_DB_PORT", 5434),
		DBUser:     env.GetEnv("USER_DB_USER"),
		DBPassword: env.GetEnv("USER_DB_PASSWORD"),
	}
}
