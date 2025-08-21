package config

import (
	"os"
	"strconv"
)

type Config struct {
	JWTSecret     string
	RedisAddr     string
	RedisPassword string
	RedisDB       int
}

func Load() *Config {
	redisDB, _ := strconv.Atoi(getEnv("REDIS_DB", "0"))

	return &Config{
		JWTSecret:     getEnv("JWT_SECRET", "your-secret-key-change-in-production"),
		RedisAddr:     getEnv("REDIS_ADDR", "localhost:6379"),
		RedisPassword: getEnv("REDIS_PASSWORD", ""),
		RedisDB:       redisDB,
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
