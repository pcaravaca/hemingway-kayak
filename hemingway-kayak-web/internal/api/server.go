package api

import (
	"log"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/websocket/v2"
	"github.com/redis/go-redis/v9"

	"hemingway-kayak-web/internal/services"
	"hemingway-kayak-web/internal/ws"
)

type Server struct {
	app               *fiber.App
	redis             *redis.Client
	userService       *services.UserService
	trackingService   *services.TrackingService
	tournamentService *services.TournamentService
	tripService       *services.TripService
	historyService    *services.HistoryService
	wsHub             *ws.Hub
}

func NewServer(
	redis *redis.Client,
	userService *services.UserService,
	trackingService *services.TrackingService,
	tournamentService *services.TournamentService,
	tripService *services.TripService,
	historyService *services.HistoryService,
) *Server {
	app := fiber.New(fiber.Config{
		AppName: "Hemingway Kayak API",
	})

	app.Use(recover.New())
	app.Use(logger.New(logger.Config{
		Format: "[${time}] ${status} - ${method} ${path} (${latency})\n",
	}))
	app.Use(cors.New(cors.Config{
		AllowOrigins: os.Getenv("CORS_ORIGINS"),
		AllowMethods: "GET,POST,PUT,DELETE,OPTIONS",
		AllowHeaders: "Origin,Content-Type,Accept,Authorization",
	}))

	return &Server{
		app:               app,
		redis:             redis,
		userService:       userService,
		trackingService:   trackingService,
		tournamentService: tournamentService,
		tripService:       tripService,
		historyService:    historyService,
		wsHub:             ws.NewHub(),
	}
}

func (s *Server) Start() error {
	s.setupRoutes()
	s.setupWebsocket()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server starting on port %s", port)
	return s.app.Listen(":" + port)
}

func (s *Server) setupRoutes() {
	api := s.app.Group("/api/v1")

	api.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status": "ok",
			"time":   time.Now().Format(time.RFC3339),
		})
	})

	auth := api.Group("/auth")
	auth.Post("/register", s.userService.Register)
	auth.Post("/login", s.userService.Login)
	auth.Post("/logout", s.userService.Logout)

	users := api.Group("/users")
	users.Get("/:id", s.userService.GetUser)
	users.Put("/:id", s.userService.UpdateUser)
	users.Get("/:id/profile", s.userService.GetProfile)

	kayaks := api.Group("/kayaks")
	kayaks.Get("/", s.userService.GetKayaks)
	kayaks.Post("/", s.userService.CreateKayak)
	kayaks.Put("/:id", s.userService.UpdateKayak)
	kayaks.Delete("/:id", s.userService.DeleteKayak)

	trips := api.Group("/trips")
	trips.Get("/", s.tripService.GetTrips)
	trips.Get("/:id", s.tripService.GetTrip)
	trips.Post("/", s.tripService.CreateTrip)
	trips.Put("/:id", s.tripService.UpdateTrip)
	trips.Delete("/:id", s.tripService.DeleteTrip)
	trips.Get("/:id/points", s.tripService.GetTripPoints)
	trips.Post("/:id/points", s.tripService.AddTripPoint)

	tournaments := api.Group("/tournaments")
	tournaments.Get("/", s.tournamentService.GetTournaments)
	tournaments.Get("/live", s.tournamentService.GetLiveTournaments)
	tournaments.Get("/:id", s.tournamentService.GetTournament)
	tournaments.Post("/", s.tournamentService.CreateTournament)
	tournaments.Put("/:id", s.tournamentService.UpdateTournament)
	tournaments.Get("/:id/participants", s.tournamentService.GetParticipants)
	tournaments.Get("/:id/ranking", s.tournamentService.GetRanking)
	tournaments.Post("/:id/join", s.tournamentService.JoinTournament)
	tournaments.Post("/:id/leave", s.tournamentService.LeaveTournament)
	tournaments.Post("/:id/catches", s.tournamentService.RecordCatch)
	tournaments.Get("/:id/catches", s.tournamentService.GetCatches)

	tracking := api.Group("/tracking")
	tracking.Get("/session/:userId", s.trackingService.GetActiveSession)
	tracking.Post("/start", s.trackingService.StartTracking)
	tracking.Post("/stop", s.trackingService.StopTracking)
	tracking.Post("/point", s.trackingService.UpdateLocation)
	tracking.Get("/follow/:userId", s.trackingService.FollowUser)
	tracking.Get("/followers", s.trackingService.GetFollowers)

	history := api.Group("/history")
	history.Get("/user/:userId", s.historyService.GetUserHistory)
	history.Get("/user/:userId/trips", s.historyService.GetTripHistory)
	history.Get("/user/:userId/tournaments", s.historyService.GetTournamentHistory)

	spots := api.Group("/spots")
	spots.Get("/", s.tripService.GetSharedSpots)
	spots.Post("/", s.tripService.CreateSpot)
	spots.Put("/:id", s.tripService.UpdateSpot)
	spots.Delete("/:id", s.tripService.DeleteSpot)
}

func (s *Server) setupWebsocket() {
	s.app.Use("/ws", websocket.New(func(c *websocket.Conn) {
		s.wsHub.HandleWebSocket(c)
	}))
	go s.wsHub.Run()
}
