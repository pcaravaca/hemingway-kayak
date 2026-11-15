package services

import (
	"context"
	"database/sql"
	"time"

	"github.com/gofiber/fiber/v2"

	"hemingway-kayak-web/internal/middleware"
)

type TripService struct {
	db *sql.DB
}

func NewTripService(db *sql.DB) *TripService {
	return &TripService{db: db}
}

func (s *TripService) GetTrips(c *fiber.Ctx) error {
	userID := c.Query("userId")
	tripType := c.Query("type")
	status := c.Query("status")

	query := `SELECT CAST(id AS VARCHAR(36)), CAST(user_id AS VARCHAR(36)), title, description, triptype, status, start_time, end_time, 
	          start_lat, start_lon, end_lat, end_lon, distance, created_at 
	          FROM trips WHERE 1=1`
	var args []interface{}

	if userID != "" {
		query += " AND user_id = @userId"
		args = append(args, sql.Named("userId", userID))
	}
	if tripType != "" {
		query += " AND triptype = @type"
		args = append(args, sql.Named("type", tripType))
	}
	if status != "" {
		query += " AND status = @status"
		args = append(args, sql.Named("status", status))
	}

	query += " ORDER BY created_at DESC"

	rows, err := s.db.QueryContext(context.Background(), query, args...)
	if err != nil {
		return c.Status(500).JSON(Error("Database error"))
	}
	defer rows.Close()

	trips := []map[string]interface{}{}
	for rows.Next() {
		var id, userID sql.NullString
		var title, description, tripType, status sql.NullString
		var startTime, endTime, createdAt sql.NullTime
		var startLat, startLon, endLat, endLon, distance sql.NullFloat64

		if err := rows.Scan(&id, &userID, &title, &description, &tripType, &status,
			&startTime, &endTime, &startLat, &startLon, &endLat, &endLon, &distance, &createdAt); err != nil {
			continue
		}

		trip := map[string]interface{}{}
		if id.Valid {
			trip["id"] = id.String
		}
		if userID.Valid {
			trip["userId"] = userID.String
		}
		if title.Valid {
			trip["title"] = title.String
		}
		if description.Valid {
			trip["description"] = description.String
		}
		if tripType.Valid {
			trip["type"] = tripType.String
		}
		if status.Valid {
			trip["status"] = status.String
		}
		if startTime.Valid {
			trip["startTime"] = startTime.Time
		}
		if endTime.Valid {
			trip["endTime"] = endTime.Time
		}
		if startLat.Valid {
			trip["startLat"] = startLat.Float64
		}
		if startLon.Valid {
			trip["startLon"] = startLon.Float64
		}
		if distance.Valid {
			trip["distance"] = distance.Float64
		}
		if createdAt.Valid {
			trip["createdAt"] = createdAt.Time
		}
		trips = append(trips, trip)
	}

	return c.JSON(trips)
}

func (s *TripService) GetTrip(c *fiber.Ctx) error {
	id := c.Params("id")

	var trip struct {
		ID          string
		UserID      string
		Title       sql.NullString
		Description sql.NullString
		Type        sql.NullString
		Status      sql.NullString
		StartTime   sql.NullTime
		EndTime     sql.NullTime
		StartLat    sql.NullFloat64
		StartLon    sql.NullFloat64
		EndLat      sql.NullFloat64
		EndLon      sql.NullFloat64
		Distance    sql.NullFloat64
		CreatedAt   sql.NullTime
	}

	err := s.db.QueryRowContext(context.Background(),
		`SELECT id, user_id, title, description, triptype, status, start_time, end_time,
		 start_lat, start_lon, end_lat, end_lon, distance, created_at 
		 FROM trips WHERE id = @id`,
		sql.Named("id", id),
	).Scan(&trip.ID, &trip.UserID, &trip.Title, &trip.Description, &trip.Type, &trip.Status,
		&trip.StartTime, &trip.EndTime, &trip.StartLat, &trip.StartLon, &trip.EndLat, &trip.EndLon,
		&trip.Distance, &trip.CreatedAt)

	if err == sql.ErrNoRows {
		return c.Status(404).JSON(Error("Trip not found"))
	}
	if err != nil {
		return c.Status(500).JSON(Error("Database error"))
	}

	return c.JSON(fiber.Map{
		"id":          trip.ID,
		"userId":      trip.UserID,
		"title":       trip.Title.String,
		"description": trip.Description.String,
		"type":        trip.Type.String,
		"status":      trip.Status.String,
		"startTime":   trip.StartTime.Time,
		"endTime":     trip.EndTime.Time,
		"startLat":    trip.StartLat.Float64,
		"startLon":    trip.StartLon.Float64,
		"endLat":      trip.EndLat.Float64,
		"endLon":      trip.EndLon.Float64,
		"distance":    trip.Distance.Float64,
		"createdAt":   trip.CreatedAt.Time,
	})
}

func (s *TripService) CreateTrip(c *fiber.Ctx) error {
	type TripInput struct {
		UserID      string  `json:"userId"`
		Title       string  `json:"title"`
		Description string  `json:"description"`
		Type        string  `json:"type"`
		StartLat    float64 `json:"startLat"`
		StartLon    float64 `json:"startLon"`
	}

	var input TripInput
	if err := c.BodyParser(&input); err != nil {
		return c.Status(400).JSON(Error("Invalid request body"))
	}
	input.UserID = middleware.UserID(c)

	var id string
	err := s.db.QueryRowContext(context.Background(),
		`INSERT INTO trips (user_id, title, description, triptype, start_lat, start_lon, status, start_time)
		 OUTPUT CAST(INSERTED.id AS VARCHAR(36))
		 VALUES (@userId, @title, @description, @type, @startLat, @startLon, 'planned', GETDATE())`,
		sql.Named("userId", input.UserID),
		sql.Named("title", input.Title),
		sql.Named("description", input.Description),
		sql.Named("type", input.Type),
		sql.Named("startLat", input.StartLat),
		sql.Named("startLon", input.StartLon),
	).Scan(&id)

	if err != nil {
		return c.Status(500).JSON(Error("Failed to create trip"))
	}

	return c.Status(201).JSON(fiber.Map{"id": id, "success": true})
}

func (s *TripService) UpdateTrip(c *fiber.Ctx) error {
	id := c.Params("id")

	type UpdateInput struct {
		Title       *string  `json:"title"`
		Description *string  `json:"description"`
		Status      *string  `json:"status"`
		EndLat      *float64 `json:"endLat"`
		EndLon      *float64 `json:"endLon"`
		Distance    *float64 `json:"distance"`
	}

	var input UpdateInput
	if err := c.BodyParser(&input); err != nil {
		return c.Status(400).JSON(Error("Invalid request body"))
	}

	query := `UPDATE trips SET 
		title = ISNULL(@title, title),
		description = ISNULL(@description, description),
		status = ISNULL(@status, status),
		end_lat = ISNULL(@endLat, end_lat),
		end_lon = ISNULL(@endLon, end_lon),
		distance = ISNULL(@distance, distance),
		end_time = CASE WHEN @status IN ('completed', 'cancelled') THEN GETDATE() ELSE end_time END
		WHERE id = @id`

	_, err := s.db.ExecContext(context.Background(), query,
		sql.Named("id", id),
		sql.Named("title", input.Title),
		sql.Named("description", input.Description),
		sql.Named("status", input.Status),
		sql.Named("endLat", input.EndLat),
		sql.Named("endLon", input.EndLon),
		sql.Named("distance", input.Distance),
	)

	if err != nil {
		return c.Status(500).JSON(Error("Failed to update trip"))
	}

	return c.JSON(Success("Trip updated"))
}

func (s *TripService) DeleteTrip(c *fiber.Ctx) error {
	id := c.Params("id")

	_, err := s.db.ExecContext(context.Background(),
		"DELETE FROM trips WHERE id = @id",
		sql.Named("id", id),
	)

	if err != nil {
		return c.Status(500).JSON(Error("Failed to delete trip"))
	}

	return c.JSON(Success("Trip deleted"))
}

func (s *TripService) GetTripPoints(c *fiber.Ctx) error {
	tripID := c.Params("id")

	rows, err := s.db.QueryContext(context.Background(),
		`SELECT id, trip_id, latitude, longitude, altitude, speed, heading, timestamp 
		 FROM trip_points WHERE trip_id = @tripId ORDER BY timestamp ASC`,
		sql.Named("tripId", tripID),
	)
	if err != nil {
		return c.Status(500).JSON(Error("Database error"))
	}
	defer rows.Close()

	points := []map[string]interface{}{}
	for rows.Next() {
		var id, tripID string
		var lat, lon, altitude, speed, heading sql.NullFloat64
		var timestamp time.Time

		if err := rows.Scan(&id, &tripID, &lat, &lon, &altitude, &speed, &heading, &timestamp); err != nil {
			continue
		}

		point := map[string]interface{}{
			"id": id, "tripId": tripID, "latitude": lat.Float64, "longitude": lon.Float64,
			"timestamp": timestamp,
		}
		if altitude.Valid {
			point["altitude"] = altitude.Float64
		}
		if speed.Valid {
			point["speed"] = speed.Float64
		}
		if heading.Valid {
			point["heading"] = heading.Float64
		}
		points = append(points, point)
	}

	return c.JSON(points)
}

func (s *TripService) AddTripPoint(c *fiber.Ctx) error {
	tripID := c.Params("id")

	type PointInput struct {
		Latitude  float64 `json:"latitude"`
		Longitude float64 `json:"longitude"`
		Altitude  float64 `json:"altitude"`
		Speed     float64 `json:"speed"`
		Heading   float64 `json:"heading"`
	}

	var input PointInput
	if err := c.BodyParser(&input); err != nil {
		return c.Status(400).JSON(Error("Invalid request body"))
	}

	_, err := s.db.ExecContext(context.Background(),
		`INSERT INTO trip_points (trip_id, latitude, longitude, altitude, speed, heading, timestamp)
		 VALUES (@tripId, @lat, @lon, @altitude, @speed, @heading, GETDATE())`,
		sql.Named("tripId", tripID),
		sql.Named("lat", input.Latitude),
		sql.Named("lon", input.Longitude),
		sql.Named("altitude", input.Altitude),
		sql.Named("speed", input.Speed),
		sql.Named("heading", input.Heading),
	)

	if err != nil {
		return c.Status(500).JSON(Error("Failed to add point"))
	}

	return c.Status(201).JSON(Success("Point added"))
}

func (s *TripService) GetSharedSpots(c *fiber.Ctx) error {
	rows, err := s.db.QueryContext(context.Background(),
		`SELECT fs.id, fs.user_id, fs.name, fs.description, fs.latitude, fs.longitude, 
		        fs.rating, fs.fish_caught, fs.is_shared, fs.created_at,
		        u.username, u.display_name
		 FROM fishing_spots fs
		 JOIN users u ON fs.user_id = u.id
		 WHERE fs.is_shared = 1
		 ORDER BY fs.rating DESC, fs.fish_caught DESC`,
	)
	if err != nil {
		return c.Status(500).JSON(Error("Database error"))
	}
	defer rows.Close()

	spots := []map[string]interface{}{}
	for rows.Next() {
		var id, userID, name, description, username, displayName sql.NullString
		var lat, lon sql.NullFloat64
		var rating, fishCaught sql.NullInt64
		var isShared bool
		var createdAt time.Time

		if err := rows.Scan(&id, &userID, &name, &description, &lat, &lon,
			&rating, &fishCaught, &isShared, &createdAt, &username, &displayName); err != nil {
			continue
		}

		spots = append(spots, map[string]interface{}{
			"id":          id.String,
			"userId":      userID.String,
			"name":        name.String,
			"description": description.String,
			"latitude":    lat.Float64,
			"longitude":   lon.Float64,
			"rating":      rating.Int64,
			"fishCaught":  fishCaught.Int64,
			"isShared":    isShared,
			"createdAt":   createdAt,
			"username":    username.String,
			"displayName": displayName.String,
		})
	}

	return c.JSON(spots)
}

func (s *TripService) CreateSpot(c *fiber.Ctx) error {
	type SpotInput struct {
		UserID      string  `json:"userId"`
		Name        string  `json:"name"`
		Description string  `json:"description"`
		Latitude    float64 `json:"latitude"`
		Longitude   float64 `json:"longitude"`
		IsShared    bool    `json:"isShared"`
	}

	var input SpotInput
	if err := c.BodyParser(&input); err != nil {
		return c.Status(400).JSON(Error("Invalid request body"))
	}
	input.UserID = middleware.UserID(c)

	var id string
	err := s.db.QueryRowContext(context.Background(),
		`INSERT INTO fishing_spots (user_id, name, description, latitude, longitude, is_shared)
		 OUTPUT CAST(INSERTED.id AS VARCHAR(36))
		 VALUES (@userId, @name, @description, @lat, @lon, @isShared)`,
		sql.Named("userId", input.UserID),
		sql.Named("name", input.Name),
		sql.Named("description", input.Description),
		sql.Named("lat", input.Latitude),
		sql.Named("lon", input.Longitude),
		sql.Named("isShared", input.IsShared),
	).Scan(&id)

	if err != nil {
		return c.Status(500).JSON(Error("Failed to create spot"))
	}

	return c.Status(201).JSON(fiber.Map{"id": id, "success": true})
}

func (s *TripService) UpdateSpot(c *fiber.Ctx) error {
	id := c.Params("id")

	type UpdateInput struct {
		Name        *string `json:"name"`
		Description *string `json:"description"`
		Rating      *int    `json:"rating"`
		FishCaught  *int    `json:"fishCaught"`
		IsShared    *bool   `json:"isShared"`
	}

	var input UpdateInput
	if err := c.BodyParser(&input); err != nil {
		return c.Status(400).JSON(Error("Invalid request body"))
	}

	_, err := s.db.ExecContext(context.Background(),
		`UPDATE fishing_spots SET 
		 name = ISNULL(@name, name),
		 description = ISNULL(@description, description),
		 rating = ISNULL(@rating, rating),
		 fish_caught = ISNULL(@fishCaught, fish_caught),
		 is_shared = ISNULL(@isShared, is_shared)
		 WHERE id = @id`,
		sql.Named("id", id),
		sql.Named("name", input.Name),
		sql.Named("description", input.Description),
		sql.Named("rating", input.Rating),
		sql.Named("fishCaught", input.FishCaught),
		sql.Named("isShared", input.IsShared),
	)

	if err != nil {
		return c.Status(500).JSON(Error("Failed to update spot"))
	}

	return c.JSON(Success("Spot updated"))
}

func (s *TripService) DeleteSpot(c *fiber.Ctx) error {
	id := c.Params("id")

	_, err := s.db.ExecContext(context.Background(),
		"DELETE FROM fishing_spots WHERE id = @id",
		sql.Named("id", id),
	)

	if err != nil {
		return c.Status(500).JSON(Error("Failed to delete spot"))
	}

	return c.JSON(Success("Spot deleted"))
}
