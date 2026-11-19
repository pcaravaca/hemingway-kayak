package services

import (
	"context"
	"database/sql"
	"strconv"
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

// tripQuery arma el SELECT de salidas con los filtros de querystring. Si
// forcedUserID no está vacío, se fuerza el dueño (para la lista de "mis salidas").
func (s *TripService) tripQuery(c *fiber.Ctx, forcedUserID string) (string, []interface{}) {
	query := `SELECT CAST(id AS VARCHAR(36)), CAST(user_id AS VARCHAR(36)), title, description, triptype, status, start_time, end_time,
	          start_lat, start_lon, end_lat, end_lon, distance, created_at
	          FROM trips WHERE 1=1`
	var args []interface{}

	owner := forcedUserID
	if owner == "" {
		owner = c.Query("userId")
	}
	if owner != "" {
		query += " AND user_id = $" + strconv.Itoa(len(args)+1)
		args = append(args, owner)
	}

	if tripType := c.Query("type"); tripType != "" {
		query += " AND triptype = $" + strconv.Itoa(len(args)+1)
		args = append(args, tripType)
	}
	if status := c.Query("status"); status != "" {
		query += " AND status = $" + strconv.Itoa(len(args)+1)
		args = append(args, status)
	}

	return query + " ORDER BY created_at DESC", args
}

// GetTrips es el feed público de salidas. No expone coordenadas: son datos de
// ubicación y el detalle de una salida ya es sólo del dueño.
func (s *TripService) GetTrips(c *fiber.Ctx) error {
	query, args := s.tripQuery(c, "")
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

// GetMyTrips devuelve sólo las salidas del usuario autenticado: se usa para los
// detalles (que son del dueño) y para el historial.
func (s *TripService) GetMyTrips(c *fiber.Ctx) error {
	query, args := s.tripQuery(c, middleware.UserID(c))
	rows, err := s.db.QueryContext(context.Background(), query, args...)
	if err != nil {
		return c.Status(500).JSON(Error("Database error"))
	}
	defer rows.Close()

	return c.JSON(scanTrips(rows))
}

// scanTrips convierte las filas de una salida a JSON, con coordenadas (es para el
// dueño; el feed público no las expone).
func scanTrips(rows *sql.Rows) []map[string]interface{} {
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
		if endLat.Valid {
			trip["endLat"] = endLat.Float64
		}
		if endLon.Valid {
			trip["endLon"] = endLon.Float64
		}
		if distance.Valid {
			trip["distance"] = distance.Float64
		}
		if createdAt.Valid {
			trip["createdAt"] = createdAt.Time
		}
		trips = append(trips, trip)
	}
	return trips
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
		 FROM trips WHERE id = $1`,
		id,
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
		 VALUES ($1, $2, $3, $4, $5, $6, 'planned', now())
		 RETURNING id::text`,
		input.UserID,
		input.Title,
		input.Description,
		input.Type,
		input.StartLat,
		input.StartLon,
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
		title = COALESCE($1, title),
		description = COALESCE($2, description),
		status = COALESCE($3, status),
		end_lat = COALESCE($4, end_lat),
		end_lon = COALESCE($5, end_lon),
		distance = COALESCE($6, distance),
		end_time = CASE WHEN $3 IN ('completed', 'cancelled') THEN now() ELSE end_time END
		WHERE id = $7`

	_, err := s.db.ExecContext(context.Background(), query,
		input.Title,
		input.Description,
		input.Status,
		input.EndLat,
		input.EndLon,
		input.Distance,
		id,
	)

	if err != nil {
		return c.Status(500).JSON(Error("Failed to update trip"))
	}

	return c.JSON(Success("Trip updated"))
}

func (s *TripService) DeleteTrip(c *fiber.Ctx) error {
	id := c.Params("id")

	_, err := s.db.ExecContext(context.Background(),
		"DELETE FROM trips WHERE id = $1",
		id,
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
		 FROM trip_points WHERE trip_id = $1 ORDER BY timestamp ASC`,
		tripID,
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
		 VALUES ($1, $2, $3, $4, $5, $6, now())`,
		tripID,
		input.Latitude,
		input.Longitude,
		input.Altitude,
		input.Speed,
		input.Heading,
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
		 WHERE fs.is_shared = TRUE
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
		 VALUES ($1, $2, $3, $4, $5, $6)
		 RETURNING id::text`,
		input.UserID,
		input.Name,
		input.Description,
		input.Latitude,
		input.Longitude,
		input.IsShared,
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
		 name = COALESCE($1, name),
		 description = COALESCE($2, description),
		 rating = COALESCE($3, rating),
		 fish_caught = COALESCE($4, fish_caught),
		 is_shared = COALESCE($5, is_shared)
		 WHERE id = $6`,
		input.Name,
		input.Description,
		input.Rating,
		input.FishCaught,
		input.IsShared,
		id,
	)

	if err != nil {
		return c.Status(500).JSON(Error("Failed to update spot"))
	}

	return c.JSON(Success("Spot updated"))
}

func (s *TripService) DeleteSpot(c *fiber.Ctx) error {
	id := c.Params("id")

	_, err := s.db.ExecContext(context.Background(),
		"DELETE FROM fishing_spots WHERE id = $1",
		id,
	)

	if err != nil {
		return c.Status(500).JSON(Error("Failed to delete spot"))
	}

	return c.JSON(Success("Spot deleted"))
}
