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
		 WHERE user_id = @userId AND is_live = 1`,
		sql.Named("userId", userID),
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
		 OUTPUT CAST(INSERTED.id AS VARCHAR(36))
		 VALUES (@userId, @tripId, @lat, @lon, 1, GETDATE())`,
		sql.Named("userId", input.UserID),
		sql.Named("tripId", input.TripID),
		sql.Named("lat", input.Latitude),
		sql.Named("lon", input.Longitude),
	).Scan(&id)

	if err != nil {
		return c.Status(500).JSON(Error("Failed to start tracking"))
	}

	_, err = s.db.ExecContext(context.Background(),
		`INSERT INTO live_tracking (user_id, trip_id, latitude, longitude, timestamp, is_active)
		 VALUES (@userId, @tripId, @lat, @lon, GETDATE(), 1)`,
		sql.Named("userId", input.UserID),
		sql.Named("tripId", input.TripID),
		sql.Named("lat", input.Latitude),
		sql.Named("lon", input.Longitude),
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

	tx, err := s.db.BeginTx(context.Background(), nil)
	if err != nil {
		return c.Status(500).JSON(Error("Transaction error"))
	}

	_, err = tx.ExecContext(context.Background(),
		`UPDATE tracking_sessions SET 
		 is_live = 0, end_time = GETDATE(),
		 end_lat = @lat, end_lon = @lon
		 WHERE user_id = @userId AND is_live = 1`,
		sql.Named("userId", input.UserID),
		sql.Named("lat", input.Latitude),
		sql.Named("lon", input.Longitude),
	)

	_, err = tx.ExecContext(context.Background(),
		`UPDATE live_tracking SET is_active = 0
		 WHERE user_id = @userId`,
		sql.Named("userId", input.UserID),
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
		 latitude = @lat, longitude = @lon,
		 speed = @speed, heading = @heading, altitude = @alt,
		 timestamp = GETDATE()
		 WHERE user_id = @userId`,
		sql.Named("userId", input.UserID),
		sql.Named("lat", input.Latitude),
		sql.Named("lon", input.Longitude),
		sql.Named("speed", input.Speed),
		sql.Named("heading", input.Heading),
		sql.Named("alt", input.Altitude),
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
		 FROM live_tracking WHERE user_id = @userId`,
		sql.Named("userId", userID),
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
	rows, err := s.db.QueryContext(context.Background(),
		`SELECT CAST(user_id AS VARCHAR(36)), CAST(trip_id AS VARCHAR(36)), latitude, longitude, speed, heading, timestamp
		 FROM live_tracking WHERE is_active = 1`,
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

		if err := rows.Scan(&userID, &tripID, &lat, &lon, &speed, &heading, &timestamp); err != nil {
			continue
		}

		tracker := map[string]interface{}{}
		if userID.Valid {
			tracker["userId"] = userID.String
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
