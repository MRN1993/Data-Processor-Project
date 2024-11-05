package config

import (
	"os"

)

type Config struct {
    RedisHost      string
    RedisPort      string
    KafkaHost      string
}

func LoadConfig() Config {

    config := Config{
        RedisHost:      getEnvWithLogging("REDIS_HOST", "127.0.0.1"),
        RedisPort:      getEnvWithLogging("REDIS_PORT", "6379"),
        KafkaHost:      getEnvWithLogging("KAFKA_HOST", "localhost:9092"),
    }

   

    return config
}

func getEnvWithLogging(key, fallback string) string {
    if value, exists := os.LookupEnv(key); exists {
        return value
    }

    return fallback
}