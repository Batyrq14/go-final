package config

import (
	"os"
)

type Config struct {
	Port        string
	DBUrl       string
	RedisUrl    string
	RabbitMQUrl string
	JWTSecret   string
	Services    ServiceConfig
}

type ServiceConfig struct {
	UserUrl        string
	MarketplaceUrl string
	ChatUrl        string
}

func Load() *Config {
	return &Config{
		Port:        getEnv("PORT", "8080"),
		DBUrl:       getEnv("DATABASE_URL", "postgres://user:password@localhost:5433/qasynda?sslmode=disable"),
		RedisUrl:    getEnv("REDIS_URL", "localhost:6379"),
		RabbitMQUrl: getEnv("RABBITMQ_URL", "amqp://user:password@localhost:5672/"),
		JWTSecret:   getEnv("JWT_SECRET", "very-secret-key"),
		Services: ServiceConfig{
			UserUrl:        getEnv("USER_SERVICE_URL", "localhost:50051"),
			MarketplaceUrl: getEnv("MARKETPLACE_SERVICE_URL", "localhost:50052"),
			ChatUrl:        getEnv("CHAT_SERVICE_URL", "localhost:50053"),
		},
	}
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

func GetUserPort() string {
	return getEnv("USER_PORT", ":50051")
}

func GetMarketplacePort() string {
	return getEnv("MARKETPLACE_PORT", ":50052")
}

func GetChatPort() string {
	return getEnv("CHAT_PORT", ":50053")
}
