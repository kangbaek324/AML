package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	AppPort string

	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string

	SourceDBName string
}

func Load() *Config {
	_ = godotenv.Load()

	return &Config{
		AppPort: getEnv("APP_PORT", "8080"),

		DBHost:     getEnv("DB_HOST", "127.0.0.1"),
		DBPort:     getEnv("DB_PORT", "3306"),
		DBUser:     getEnv("DB_USER", "root"),
		DBPassword: getEnv("DB_PASSWORD", ""),
		DBName:     getEnv("DB_NAME", "aml"),

		SourceDBName: getEnv("SOURCE_DB_NAME", "stock"),
	}
}

// DSN builds a MySQL data source name for the AML database.
func (c *Config) DSN() string {
	return c.dsnFor(c.DBName)
}

// SourceDSN builds a MySQL data source name for the Kronex source database.
// AML DB와 Source DB는 같은 인스턴스를 사용하므로 접속 계정은 공유한다.
func (c *Config) SourceDSN() string {
	return c.dsnFor(c.SourceDBName)
}

func (c *Config) dsnFor(dbName string) string {
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true",
		c.DBUser, c.DBPassword, c.DBHost, c.DBPort, dbName)
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
