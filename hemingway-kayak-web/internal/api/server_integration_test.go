package api

import (
	"encoding/json"
	"bytes"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"

	"hemingway-kayak-web/internal/services"
	"hemingway-kayak-web/pkg/auth"
	"hemingway-kayak-web/pkg/config"
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

	// Cada prueba sube las fotos a su propia carpeta temporal.
	t.Setenv("UPLOAD_DIR", t.TempDir())

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
		services.NewUploadService(),
		services.NewWeatherService(),
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

	// El feed público no expone coordenadas: son datos de ubicación.
	_, body := do(app, "GET", "/api/v1/trips", "", "")
	if strings.Contains(body, "startLat") || strings.Contains(body, "startLon") {
		t.Error("el feed público de salidas está exponiendo coordenadas")
	}

	// La lista de torneos tiene que traer el conteo de participantes: el detalle
	// muestra el ranking y la lista decía 0 siempre.
	_, body = do(app, "GET", "/api/v1/tournaments", "", "")
	var tournaments []map[string]interface{}
	if err := json.Unmarshal([]byte(body), &tournaments); err != nil {
		t.Fatalf("/tournaments no es JSON: %v", err)
	}
	if len(tournaments) == 0 {
		t.Fatal("/tournaments no devolvió ningún torneo")
	}
	for _, t0 := range tournaments {
		count, ok := t0["participantsCount"].(float64)
		if !ok {
			t.Errorf("%s: participantsCount ausente o no numérico (%T)", t0["name"], t0["participantsCount"])
			continue
		}
		if count < 0 {
			t.Errorf("%s: participantsCount negativo", t0["name"])
		}
	}
}

func TestMyTripsRequiresToken(t *testing.T) {
	app := integrationApp(t)

	code, _ := do(app, "GET", "/api/v1/trips/mine", "", "")
	if code != 401 {
		t.Errorf("/trips/mine sin token: esperaba 401, llegó %d", code)
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

func testToken(t *testing.T, userID, email string) string {
	t.Helper()
	token, err := auth.Generate([]byte(config.Load().JWTSecret), userID, email, time.Hour)
	if err != nil {
		t.Fatalf("no se pudo firmar el token de prueba: %v", err)
	}
	return token
}

// multipartRequest arma un POST /uploads con un archivo.
func multipartRequest(t *testing.T, field, filename string, content []byte) *http.Request {
	t.Helper()
	var buf bytes.Buffer
	w := multipart.NewWriter(&buf)
	part, err := w.CreateFormFile(field, filename)
	if err != nil {
		t.Fatalf("CreateFormFile: %v", err)
	}
	if _, err := part.Write(content); err != nil {
		t.Fatalf("write: %v", err)
	}
	_ = w.Close()
	req := httptest.NewRequest("POST", "/api/v1/uploads", &buf)
	req.Header.Set("Content-Type", w.FormDataContentType())
	return req
}

// El PNG más chico posible (1x1 transparente); http.DetectContentType lo reconoce.
var tinyPNG = []byte{
	0x89, 0x50, 0x4e, 0x47, 0x0d, 0x0a, 0x1a, 0x0a, 0x00, 0x00, 0x00, 0x0d,
	0x49, 0x48, 0x44, 0x52, 0x00, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x01,
	0x08, 0x06, 0x00, 0x00, 0x00, 0x1f, 0x15, 0xc4, 0x89, 0x00, 0x00, 0x00,
	0x0a, 0x49, 0x44, 0x41, 0x54, 0x78, 0x9c, 0x63, 0x00, 0x01, 0x00, 0x00,
	0x05, 0x00, 0x01, 0x0d, 0x0a, 0x2d, 0xb4, 0x00, 0x00, 0x00, 0x00, 0x49,
	0x45, 0x4e, 0x44, 0xae, 0x42, 0x60, 0x82,
}

func TestUploadPhoto(t *testing.T) {
	app := integrationApp(t)
	token := testToken(t, "00000000-0000-0000-0000-000000000001", "peter@hemingway.com")

	// Sin sesión no se puede subir.
	code, _ := do(app, "POST", "/api/v1/uploads", "", "")
	if code != 401 {
		t.Errorf("upload sin token: esperaba 401, llegó %d", code)
	}

	// Content-Type Trampa: el contenido manda, no el header.
	req := multipartRequest(t, "photo", "foto.txt", []byte("esto no es una imagen"))
	req.Header.Set("Authorization", "Bearer "+token)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("app.Test: %v", err)
	}
	resp.Body.Close()
	if resp.StatusCode != 415 {
		t.Errorf("upload de texto: esperaba 415, llegó %d", resp.StatusCode)
	}

	// Una foto válida se guarda y queda servida.
	req = multipartRequest(t, "photo", "captura.png", tinyPNG)
	req.Header.Set("Authorization", "Bearer "+token)
	resp, err = app.Test(req)
	if err != nil {
		t.Fatalf("app.Test: %v", err)
	}
	defer resp.Body.Close()
	raw, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != 201 {
		t.Fatalf("upload válido: esperaba 201, llegó %d (%s)", resp.StatusCode, raw)
	}
	if !strings.Contains(string(raw), "/uploads/") {
		t.Errorf("la respuesta no devuelve la URL: %s", raw)
	}
}
