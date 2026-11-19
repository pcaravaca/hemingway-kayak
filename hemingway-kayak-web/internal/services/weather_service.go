package services

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"hemingway-kayak-web/pkg/config"

	"github.com/gofiber/fiber/v2"
)

// WeatherService expone el clima actual de OpenWeatherMap sin mandar la clave al
// navegador: el frontend llama /api/v1/weather y el backend hace el request.
type WeatherService struct {
	apiKey string
	client *http.Client
}

func NewWeatherService() *WeatherService {
	return &WeatherService{
		apiKey: config.Load().OpenWeatherAPIKey,
		client: &http.Client{Timeout: 10 * time.Second},
	}
}

type openWeatherResponse struct {
	Main struct {
		Temp      float64 `json:"temp"`
		FeelsLike float64 `json:"feels_like"`
		Humidity  int     `json:"humidity"`
		Pressure  int     `json:"pressure"`
	} `json:"main"`
	Wind struct {
		Speed float64 `json:"speed"`
		Deg   float64 `json:"deg"`
	} `json:"wind"`
	Clouds struct {
		All int `json:"all"`
	} `json:"clouds"`
	Weather []struct {
		Description string `json:"description"`
		Icon        string `json:"icon"`
	} `json:"weather"`
	Sys struct {
		Sunrise int64 `json:"sunrise"`
		Sunset  int64 `json:"sunset"`
	} `json:"sys"`
	Name string `json:"name"`
}

// GetWeather devuelve el clima actual de un punto. Si el servidor no tiene
// OPENWEATHER_API_KEY responde 503 con un mensaje claro (no un 500 pelado).
func (s *WeatherService) GetWeather(c *fiber.Ctx) error {
	if s.apiKey == "" {
		return c.Status(503).JSON(Error("El servidor no tiene OPENWEATHER_API_KEY configurada"))
	}

	lat := c.Query("lat")
	lon := c.Query("lon")
	if lat == "" || lon == "" {
		return c.Status(400).JSON(Error("Faltan los parámetros lat y lon"))
	}

	url := fmt.Sprintf(
		"https://api.openweathermap.org/data/2.5/weather?lat=%s&lon=%s&units=metric&lang=es&appid=%s",
		lat, lon, s.apiKey,
	)

	resp, err := s.client.Get(url)
	if err != nil {
		return c.Status(502).JSON(Error("No se pudo consultar OpenWeatherMap"))
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return c.Status(502).JSON(Error("OpenWeatherMap respondió " + resp.Status))
	}

	var raw openWeatherResponse
	if err := json.NewDecoder(resp.Body).Decode(&raw); err != nil {
		return c.Status(502).JSON(Error("No se pudo interpretar la respuesta de OpenWeatherMap"))
	}

	out := map[string]interface{}{
		"latitude":  lat,
		"longitude": lon,
		"place":     raw.Name,
		"temperature": raw.Main.Temp,
		"feelsLike":   raw.Main.FeelsLike,
		"humidity":    raw.Main.Humidity,
		"pressure":    raw.Main.Pressure,
		"windSpeed":   raw.Wind.Speed,
		"windDeg":     raw.Wind.Deg,
		"cloudCover":  raw.Clouds.All,
		"sunrise":     raw.Sys.Sunrise,
		"sunset":      raw.Sys.Sunset,
		"fetchedAt":   time.Now().UTC(),
	}
	if len(raw.Weather) > 0 {
		out["description"] = raw.Weather[0].Description
		out["icon"] = raw.Weather[0].Icon
	}

	return c.JSON(out)
}
