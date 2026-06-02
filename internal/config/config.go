package config

import (
	"fmt"
	"github.com/joho/godotenv"
)

func Load() (*Config, error) {
	godotenv.Load(".env")
	return &Config{
		Server: ServerConfig{
			Port: ":" + GetEnv("PORT", "8000"),
			Passkey: GetEnv("PASSKEY", ""),
		},
		Database: DatabaseConfig{
			Name: GetEnv("DB_NAME", "sendin"),
			Username: GetEnv("DB_USER", "postgres"),
			Password: GetEnv("DB_PASSWORD", "password"),
			Host: GetEnv("DB_HOST", "localhost"),
			Port: GetEnv("DB_PORT", "5432"),
			SSL: GetEnv("DB_SSLMODE", "disable"),
		},
	}, nil
}

func (cfg *DatabaseConfig) GetDatabaseDSN() string {
	return fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
        cfg.Host,
        cfg.Port,
        cfg.Username,
        cfg.Password,
        cfg.Name,
        cfg.SSL,
	)
}

func (cfg *DatabaseConfig) GetAdminDSN() string {
	return fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=postgres sslmode=%s",
		cfg.Host,
		cfg.Port,
		cfg.Username,
		cfg.Password,
		cfg.SSL,
	)
}