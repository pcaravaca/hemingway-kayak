package api

import (
	"io"
	"log"
	"net/http/httptest"
	"os"
	"strings"
	"sync"
	"testing"

	"github.com/gofiber/fiber/v2"

	"hemingway-kayak-web/internal/services"
	"hemingway-kayak-web/pkg/database"
)

var once sync.Once

// integrationApp levanta el servidor real (mismas rutas y middleware que en
// producción) contra la base de datos. Los tests de integración sólo corren si
// HK_INTEGRATION=1 y hay base disponible; si no, se saltan para que
// `go test ./...` siga siendo verde en una máquina sin Postgres.
func integrationApp(t *testing.T) *fiber.App {
	t.Helper()

	if os.Getenv("HK_INTEGRATION") == "" {
		t.Skip("HK_INTEGRATION no definida: se omiten los tests de integracion")
	}

	db, err := database.Connect()
	if err != nil {
		t.Skipf("no se pudo conectar a la DB de prueba: %v", err)
	}
	t.Cleanup(func() { _ = db.Close() })

	// Una sola vez por proceso: asegura esquema y seed. En la DB de desarrollo ya
	// están aplicados y queda en un no-op.
	once.Do(func() {
		if err := database.RunMigrations(); err != nil {
			t.Fatalf("migraciones: %v", err)
		}
		if err := database.SeedData(); err != nil {
			log.Printf("seed: %v", err)
		}
	})

	srv := NewServer(
		db,
		nil,
		services.NewUserService(db),
		services.NewTrackingService(db, nil),
		services.NewTournamentService(db),
		services.NewTripService(db),
		services.NewHistoryService(db),
		services.NewSpeciesService(db),
	)
	// NewServer sólo arma la app; las rutas se registran en Start(), que los tests
	// no llaman porque usan app.Test directamente.
	srv.setupRoutes()
	return srv.App()
}

func do(app *fiber.App, method, path, token, body string) (int, string) {
	req := httptest.NewRequest(method, path, strings.NewReader(body))
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}
	if body != "" {
		req.Header.Set("Content-Type", "application/json")
	}
	resp, err := app.Test(req)
	if err != nil {
		return 0, err.Error()
	}
	defer resp.Body.Close()
	raw, _ := io.ReadAll(resp.Body)
	return resp.StatusCode, string(raw)
}

func TestHealth(t *testing.T) {
	app := integrationApp(t)
	code, body := do(app, "GET", "/api/v1/health", "", "")
	if code != 200 || !strings.Contains(body, `"status":"ok"`) {
		t.Fatalf("health: code=%d body=%s", code, body)
	}
}

// TestSpeciesEndpoint cubre la DB de peces (migraciones 004/005/006): 19 especies
// y el filtro por hábitat.
func TestSpeciesEndpoint(t *testing.T) {
	app := integrationApp(t)

	code, body := do(app, "GET", "/api/v1/species", "", "")
	if code != 200 {
		t.Fatalf("species: code=%d body=%s", code, body)
	}
	for _, want := range []string{"Robalo / Snook", "Tarpon / Atlantic Tarpon", "Guapote / Rainbow bass"} {
		if !strings.Contains(body, want) {
			t.Errorf("falta %q en /species", want)
		}
	}
	if !strings.Contains(body, "TLPC") {
		t.Error("la respuesta de /species no trae la talla legal (incopesca)")
	}

	code, body = do(app, "GET", "/api/v1/species?habitat=rio", "", "")
	if code != 200 {
		t.Fatalf("species?habitat=rio: code=%d body=%s", code, body)
	}
	if strings.Contains(body, `"habitat":"mar"`) || strings.Contains(body, `"habitat":"estuario"`) {
		t.Error("el filtro por habitat no está filtrando")
	}

	code, _ = do(app, "GET", "/api/v1/species/00000000-0000-0000-0000-000000000000", "", "")
	if code != 404 {
		t.Errorf("especie inexistente: esperaba 404, llegó %d", code)
	}
}

func TestPublicReadEndpoints(t *testing.T) {
	app := integrationApp(t)

	for _, path := range []string{"/api/v1/tournaments", "/api/v1/tournaments/live", "/api/v1/trips", "/api/v1/spots"} {
		code, body := do(app, "GET", path, "", "")
		if code != 200 {
			t.Errorf("%s: esperaba 200, llegó %d (%s)", path, code, body)
		}
	}
}

func TestProtectedRoutesRequireToken(t *testing.T) {
	app := integrationApp(t)

	for _, path := range []string{"/api/v1/users/me", "/api/v1/users/me/stats", "/api/v1/tracking/followers"} {
		code, _ := do(app, "GET", path, "", "")
		if code != 401 {
			t.Errorf("%s sin token: esperaba 401, llegó %d", path, code)
		}
	}
}

func TestLoginRejectsBadPassword(t *testing.T) {
	app := integrationApp(t)

	code, body := do(app, "POST", "/api/v1/auth/login", "",
		`{"email":"peter@hemingway.com","password":"incorrecta"}`)
	if code != 401 {
		t.Fatalf("login con clave mala: esperaba 401, llegó %d (%s)", code, body)
	}
}
