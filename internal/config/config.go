package config

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

type Config struct {
	Server struct {
		Addr            string
		ReadTimeout     time.Duration
		WriteTimeout    time.Duration
		ShutdownTimeout time.Duration
		EnableSwagger   bool
	}
	DB struct {
		DSN             string
		Host            string
		Port            int
		User            string
		Password        string
		Name            string
		SSLMode         string
		MaxConns        int32
		MinConns        int32
		ConnMaxLifetime time.Duration
		ConnMaxIdleTime time.Duration
	}
	Log struct {
		Level  string
		Format string
	}
}

func Load() (*Config, error) {
	var c Config

	// Считываем из env, если пусто, то ставим дефолтные значения

	//Server
	c.Server.Addr = getEnv("SERVER_ADDR", ":8080")
	c.Server.ReadTimeout = getEnvDur("SERVER_READ_TIMEOUT", 5*time.Second)
	c.Server.WriteTimeout = getEnvDur("SERVER_WRITE_TIMEOUT", 10*time.Second)
	c.Server.ShutdownTimeout = getEnvDur("SERVER_SHUTDOWN_TIMEOUT", 5*time.Second)

	//DB
	c.DB.DSN = os.Getenv("DB_DSN")
	c.DB.Host = getEnv("DB_HOST", "localhost")
	c.DB.Port = getEnvInt("DB_PORT", 5435)
	c.DB.User = getEnv("DB_USER", "subs")
	c.DB.Password = getEnv("DB_PASSWORD", "subs")
	c.DB.Name = getEnv("DB_NAME", "subs")
	c.DB.SSLMode = getEnv("DB_SSLMODE", "disable")
	c.DB.MaxConns = int32(getEnvInt("DB_MAX_CONNS", 10))
	c.DB.MinConns = int32(getEnvInt("DB_MIN_CONNS", 0))
	c.DB.ConnMaxLifetime = getEnvDur("DB_CONN_MAX_LIFETIME", 30*time.Minute)
	c.DB.ConnMaxIdleTime = getEnvDur("DB_CONN_MAX_IDLE_TIME", 10*time.Minute)

	// Собираем DSN, если не дан целиком
	if strings.TrimSpace(c.DB.DSN) == "" {
		c.DB.DSN = fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=%s",
			c.DB.User, c.DB.Password, c.DB.Host, c.DB.Port, c.DB.Name, c.DB.SSLMode)
	}

	//Log
	c.Log.Level = getEnv("LOG_LEVEL", "info")
	c.Log.Format = getEnv("LOG_FORMAT", "json")

	if c.DB.DSN == "" {
		return nil, errors.New("empty DB_DSN")
	}
	return &c, nil
}

// Функции проверяют переменные из env, и если пусто, то возвращают дефолтные значения разных типов
func getEnv(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}

func getEnvInt(key string, def int) int {
	if v := os.Getenv(key); v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			return n
		}
	}
	return def
}

func getEnvDur(key string, def time.Duration) time.Duration {
	if v := os.Getenv(key); v != "" {
		if d, err := time.ParseDuration(v); err == nil {
			return d
		}
	}
	return def
}
