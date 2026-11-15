// Package config centraliza la configuración del backend. Todo lo sensible se
// lee de variables de entorno (ver .env.example); no hay secretos en el código.
package config

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// Config es la configuración inmutable que usa todo el backend.
type Config struct {
	DBServer      string
	DBPort        int
	DBUser        string
	DBPassword    string
	DBName        string
	JWTSecret     string
	RedisAddr     string
	RedisPassword string
	Port          string
	CORSOrigins   string
}

// Load lee la configuración del entorno. Los valores por defecto son los de un
// entorno de desarrollo local; lo sensible (DB_PASSWORD, JWT_SECRET) no tiene
// valor por defecto y la app falla rápido si falta.
func Load() *Config {
	loadDotEnv()
	return &Config{
		DBServer:      getenv("DB_SERVER", "10.10.1.21"),
		DBPort:        getenvInt("DB_PORT", 1433),
		DBUser:        getenv("DB_USER", "sa"),
		DBPassword:    os.Getenv("DB_PASSWORD"),
		DBName:        getenv("DB_NAME", "Hemingway"),
		JWTSecret:     os.Getenv("JWT_SECRET"),
		RedisAddr:     getenv("REDIS_ADDR", "localhost:6379"),
		RedisPassword: os.Getenv("REDIS_PASSWORD"),
		Port:          getenv("PORT", "8080"),
		CORSOrigins:   getenv("CORS_ORIGINS", "http://localhost:5173"),
	}
}

// DSN devuelve el connection string para go-mssqldb.
func (c *Config) DSN(database string) string {
	return "server=" + c.DBServer +
		";port=" + strconv.Itoa(c.DBPort) +
		";database=" + database +
		";user id=" + c.DBUser +
		";password=" + c.DBPassword +
		";encrypt=disable"
}

// Validate falla rápido si falta configuración sensible.
func (c *Config) Validate() error {
	var missing []string
	if c.DBPassword == "" {
		missing = append(missing, "DB_PASSWORD")
	}
	if c.JWTSecret == "" {
		missing = append(missing, "JWT_SECRET")
	}
	if len(missing) > 0 {
		return fmt.Errorf("faltan variables de entorno requeridas: %s (copia .env.example a .env)",
			strings.Join(missing, ", "))
	}
	return nil
}

// loadDotEnv carga el archivo .env del directorio actual si existe. Las variables
// que ya están definidas en el entorno ganan, para que Docker/systemd puedan
// sobrescribirlas. Sin dependencias externas.
func loadDotEnv() {
	f, err := os.Open(".env")
	if err != nil {
		return
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		line = strings.TrimPrefix(line, "export ")

		key, value, ok := strings.Cut(line, "=")
		if !ok {
			continue
		}
		key = strings.TrimSpace(key)
		value = strings.Trim(strings.TrimSpace(value), `"'`)

		if key != "" && os.Getenv(key) == "" {
			_ = os.Setenv(key, value)
		}
	}
}

func getenv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func getenvInt(key string, fallback int) int {
	if v := os.Getenv(key); v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			return n
		}
	}
	return fallback
}
