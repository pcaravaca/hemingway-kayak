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

type TrackingService struct {
	db    *sql.DB
	redis *redis.Client
}

func NewTrackingService(db *sql.DB, redis *redis.Client) *TrackingService {
	return &TrackingService{db: db, redis: redis}
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

	var id string
	err := s.db.QueryRowContext(context.Background(),
		`INSERT INTO tracking_sessions (user_id, trip_id, start_lat, start_lon, is_live, start_time)
		 VALUES ($1, $2, $3, $4, TRUE, now())
		 RETURNING id::text`,
		input.UserID,
		input.TripID,
		input.Latitude,
		input.Longitude,
	).Scan(&id)

	if err != nil {
		return c.Status(500).JSON(Error("Failed to start tracking"))
	}

	_, err = s.db.ExecContext(context.Background(),
		`INSERT INTO live_tracking (user_id, trip_id, latitude, longitude, timestamp, is_active)
		 VALUES ($1, $2, $3, $4, now(), 1)`,
		input.UserID,
		input.TripID,
		input.Latitude,
		input.Longitude,
	)

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

	_, err := s.db.ExecContext(context.Background(),
		`UPDATE live_tracking SET 
		 latitude = $1, longitude = $2,
		 speed = $3, heading = $4, altitude = $5,
		 timestamp = now()
		 WHERE user_id = $6`,
		input.Latitude,
		input.Longitude,
		input.Speed,
		input.Heading,
		input.Altitude,
		input.UserID,
	)

	if err != nil {
		return c.Status(500).JSON(Error("Failed to update"))
	}

	ctx := context.Background()
	data, _ := json.Marshal(fiber.Map{
		"userId": input.UserID, "latitude": input.Latitude,
		"longitude": input.Longitude, "speed": input.Speed,
		"heading": input.Heading, "timestamp": time.Now(),
	})
	if s.redis != nil {
		if err := s.redis.Publish(ctx, "tracking:live", data).Err(); err != nil {
			log.Printf("Warning: failed to publish tracking update: %v", err)
		}
	}

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
