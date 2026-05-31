package config

import (
	"fmt"
	"log"
	"time"
)

func Parse() (*Config, error) {
	cfg := &Config{
		Env: getEnv("APP_ENV", "development"),
		HTTP: HTTPConfig{
			Port:         getEnv("HTTP_PORT", "8080"),
			ReadTimeout:  getEnvAsDuration("HTTP_READ_TIMEOUT", 5*time.Second),
			WriteTimeout: getEnvAsDuration("HTTP_WRITE_TIMEOUT", 10*time.Second),
			IdleTimeout:  getEnvAsDuration("HTTP_IDLE_TIMEOUT", time.Minute),
		},
		Database: DatabaseConfig{
			Host:            getEnv("DB_HOST", "localhost"),
			Port:            getEnvAsInt("DB_PORT", 5432),
			User:            getEnv("DB_USER", "postgres"),
			Password:        getEnv("DB_PASSWORD", ""),
			DBName:          getEnv("DB_NAME", "myapp"),
			SSLMode:         getEnv("DB_SSLMODE", "disable"),
			MaxOpenConns:    getEnvAsInt("DB_MAX_OPEN_CONNS", 25),
			MaxIdleConns:    getEnvAsInt("DB_MAX_IDLE_CONNS", 25),
			ConnMaxLifetime: getEnvAsDuration("DB_CONN_MAX_LIFETIME", 5*time.Minute),
		},
		JWT: JWTConfig{
			Secret:    getEnv("JWT_SECRET", ""),
			ExpiresIn: getEnvAsDuration("JWT_EXPIRES_IN", 24*time.Hour),
			RefreshIn: getEnvAsDuration("JWT_REFRESH_IN", 72*time.Hour),
		},
		Logger: LoggerConfig{
			Level:      getEnv("LOG_LEVEL", "info"),
			JSONFormat: getEnvAsBool("LOG_JSON_FORMAT", true),
			LogFile:    getEnv("LOG_FILE", "./logs/app.log"),
			MaxSize:    getEnvAsInt("LOG_MAX_SIZE", 100),
			MaxBackups: getEnvAsInt("LOG_MAX_BACKUPS", 5),
			MaxAge:     getEnvAsInt("LOG_MAX_AGE", 25),
			Compress:   getEnvAsBool("LOG_COMPRESS", true),
		},
	}

	// validation
	if cfg.JWT.Secret == "" && cfg.Env != "test" {
		return nil, fmt.Errorf("JWT_SECRET is required in non-test environment")
	}
	if cfg.Database.Password == "" && cfg.Env != "test" {
		log.Println("config warning: DB_PASSWORD is empty, this might cause issues")
	}
	return cfg, nil
}
