package api

import (
	"context"
	"database/sql"
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/websocket/v2"
	"github.com/redis/go-redis/v9"

	"hemingway-kayak-web/internal/middleware"
	"hemingway-kayak-web/internal/services"
	"hemingway-kayak-web/internal/ws"
	"hemingway-kayak-web/pkg/auth"
	"hemingway-kayak-web/pkg/config"
)

type Server struct {
	app               *fiber.App
	db                *sql.DB
	redis             *redis.Client
	userService       *services.UserService
	trackingService   *services.TrackingService
	tournamentService *services.TournamentService
	tripService       *services.TripService
	historyService  *services.HistoryService
	speciesService  *services.SpeciesService
	uploadService   *services.UploadService
	weatherService  *services.WeatherService
	wsHub           *ws.Hub
	jwtSecret         []byte
	revoked           *auth.Revocation
}

func NewServer(
	db *sql.DB,
	redis *redis.Client,
	userService *services.UserService,
	trackingService *services.TrackingService,
	tournamentService *services.TournamentService,
	tripService *services.TripService,
	historyService *services.HistoryService,
	speciesService *services.SpeciesService,
	uploadService *services.UploadService,
	weatherService *services.WeatherService,
) *Server {
	cfg := config.Load()
	app := fiber.New(fiber.Config{
		AppName: "Hemingway Kayak API",
	})

	app.Use(recover.New())
	app.Use(logger.New(logger.Config{
		Format: "[${time}] ${status} - ${method} ${path} (${latency})\n",
	}))
	app.Use(cors.New(cors.Config{
		AllowOrigins: cfg.CORSOrigins,
		AllowMethods: "GET,POST,PUT,DELETE,OPTIONS",
		AllowHeaders: "Origin,Content-Type,Accept,Authorization",
	}))

	hub := ws.NewHub()
	// El tracking deja de escribir solo en la DB: ahora también empuja cada
	// punto a los browsers conectados por WebSocket.
	trackingService.SetBroadcaster(hub)
	// Igual para el ranking de torneo: join/leave/captura actualizan la tabla
	// en vivo de quien esté suscripto a la sala del torneo.
	tournamentService.SetBroadcaster(hub)

	return &Server{
		app:               app,
		db:                db,
		redis:             redis,
		userService:       userService,
		trackingService:   trackingService,
		tournamentService: tournamentService,
		tripService:       tripService,
		historyService:    historyService,
		speciesService:    speciesService,
		uploadService:     uploadService,
		weatherService:    weatherService,
		wsHub:             hub,
		jwtSecret:         []byte(cfg.JWTSecret),
		revoked:           auth.NewRevocation(),
	}
}

// App expone la aplicación Fiber; lo usan los tests de integración (app.Test).
func (s *Server) App() *fiber.App { return s.app }

// Start levanta el servidor (bloquea hasta que el listener falle).
func (s *Server) Start() error {
	s.setupRoutes()
	s.setupWebsocket()

	log.Printf("Server starting on port %s", config.Load().Port)
	return s.app.Listen(":" + config.Load().Port)
}

// Shutdown detiene el servidor de forma grácil.
func (s *Server) Shutdown() error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	return s.app.ShutdownWithContext(ctx)
}

func (s *Server) setupRoutes() {
	api := s.app.Group("/api/v1")

	api.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status": "ok",
			"time":   time.Now().Format(time.RFC3339),
		})
	})

	// Las fotos de las capturas se sirven estáticas desde el mismo origen.
	s.app.Static("/uploads", services.UploadDir())

	// --- Público ---
	api.Post("/auth/register", s.userService.Register)
	api.Post("/auth/login", s.userService.Login)

	// Lectura pública (feed)
	api.Get("/tournaments", s.tournamentService.GetTournaments)
	api.Get("/tournaments/live", s.tournamentService.GetLiveTournaments)
	api.Get("/tournaments/:id", s.tournamentService.GetTournament)
	api.Get("/tournaments/:id/participants", s.tournamentService.GetParticipants)
	api.Get("/tournaments/:id/ranking", s.tournamentService.GetRanking)
	api.Get("/tournaments/:id/catches", s.tournamentService.GetCatches)
	api.Get("/trips", s.tripService.GetTrips)
	api.Get("/spots", s.tripService.GetSharedSpots)
	api.Get("/weather", s.weatherService.GetWeather)
	api.Get("/species", s.speciesService.GetSpecies)
	api.Get("/species/:id", s.speciesService.GetSpeciesByID)

	// --- Autenticado ---
	sec := api.Group("", middleware.JWTAuth(s.jwtSecret, s.revoked))

	sec.Post("/auth/logout", s.userService.Logout)

	// Fotos de captura (multipart). Se sirven estáticas desde /uploads.
	sec.Post("/uploads", s.uploadService.UploadPhoto)

	sec.Get("/users/me", s.userService.GetMe)
	sec.Get("/users/me/stats", s.userService.GetMyStats)
	sec.Get("/users/:id", middleware.RequireOwnership("id"), s.userService.GetUser)
	sec.Put("/users/:id", middleware.RequireOwnership("id"), s.userService.UpdateUser)
	sec.Get("/users/:id/profile", middleware.RequireOwnership("id"), s.userService.GetProfile)

	sec.Get("/kayaks/", s.userService.GetKayaks)
	sec.Post("/kayaks/", s.userService.CreateKayak)
	sec.Put("/kayaks/:id", middleware.RequireResourceOwnership(s.db, "kayaks", "user_id", "id"), s.userService.UpdateKayak)
	sec.Delete("/kayaks/:id", middleware.RequireResourceOwnership(s.db, "kayaks", "user_id", "id"), s.userService.DeleteKayak)

	// Salidas propias: mis coordenadas, para el historial y los detalles.
	sec.Get("/trips/mine", s.tripService.GetMyTrips)
	sec.Get("/trips/:id", middleware.RequireResourceOwnership(s.db, "trips", "user_id", "id"), s.tripService.GetTrip)
	sec.Post("/trips/", s.tripService.CreateTrip)
	sec.Put("/trips/:id", middleware.RequireResourceOwnership(s.db, "trips", "user_id", "id"), s.tripService.UpdateTrip)
	sec.Delete("/trips/:id", middleware.RequireResourceOwnership(s.db, "trips", "user_id", "id"), s.tripService.DeleteTrip)
	sec.Get("/trips/:id/points", middleware.RequireResourceOwnership(s.db, "trips", "user_id", "id"), s.tripService.GetTripPoints)
	sec.Post("/trips/:id/points", middleware.RequireResourceOwnership(s.db, "trips", "user_id", "id"), s.tripService.AddTripPoint)

	sec.Post("/tournaments/", s.tournamentService.CreateTournament)
	sec.Put("/tournaments/:id", middleware.RequireResourceOwnership(s.db, "tournaments", "organizer_id", "id"), s.tournamentService.UpdateTournament)
	sec.Post("/tournaments/:id/join", s.tournamentService.JoinTournament)
	sec.Post("/tournaments/:id/leave", s.tournamentService.LeaveTournament)
	sec.Post("/tournaments/:id/catches", s.tournamentService.RecordCatch)

	sec.Get("/tracking/session/:userId", middleware.RequireOwnership("userId"), s.trackingService.GetActiveSession)
	sec.Post("/tracking/start", s.trackingService.StartTracking)
	sec.Post("/tracking/stop", s.trackingService.StopTracking)
	sec.Post("/tracking/point", s.trackingService.UpdateLocation)
	sec.Get("/tracking/follow/:userId", s.trackingService.FollowUser)
	sec.Get("/tracking/followers", s.trackingService.GetFollowers)

	sec.Get("/history/user/:userId", middleware.RequireOwnership("userId"), s.historyService.GetUserHistory)
	sec.Get("/history/user/:userId/trips", middleware.RequireOwnership("userId"), s.historyService.GetTripHistory)
	sec.Get("/history/user/:userId/tournaments", middleware.RequireOwnership("userId"), s.historyService.GetTournamentHistory)

	sec.Post("/spots/", s.tripService.CreateSpot)
	sec.Put("/spots/:id", middleware.RequireResourceOwnership(s.db, "fishing_spots", "user_id", "id"), s.tripService.UpdateSpot)
	sec.Delete("/spots/:id", middleware.RequireResourceOwnership(s.db, "fishing_spots", "user_id", "id"), s.tripService.DeleteSpot)
}

func (s *Server) setupWebsocket() {
	s.app.Use("/ws", websocket.New(func(c *websocket.Conn) {
		s.wsHub.HandleWebSocket(c)
	}))
	go s.wsHub.Run()
}
