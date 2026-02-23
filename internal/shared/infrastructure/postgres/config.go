// internal/shared/infrastructure/postgres/config.go
package postgres

import (
	"fmt"
	"os"
)

// DBConfig holds the PostgreSQL connection configuration.
type DBConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
}

// LoadDBConfigFromEnv reads DB connection settings from environment variables.
// Required: DB_HOST, DB_PORT, DB_USER, DB_PASSWORD, DB_NAME.
func LoadDBConfigFromEnv() (DBConfig, error) {
	cfg := DBConfig{
		Host:     os.Getenv("DB_HOST"),
		Port:     os.Getenv("DB_PORT"),
		User:     os.Getenv("DB_USER"),
		Password: os.Getenv("DB_PASSWORD"),
		DBName:   os.Getenv("DB_NAME"),
	}

	if cfg.Host == "" {
		return DBConfig{}, fmt.Errorf("DB_HOST no configurado")
	}
	if cfg.Port == "" {
		cfg.Port = "5432"
	}
	if cfg.User == "" {
		return DBConfig{}, fmt.Errorf("DB_USER no configurado")
	}
	if cfg.DBName == "" {
		return DBConfig{}, fmt.Errorf("DB_NAME no configurado")
	}

	return cfg, nil
}

// DSN returns the PostgreSQL connection string (Data Source Name).
func (c DBConfig) DSN() string {
	return fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		c.Host, c.Port, c.User, c.Password, c.DBName,
	)
}
