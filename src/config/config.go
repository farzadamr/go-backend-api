package config

import (
	"time"
)

type Config struct {
	Env      string
	HTTP     HTTPConfig
	Database DatabaseConfig
	JWT      JWTConfig
	Logger   LoggerConfig
}

type HTTPConfig struct {
	Port         string
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	IdleTimeout  time.Duration
}

type DatabaseConfig struct {
	Host            string
	Port            int
	User            string
	Password        string
	DBName          string
	SSLMode         string
	MaxOpenConns    int
	MaxIdleConns    int
	ConnMaxLifetime time.Duration
}

type JWTConfig struct {
	Secret    string
	ExpiresIn time.Duration
	RefreshIn time.Duration
}

type LoggerConfig struct {
	Level      string
	JSONFormat bool
	LogFile    string
	MaxSize    int
	MaxBackups int
	MaxAge     int
	Compress   bool
}
