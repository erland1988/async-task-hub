package global

import (
	"fmt"
	"github.com/joho/godotenv"
	"os"
	"strconv"
)

type Config struct {
	BASE_PATH               string
	APP_ENV                 string
	APP_LOG_MODE            string
	APP_LOG_FILENAME        string
	DATABASE_DRIVER         string
	DATABASE_HOST           string
	DATABASE_PORT           int
	DATABASE_DATABASE       string
	DATABASE_USERNAME       string
	DATABASE_PASSWORD       string
	DATABASE_CHARSET        string
	DATABASE_MAX_IDLE_CONNS int
	DATABASE_MAX_OPEN_CONNS int
	DATABASE_LOG_MODE       string
	REDIS_HOST              string
	REDIS_PORT              string
	REDIS_PASSWORD          string
	REDIS_DB                int
	REDIS_PREFIX            string
}

var CONFIG Config

func InitializeConfig() {
	// 尝试加载 .env 文件，但不报错
	if err := godotenv.Load(".env"); err != nil {
		fmt.Println("No .env file found, using environment variables")
	}

	CONFIG = Config{
		BASE_PATH:               getEnv("BASE_PATH", ""),
		APP_ENV:                 getEnv("APP_ENV", "development"),
		APP_LOG_MODE:            getEnv("APP_LOG_MODE", "warn"),
		APP_LOG_FILENAME:        getEnv("APP_LOG_FILENAME", ""),
		DATABASE_DRIVER:         getEnv("DATABASE_DRIVER", "mysql"),
		DATABASE_HOST:           getEnv("DATABASE_HOST", "localhost"),
		DATABASE_PORT:           getEnvAsInt("DATABASE_PORT", 3306),
		DATABASE_DATABASE:       getEnv("DATABASE_DATABASE", "testdb"),
		DATABASE_USERNAME:       getEnv("DATABASE_USERNAME", "root"),
		DATABASE_PASSWORD:       getEnv("DATABASE_PASSWORD", ""),
		DATABASE_CHARSET:        getEnv("DATABASE_CHARSET", "utf8mb4"),
		DATABASE_MAX_IDLE_CONNS: getEnvAsInt("DATABASE_MAX_IDLE_CONNS", 10),
		DATABASE_MAX_OPEN_CONNS: getEnvAsInt("DATABASE_MAX_OPEN_CONNS", 100),
		DATABASE_LOG_MODE:       getEnv("DATABASE_LOG_MODE", "false"),
		REDIS_HOST:              getEnv("REDIS_HOST", "localhost"),
		REDIS_PORT:              getEnv("REDIS_PORT", "6379"),
		REDIS_PASSWORD:          getEnv("REDIS_PASSWORD", ""),
		REDIS_DB:                getEnvAsInt("REDIS_DB", 0),
		REDIS_PREFIX:            getEnv("REDIS_PREFIX", ""),
	}
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

func getEnvAsInt(name string, defaultValue int) int {
	if valueStr, exists := os.LookupEnv(name); exists {
		if value, err := strconv.Atoi(valueStr); err == nil {
			return value
		}
	}
	return defaultValue
}

func getEnvAsBool(name string, defaultValue bool) bool {
	if valueStr, exists := os.LookupEnv(name); exists {
		if value, err := strconv.ParseBool(valueStr); err == nil {
			return value
		}
	}
	return defaultValue
}
