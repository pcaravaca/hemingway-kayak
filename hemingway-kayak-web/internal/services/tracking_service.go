package services

import (
	"context"
	"database/sql"
	"encoding/json"
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/redis/go-redis/v9"

	"hemingway-kayak-web/internal/middleware"
)

// wsBroadcaster es lo único que TrackingService necesita del hub de WebSocket.
// Se define acá para no acoplar services con internal/ws.
type wsBroadcaster interface {
	BroadcastTo(room string, eventType string, data interface{})
}

type TrackingService struct {
	db    *sql.DB
	redis *redis.Client
	hub   wsBroadcaster
}

func NewTrackingService(db *sql.DB, redis *redis.Client) *TrackingService {
	return &TrackingService{db: db, redis: redis}
}

// SetBroadcaster engancha el hub de WebSocket. Sin esto, la posición en vivo se
// guarda en la DB pero ningún browser se entera.
func (s *TrackingService) SetBroadcaster(hub wsBroadcaster) {
	s.hub = hub
}

// push transmite la posición a todos los clientes conectados.
func (s *TrackingService) push(event string, payload fiber.Map) {
	if s.hub != nil {
		s.hub.BroadcastTo("", event, payload)
	}
}

func (s *TrackingService) GetActiveSession(c *fiber.Ctx) error {
	userID := c.Params("userId")

	var session struct {
		ID            string
		UserID        string
		TripID        sql.NullString
		StartTime     time.Time
		StartLat      float64
		StartLon      float64
		TotalDistance float64
		IsLive        bool
		ViewerCount   int
	}

	err := s.db.QueryRowContext(context.Background(),
		`SELECT id, user_id, trip_id, start_time, start_lat, start_lon, 
		        total_distance, is_live, viewer_count
		 FROM tracking_sessions
		 WHERE user_id = $1 AND is_live = TRUE`,
		userID,
	).Scan(&session.ID, &session.UserID, &session.TripID, &session.StartTime,
		&session.StartLat, &session.StartLon, &session.TotalDistance,
		&session.IsLive, &session.ViewerCount)

	if err == sql.ErrNoRows {
		return c.Status(404).JSON(Error("No active session"))
	}
	if err != nil {
		return c.Status(500).JSON(Error("Database error"))
	}

	return c.JSON(fiber.Map{
		"id":            session.ID,
		"userId":        session.UserID,
		"tripId":        session.TripID.String,
		"startTime":     session.StartTime,
		"startLat":      session.StartLat,
		"startLon":      session.StartLon,
		"totalDistance": session.TotalDistance,
		"isLive":        session.IsLive,
		"viewerCount":   session.ViewerCount,
	})
}

func (s *TrackingService) StartTracking(c *fiber.Ctx) error {
	type Input struct {
		UserID    string  `json:"userId"`
		TripID    string  `json:"tripId"`
		Latitude  float64 `json:"latitude"`
		Longitude float64 `json:"longitude"`
	}

	var input Input
	if err := c.BodyParser(&input); err != nil {
		return c.Status(400).JSON(Error("Invalid request"))
	}
	input.UserID = middleware.UserID(c)

	// tripId es opcional: si viene vacío o null, se guarda NULL y no un '' que
	// Postgres rechaza como UUID (eso devolvía un 500 sin mensaje).
	var tripID sql.NullString
	if input.TripID != "" {
		tripID = sql.NullString{String: input.TripID, Valid: true}
	}

	var id string
	err := s.db.QueryRowContext(context.Background(),
		`INSERT INTO tracking_sessions (user_id, trip_id, start_lat, start_lon, is_live, start_time)
		 VALUES ($1, $2, $3, $4, TRUE, now())
		 RETURNING id::text`,
		input.UserID,
		tripID,
		input.Latitude,
		input.Longitude,
	).Scan(&id)

	if err != nil {
		return c.Status(500).JSON(Error("Failed to start tracking"))
	}

	_, err = s.db.ExecContext(context.Background(),
		`INSERT INTO live_tracking (user_id, trip_id, latitude, longitude, timestamp, is_active)
		 VALUES ($1, $2, $3, $4, now(), 1)
		 ON CONFLICT (user_id) DO UPDATE SET
		   trip_id = EXCLUDED.trip_id,
		   latitude = EXCLUDED.latitude,
		   longitude = EXCLUDED.longitude,
		   timestamp = now(),
		   is_active = TRUE`,
		input.UserID,
		tripID,
		input.Latitude,
		input.Longitude,
	)

	s.push("tracking", fiber.Map{
		"userId": input.UserID, "latitude": input.Latitude,
		"longitude": input.Longitude, "speed": 0, "heading": 0, "timestamp": time.Now(),
	})

	return c.Status(201).JSON(fiber.Map{"id": id, "success": true})
}

func (s *TrackingService) StopTracking(c *fiber.Ctx) error {
	type Input struct {
		UserID    string  `json:"userId"`
		Latitude  float64 `json:"latitude"`
		Longitude float64 `json:"longitude"`
	}

	var input Input
	if err := c.BodyParser(&input); err != nil {
		return c.Status(400).JSON(Error("Invalid request"))
	}
	input.UserID = middleware.UserID(c)

	tx, err := s.db.BeginTx(context.Background(), nil)
	if err != nil {
		return c.Status(500).JSON(Error("Transaction error"))
	}

	_, err = tx.ExecContext(context.Background(),
		`UPDATE tracking_sessions SET 
		 is_live = FALSE, end_time = now(),
		 end_lat = $1, end_lon = $2
		 WHERE user_id = $3 AND is_live = TRUE`,
		input.Latitude,
		input.Longitude,
		input.UserID,
	)

	_, err = tx.ExecContext(context.Background(),
		`UPDATE live_tracking SET is_active = FALSE
		 WHERE user_id = $1`,
		input.UserID,
	)

	if err != nil {
		tx.Rollback()
		return c.Status(500).JSON(Error("Failed to stop"))
	}

	tx.Commit()

	s.push("tracking-stop", fiber.Map{
		"userId": input.UserID, "timestamp": time.Now(),
	})

	return c.JSON(Success("Tracking stopped"))
}

func (s *TrackingService) UpdateLocation(c *fiber.Ctx) error {
	type Input struct {
		UserID    string  `json:"userId"`
		Latitude  float64 `json:"latitude"`
		Longitude float64 `json:"longitude"`
		Speed     float64 `json:"speed"`
		Heading   float64 `json:"heading"`
		Altitude  float64 `json:"altitude"`
	}

	var input Input
	if err := c.BodyParser(&input); err != nil {
		return c.Status(400).JSON(Error("Invalid request"))
	}
	input.UserID = middleware.UserID(c)

	// Upsert: antes era un UPDATE puro, así que el primer punto de una sesión
	// (si no pasó por /tracking/start) no dejaba fila y el mapa nunca lo veía.
	_, err := s.db.ExecContext(context.Background(),
		`INSERT INTO live_tracking
		   (user_id, trip_id, latitude, longitude, speed, heading, altitude, timestamp, is_active)
		 VALUES ($1, NULL, $2, $3, $4, $5, $6, now(), TRUE)
		 ON CONFLICT (user_id) DO UPDATE SET
		   latitude = EXCLUDED.latitude,
		   longitude = EXCLUDED.longitude,
		   speed = EXCLUDED.speed,
		   heading = EXCLUDED.heading,
		   altitude = EXCLUDED.altitude,
		   timestamp = now(),
		   is_active = TRUE`,
		input.UserID,
		input.Latitude,
		input.Longitude,
		input.Speed,
		input.Heading,
		input.Altitude,
	)

	if err != nil {
		return c.Status(500).JSON(Error("Failed to update"))
	}

	payload := fiber.Map{
		"userId": input.UserID, "latitude": input.Latitude,
		"longitude": input.Longitude, "speed": input.Speed,
		"heading": input.Heading, "timestamp": time.Now(),
	}

	ctx := context.Background()
	data, _ := json.Marshal(payload)
	if s.redis != nil {
		if err := s.redis.Publish(ctx, "tracking:live", data).Err(); err != nil {
			log.Printf("Warning: failed to publish tracking update: %v", err)
		}
	}

	// Es el que de verdad actualiza los marcadores del mapa en el browser.
	s.push("tracking", payload)

	return c.JSON(Success("Location updated"))
}

func (s *TrackingService) FollowUser(c *fiber.Ctx) error {
	userID := c.Params("userId")

	var tracking struct {
		UserID    string
		TripID    sql.NullString
		Latitude  float64
		Longitude float64
		Speed     float64
		Heading   float64
		Timestamp time.Time
		IsActive  bool
	}

	err := s.db.QueryRowContext(context.Background(),
		`SELECT user_id, trip_id, latitude, longitude, speed, heading, timestamp, is_active
		 FROM live_tracking WHERE user_id = $1`,
		userID,
	).Scan(&tracking.UserID, &tracking.TripID, &tracking.Latitude,
		&tracking.Longitude, &tracking.Speed, &tracking.Heading,
		&tracking.Timestamp, &tracking.IsActive)

	if err == sql.ErrNoRows || !tracking.IsActive {
		return c.Status(404).JSON(Error("User not sharing location"))
	}
	if err != nil {
		return c.Status(500).JSON(Error("Database error"))
	}

	return c.JSON(fiber.Map{
		"userId":    tracking.UserID,
		"tripId":    tracking.TripID.String,
		"latitude":  tracking.Latitude,
		"longitude": tracking.Longitude,
		"speed":     tracking.Speed,
		"heading":   tracking.Heading,
		"timestamp": tracking.Timestamp,
	})
}

func (s *TrackingService) GetFollowers(c *fiber.Ctx) error {
	// Se une con users para que el panel y el home muestren un nombre legible
	// en vez del UUID del kayaker.
	rows, err := s.db.QueryContext(context.Background(),
		`SELECT CAST(lt.user_id AS VARCHAR(36)), CAST(lt.trip_id AS VARCHAR(36)),
		        lt.latitude, lt.longitude, lt.speed, lt.heading, lt.timestamp,
		        COALESCE(u.display_name, u.username, '')
		 FROM live_tracking lt
		 LEFT JOIN users u ON u.id = lt.user_id
		 WHERE lt.is_active = TRUE`,
	)
	if err != nil {
		return c.Status(500).JSON(Error("Database error"))
	}
	defer rows.Close()

	trackers := []map[string]interface{}{}
	for rows.Next() {
		var userID, tripID sql.NullString
		var lat, lon, speed, heading sql.NullFloat64
		var timestamp sql.NullTime
		var name sql.NullString

		if err := rows.Scan(&userID, &tripID, &lat, &lon, &speed, &heading, &timestamp, &name); err != nil {
			continue
		}

		tracker := map[string]interface{}{}
		if userID.Valid {
			tracker["userId"] = userID.String
		}
		if name.Valid {
			tracker["displayName"] = name.String
		}
		if lat.Valid {
			tracker["latitude"] = lat.Float64
		}
		if lon.Valid {
			tracker["longitude"] = lon.Float64
		}
		if timestamp.Valid {
			tracker["timestamp"] = timestamp.Time
		}
		if tripID.Valid {
			tracker["tripId"] = tripID.String
		}
		if speed.Valid {
			tracker["speed"] = speed.Float64
		}
		if heading.Valid {
			tracker["heading"] = heading.Float64
		}
		trackers = append(trackers, tracker)
	}

	return c.JSON(trackers)
}
