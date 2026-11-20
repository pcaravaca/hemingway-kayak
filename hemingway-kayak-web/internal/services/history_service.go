package services

import (
	"context"
	"database/sql"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

type HistoryService struct {
	db *sql.DB
}

func NewHistoryService(db *sql.DB) *HistoryService {
	return &HistoryService{db: db}
}

func (s *HistoryService) GetUserHistory(c *fiber.Ctx) error {
	userID := c.Params("userId")
	historyType := c.Query("type")

	query := `SELECT id, user_id, historytype, reference_id, title, description, data, occurred_at, created_at
	          FROM history_entries WHERE user_id = $1`

	args := []interface{}{userID}
	if historyType != "" {
		query += " AND historytype = $" + strconv.Itoa(len(args)+1)
		args = append(args, historyType)
	}
	query += " ORDER BY occurred_at DESC"

	rows, err := s.db.QueryContext(context.Background(), query, args...)
	if err != nil {
		return c.Status(500).JSON(Error("Database error"))
	}
	defer rows.Close()

	entries := []map[string]interface{}{}
	for rows.Next() {
		var id, userID, historyType, title, description, data sql.NullString
		var referenceID sql.NullString
		var occurredAt, createdAt sql.NullTime

		if err := rows.Scan(&id, &userID, &historyType, &referenceID, &title, &description, &data, &occurredAt, &createdAt); err != nil {
			continue
		}

		entry := map[string]interface{}{}
		if id.Valid {
			entry["id"] = id.String
		}
		if userID.Valid {
			entry["userId"] = userID.String
		}
		if historyType.Valid {
			entry["type"] = historyType.String
		}
		if referenceID.Valid {
			entry["referenceId"] = referenceID.String
		}
		if title.Valid {
			entry["title"] = title.String
		}
		if description.Valid {
			entry["description"] = description.String
		}
		if data.Valid {
			entry["data"] = data.String
		}
		if occurredAt.Valid {
			entry["occurredAt"] = occurredAt.Time
		}
		if createdAt.Valid {
			entry["createdAt"] = createdAt.Time
		}
		entries = append(entries, entry)
	}

	return c.JSON(entries)
}

func (s *HistoryService) GetTripHistory(c *fiber.Ctx) error {
	userID := c.Params("userId")

	rows, err := s.db.QueryContext(context.Background(),
		`SELECT t.id, t.user_id, t.title, t.triptype, UPPER(t.status), t.start_time, t.end_time, 
	        t.distance, t.created_at, t.start_lat, t.start_lon,
	        (SELECT COUNT(*) FROM trip_points WHERE trip_id = t.id) as point_count
	 	FROM trips t
	 	WHERE t.user_id = $1
	 	ORDER BY t.created_at DESC`,
		userID,
	)
	if err != nil {
		return c.Status(500).JSON(Error("Database error"))
	}
	defer rows.Close()

	trips := []map[string]interface{}{}
	for rows.Next() {
		var id, userID, title, tripType, status sql.NullString
		var startTime, endTime, createdAt sql.NullTime
		var distance, startLat, startLon sql.NullFloat64
		var pointCount int

		if err := rows.Scan(&id, &userID, &title, &tripType, &status, &startTime, &endTime, &distance, &createdAt, &startLat, &startLon, &pointCount); err != nil {
			continue
		}

		trip := map[string]interface{}{}
		if id.Valid {
			trip["id"] = id.String
		}
		if title.Valid {
			trip["title"] = title.String
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
		if startLat.Valid {
			trip["startLat"] = startLat.Float64
		}
		if startLon.Valid {
			trip["startLon"] = startLon.Float64
		}
		trip["pointCount"] = pointCount
		trips = append(trips, trip)
	}

	return c.JSON(trips)
}

func (s *HistoryService) GetTournamentHistory(c *fiber.Ctx) error {
	userID := c.Params("userId")

	rows, err := s.db.QueryContext(context.Background(),
		`SELECT tp.tournament_id, t.name, t.tournamenttype, t.status, t.start_time, t.end_time,
	        tp.score, tp.catches_count, tp.joined_at
	 	FROM tournament_participants tp
	 	JOIN tournaments t ON tp.tournament_id = t.id
	 	WHERE tp.user_id = $1
	 	ORDER BY tp.joined_at DESC`,
		userID,
	)
	if err != nil {
		return c.Status(500).JSON(Error("Database error"))
	}
	defer rows.Close()

	history := []map[string]interface{}{}
	for rows.Next() {
		var tournamentID, name, tournamentType, status sql.NullString
		var startTime, endTime, joinedAt sql.NullTime
		var score sql.NullFloat64
		var catchesCount sql.NullInt64

		if err := rows.Scan(&tournamentID, &name, &tournamentType, &status, &startTime, &endTime, &score, &catchesCount, &joinedAt); err != nil {
			continue
		}

		entry := map[string]interface{}{}
		if tournamentID.Valid {
			entry["tournamentId"] = tournamentID.String
		}
		if name.Valid {
			entry["name"] = name.String
		}
		if tournamentType.Valid {
			entry["type"] = tournamentType.String
		}
		if status.Valid {
			entry["status"] = status.String
		}
		if startTime.Valid {
			entry["startTime"] = startTime.Time
		}
		if score.Valid {
			entry["score"] = score.Float64
		}
		if catchesCount.Valid {
			entry["catchesCount"] = catchesCount.Int64
		}
		history = append(history, entry)
	}

	return c.JSON(history)
}
