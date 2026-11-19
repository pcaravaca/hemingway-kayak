package services

import (
	"io"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
)

// Sin clave el endpoint tiene que responder 503 con un mensaje claro, no un 500.
func TestWeatherWithoutKey(t *testing.T) {
	svc := &WeatherService{apiKey: "", client: nil}

	app := fiber.New()
	app.Get("/weather", svc.GetWeather)

	req := httptest.NewRequest("GET", "/weather?lat=9.9&lon=-84.8", nil)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("app.Test: %v", err)
	}
	if resp.StatusCode != 503 {
		t.Errorf("sin clave: esperaba 503, llegó %d", resp.StatusCode)
	}
	body, _ := io.ReadAll(resp.Body)
	if !contains(string(body), "OPENWEATHER_API_KEY") {
		t.Errorf("el error no explica qué falta: %s", body)
	}
}

// Con clave pero sin coordenadas tiene que ser 400.
func TestWeatherRequiresCoords(t *testing.T) {
	svc := &WeatherService{apiKey: "test-key", client: nil}

	app := fiber.New()
	app.Get("/weather", svc.GetWeather)

	req := httptest.NewRequest("GET", "/weather", nil)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("app.Test: %v", err)
	}
	if resp.StatusCode != 400 {
		t.Errorf("sin lat/lon: esperaba 400, llegó %d", resp.StatusCode)
	}
}

func contains(s, sub string) bool {
	for i := 0; i+len(sub) <= len(s); i++ {
		if s[i:i+len(sub)] == sub {
			return true
		}
	}
	return false
}
