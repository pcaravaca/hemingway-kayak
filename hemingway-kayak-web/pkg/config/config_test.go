package config

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestLoadDefaults(t *testing.T) {
	t.Chdir(t.TempDir()) // sin .env: se prueban los valores por defecto
	t.Setenv("DB_PASSWORD", "")
	t.Setenv("JWT_SECRET", "")
	t.Setenv("DB_SERVER", "")
	t.Setenv("DB_PORT", "")
	t.Setenv("DB_NAME", "")
	t.Setenv("CORS_ORIGINS", "")

	cfg := Load()

	if cfg.DBServer != "100.89.167.67" {
		t.Errorf("DBServer = %q, want default", cfg.DBServer)
	}
	if cfg.DBPort != 5499 {
		t.Errorf("DBPort = %d, want 5499", cfg.DBPort)
	}
	if cfg.Port != "8080" {
		t.Errorf("Port = %q, want 8080", cfg.Port)
	}
}

func TestValidateRequiresSecrets(t *testing.T) {
	t.Chdir(t.TempDir()) // sin .env
	t.Setenv("DB_PASSWORD", "")
	t.Setenv("JWT_SECRET", "")

	cfg := Load()
	err := cfg.Validate()
	if err == nil {
		t.Fatal("Validate() should fail without secrets")
	}
	for _, want := range []string{"DB_PASSWORD", "JWT_SECRET"} {
		if !strings.Contains(err.Error(), want) {
			t.Errorf("error should mention %s: %v", want, err)
		}
	}
}

func TestValidateOK(t *testing.T) {
	t.Setenv("DB_PASSWORD", "x")
	t.Setenv("JWT_SECRET", "y")

	if err := Load().Validate(); err != nil {
		t.Fatalf("Validate() error = %v", err)
	}
}

func TestDSN(t *testing.T) {
	cfg := &Config{DBServer: "h", DBPort: 5499, DBUser: "hemingway", DBPassword: "p", DBName: "hemingway_dev"}
	got := cfg.DSN(cfg.DBName)
	for _, want := range []string{"host=h", "port=5499", "dbname=hemingway_dev", "user=hemingway", "password=p", "sslmode=disable"} {
		if !strings.Contains(got, want) {
			t.Errorf("DSN %q missing %q", got, want)
		}
	}
}

func TestDSNDatabaseURLWins(t *testing.T) {
	url := "postgres://u:p@db:5432/app?sslmode=require"
	cfg := &Config{DatabaseURL: url, DBServer: "h", DBPort: 5499, DBUser: "x", DBPassword: "y", DBName: "z"}
	if got := cfg.DSN("ignored"); got != url {
		t.Errorf("DSN = %q, want DATABASE_URL %q", got, url)
	}
}

func TestDSNEscapesPassword(t *testing.T) {
	cfg := &Config{DBServer: "h", DBPort: 5499, DBUser: "u", DBPassword: `a'b\c`, DBName: "d"}
	got := cfg.DSN(cfg.DBName)
	if !strings.Contains(got, `password='a\'b\\c'`) {
		t.Errorf("DSN no escapa el password: %q", got)
	}
}

func TestLoadDotEnv(t *testing.T) {
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, ".env"), []byte("# comment\nHK_TEST_VAR=hello\n"), 0o600); err != nil {
		t.Fatalf("write .env: %v", err)
	}
	t.Chdir(dir)
	t.Setenv("HK_TEST_VAR", "")

	loadDotEnv()

	if got := os.Getenv("HK_TEST_VAR"); got != "hello" {
		t.Errorf("HK_TEST_VAR = %q, want %q", got, "hello")
	}
}

func TestLoadDotEnvDoesNotOverrideEnvironment(t *testing.T) {
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, ".env"), []byte("HK_TEST_VAR=from-file\n"), 0o600); err != nil {
		t.Fatalf("write .env: %v", err)
	}
	t.Chdir(dir)
	t.Setenv("HK_TEST_VAR", "from-env")

	loadDotEnv()

	if got := os.Getenv("HK_TEST_VAR"); got != "from-env" {
		t.Errorf("HK_TEST_VAR = %q, want %q", got, "from-env")
	}
}
