package services

import (
	"context"
	"database/sql"

	"github.com/gofiber/fiber/v2"
)

type TournamentService struct {
	db *sql.DB
}

func NewTournamentService(db *sql.DB) *TournamentService {
	return &TournamentService{db: db}
}

func (s *TournamentService) GetTournaments(c *fiber.Ctx) error {
	status := c.Query("status")

	query := `SELECT CAST(t.id AS VARCHAR(36)), CAST(t.organizer_id AS VARCHAR(36)), t.name, t.description, t.start_time, t.end_time,
	          t.location_name, t.start_lat, t.start_lon, t.status, t.tournamenttype,
	          t.max_participants, t.created_at,
	          u.username as organizer_username
	          FROM tournaments t
	          JOIN users u ON t.organizer_id = u.id WHERE 1=1`

	var args []interface{}
	if status != "" {
		query += " AND t.status = @status"
		args = append(args, sql.Named("status", status))
	}
	query += " ORDER BY t.start_time DESC"

	rows, err := s.db.QueryContext(context.Background(), query, args...)
	if err != nil {
		return c.Status(500).JSON(Error("Database error"))
	}
	defer rows.Close()

	tournaments := []map[string]interface{}{}
	for rows.Next() {
		var id, name, description, organizerID, locationName, status, tournamentType sql.NullString
		var organizerUsername sql.NullString
		var startTime, endTime, createdAt sql.NullTime
		var startLat, startLon sql.NullFloat64
		var maxParticipants sql.NullInt64

		if err := rows.Scan(&id, &organizerID, &name, &description, &startTime, &endTime,
			&locationName, &startLat, &startLon, &status, &tournamentType,
			&maxParticipants, &createdAt, &organizerUsername); err != nil {
			continue
		}

		t := map[string]interface{}{}
		if id.Valid {
			t["id"] = id.String
		}
		if organizerID.Valid {
			t["organizerId"] = organizerID.String
		}
		if name.Valid {
			t["name"] = name.String
		}
		if description.Valid {
			t["description"] = description.String
		}
		if status.Valid {
			t["status"] = status.String
		}
		if tournamentType.Valid {
			t["type"] = tournamentType.String
		}
		if locationName.Valid {
			t["locationName"] = locationName.String
		}
		if startTime.Valid {
			t["startTime"] = startTime.Time
		}
		if endTime.Valid {
			t["endTime"] = endTime.Time
		}
		if startLat.Valid {
			t["startLat"] = startLat.Float64
			t["startLon"] = startLon.Float64
		}
		if maxParticipants.Valid {
			t["maxParticipants"] = maxParticipants.Int64
		}
		if createdAt.Valid {
			t["createdAt"] = createdAt.Time
		}
		if organizerUsername.Valid {
			t["organizerUsername"] = organizerUsername.String
		}
		tournaments = append(tournaments, t)
	}

	return c.JSON(tournaments)
}

func (s *TournamentService) GetLiveTournaments(c *fiber.Ctx) error {
	rows, err := s.db.QueryContext(context.Background(),
		`SELECT CAST(t.id AS VARCHAR(36)), CAST(t.organizer_id AS VARCHAR(36)), t.name, t.description, t.start_time, t.end_time,
		        t.location_name, t.start_lat, t.start_lon, t.tournamenttype, t.max_participants,
		        u.username as organizer_username,
		        (SELECT COUNT(*) FROM tournament_participants WHERE tournament_id = t.id) as participant_count
		 	FROM tournaments t
		 	JOIN users u ON t.organizer_id = u.id
		 	WHERE t.status = 'LIVE'
		 	ORDER BY t.start_time DESC`,
	)
	if err != nil {
		return c.Status(500).JSON(Error("Database error"))
	}
	defer rows.Close()

	tournaments := []map[string]interface{}{}
	for rows.Next() {
		var id, name, description, organizerID, locationName, tournamentType sql.NullString
		var organizerUsername sql.NullString
		var startTime, endTime sql.NullTime
		var startLat, startLon sql.NullFloat64
		var maxParticipants, participantCount sql.NullInt64

		if err := rows.Scan(&id, &organizerID, &name, &description, &startTime, &endTime,
			&locationName, &startLat, &startLon, &tournamentType,
			&maxParticipants, &organizerUsername, &participantCount); err != nil {
			continue
		}

		t := map[string]interface{}{}
		if id.Valid {
			t["id"] = id.String
		}
		if organizerID.Valid {
			t["organizerId"] = organizerID.String
		}
		if name.Valid {
			t["name"] = name.String
		}
		if description.Valid {
			t["description"] = description.String
		}
		if startTime.Valid {
			t["startTime"] = startTime.Time
		}
		if endTime.Valid {
			t["endTime"] = endTime.Time
		}
		if locationName.Valid {
			t["locationName"] = locationName.String
		}
		if startLat.Valid {
			t["startLat"] = startLat.Float64
		}
		if startLon.Valid {
			t["startLon"] = startLon.Float64
		}
		if tournamentType.Valid {
			t["type"] = tournamentType.String
		}
		if maxParticipants.Valid {
			t["maxParticipants"] = maxParticipants.Int64
		}
		if organizerUsername.Valid {
			t["organizerUsername"] = organizerUsername.String
		}
		if participantCount.Valid {
			t["participantsCount"] = participantCount.Int64
		}
		tournaments = append(tournaments, t)
	}

	return c.JSON(tournaments)
}

func (s *TournamentService) GetTournament(c *fiber.Ctx) error {
	id := c.Params("id")

	var tournament struct {
		ID, Name, Description, OrganizerID, LocationName, Status, Type string
		StartTime, EndTime                                             sql.NullTime
		StartLat, StartLon                                             sql.NullFloat64
		MaxParticipants                                                sql.NullInt64
	}

	err := s.db.QueryRowContext(context.Background(),
		`SELECT id, name, description, organizer_id, start_time, end_time,
		        location_name, start_lat, start_lon, status, tournamenttype, max_participants
		 FROM tournaments WHERE id = @id`,
		sql.Named("id", id),
	).Scan(&tournament.ID, &tournament.Name, &tournament.Description,
		&tournament.OrganizerID, &tournament.StartTime, &tournament.EndTime,
		&tournament.LocationName, &tournament.StartLat, &tournament.StartLon,
		&tournament.Status, &tournament.Type, &tournament.MaxParticipants)

	if err == sql.ErrNoRows {
		return c.Status(404).JSON(Error("Tournament not found"))
	}
	if err != nil {
		return c.Status(500).JSON(Error("Database error"))
	}

	return c.JSON(fiber.Map{
		"id":              tournament.ID,
		"name":            tournament.Name,
		"description":     tournament.Description,
		"organizerId":     tournament.OrganizerID,
		"startTime":       tournament.StartTime.Time,
		"endTime":         tournament.EndTime.Time,
		"locationName":    tournament.LocationName,
		"startLat":        tournament.StartLat.Float64,
		"startLon":        tournament.StartLon.Float64,
		"status":          tournament.Status,
		"type":            tournament.Type,
		"maxParticipants": tournament.MaxParticipants.Int64,
	})
}

func (s *TournamentService) CreateTournament(c *fiber.Ctx) error {
	type Input struct {
		OrganizerID     string  `json:"organizerId"`
		Name            string  `json:"name"`
		Description     string  `json:"description"`
		StartTime       string  `json:"startTime"`
		EndTime         string  `json:"endTime"`
		LocationName    string  `json:"locationName"`
		StartLat        float64 `json:"startLat"`
		StartLon        float64 `json:"startLon"`
		Type            string  `json:"type"`
		MaxParticipants int     `json:"maxParticipants"`
	}

	var input Input
	if err := c.BodyParser(&input); err != nil {
		return c.Status(400).JSON(Error("Invalid request"))
	}

	var id string
	err := s.db.QueryRowContext(context.Background(),
		`INSERT INTO tournaments (name, description, organizer_id, start_time, end_time, 
		 location_name, start_lat, start_lon, tournamenttype, max_participants, status)
		 OUTPUT CAST(INSERTED.id AS VARCHAR(36))
		 VALUES (@name, @desc, @orgId, @start, @end, @loc, @lat, @lon, @type, @max, 'planned')`,
		sql.Named("name", input.Name),
		sql.Named("desc", input.Description),
		sql.Named("orgId", input.OrganizerID),
		sql.Named("start", input.StartTime),
		sql.Named("end", input.EndTime),
		sql.Named("loc", input.LocationName),
		sql.Named("lat", input.StartLat),
		sql.Named("lon", input.StartLon),
		sql.Named("type", input.Type),
		sql.Named("max", input.MaxParticipants),
	).Scan(&id)

	if err != nil {
		return c.Status(500).JSON(Error("Failed to create tournament"))
	}

	return c.Status(201).JSON(fiber.Map{"id": id, "success": true})
}

func (s *TournamentService) UpdateTournament(c *fiber.Ctx) error {
	id := c.Params("id")

	type Input struct {
		Name        *string `json:"name"`
		Description *string `json:"description"`
		Status      *string `json:"status"`
		StartTime   *string `json:"startTime"`
		EndTime     *string `json:"endTime"`
	}

	var input Input
	if err := c.BodyParser(&input); err != nil {
		return c.Status(400).JSON(Error("Invalid request"))
	}

	_, err := s.db.ExecContext(context.Background(),
		`UPDATE tournaments SET 
		 name = ISNULL(@name, name),
		 description = ISNULL(@desc, description),
		 status = ISNULL(@status, status),
		 start_time = ISNULL(@start, start_time),
		 end_time = ISNULL(@end, end_time)
		 WHERE id = @id`,
		sql.Named("id", id),
		sql.Named("name", input.Name),
		sql.Named("desc", input.Description),
		sql.Named("status", input.Status),
		sql.Named("start", input.StartTime),
		sql.Named("end", input.EndTime),
	)

	if err != nil {
		return c.Status(500).JSON(Error("Failed to update"))
	}

	return c.JSON(Success("Tournament updated"))
}

func (s *TournamentService) GetParticipants(c *fiber.Ctx) error {
	tournamentID := c.Params("id")

	rows, err := s.db.QueryContext(context.Background(),
		`SELECT CAST(p.id AS VARCHAR(36)), CAST(p.user_id AS VARCHAR(36)), CAST(p.kayak_id AS VARCHAR(36)),
		        p.status, p.current_lat, p.current_lon,
		        p.score, p.catches_count, p.last_update_time, p.joined_at,
		        u.username, u.display_name, k.name as kayak_name
		 FROM tournament_participants p
		 JOIN users u ON p.user_id = u.id
		 LEFT JOIN kayaks k ON p.kayak_id = k.id
		 WHERE p.tournament_id = @tournamentId
		 ORDER BY p.score DESC, p.catches_count DESC`,
		sql.Named("tournamentId", tournamentID),
	)
	if err != nil {
		return c.Status(500).JSON(Error("Database error"))
	}
	defer rows.Close()

	participants := []map[string]interface{}{}
	for rows.Next() {
		var id, userID, kayakID, status, username, displayName, kayakName sql.NullString
		var currentLat, currentLon, score sql.NullFloat64
		var catchesCount sql.NullInt64
		var lastUpdateTime, joinedAt sql.NullTime

		if err := rows.Scan(&id, &userID, &kayakID, &status, &currentLat, &currentLon,
			&score, &catchesCount, &lastUpdateTime, &joinedAt,
			&username, &displayName, &kayakName); err != nil {
			continue
		}

		p := map[string]interface{}{}
		if id.Valid {
			p["id"] = id.String
		}
		if userID.Valid {
			p["userId"] = userID.String
		}
		if kayakID.Valid {
			p["kayakId"] = kayakID.String
		}
		if status.Valid {
			p["status"] = status.String
		}
		if currentLat.Valid {
			p["currentLat"] = currentLat.Float64
		}
		if currentLon.Valid {
			p["currentLon"] = currentLon.Float64
		}
		if score.Valid {
			p["score"] = score.Float64
		}
		if catchesCount.Valid {
			p["catchesCount"] = catchesCount.Int64
		}
		if lastUpdateTime.Valid {
			p["lastUpdateTime"] = lastUpdateTime.Time
		}
		if joinedAt.Valid {
			p["joinedAt"] = joinedAt.Time
		}
		if username.Valid {
			p["username"] = username.String
		}
		if displayName.Valid {
			p["displayName"] = displayName.String
		}
		if kayakName.Valid {
			p["kayakName"] = kayakName.String
		}
		participants = append(participants, p)
	}

	return c.JSON(participants)
}

func (s *TournamentService) GetRanking(c *fiber.Ctx) error {
	tournamentID := c.Params("id")

	rows, err := s.db.QueryContext(context.Background(),
		`SELECT CAST(p.id AS VARCHAR(36)), CAST(p.user_id AS VARCHAR(36)), p.score, p.catches_count,
		        p.current_lat, p.current_lon,
		        u.username, u.display_name, k.name as kayak_name,
		        ROW_NUMBER() OVER (ORDER BY p.score DESC, p.catches_count DESC) as rank
		 FROM tournament_participants p
		 JOIN users u ON p.user_id = u.id
		 LEFT JOIN kayaks k ON p.kayak_id = k.id
		 WHERE p.tournament_id = @tournamentId AND p.status = 'active'`,
		sql.Named("tournamentId", tournamentID),
	)
	if err != nil {
		return c.Status(500).JSON(Error("Database error"))
	}
	defer rows.Close()

	ranking := []map[string]interface{}{}
	for rows.Next() {
		var id, userID, username, displayName, kayakName sql.NullString
		var score, currentLat, currentLon sql.NullFloat64
		var catchesCount, rank sql.NullInt64

		if err := rows.Scan(&id, &userID, &score, &catchesCount, &currentLat, &currentLon,
			&username, &displayName, &kayakName, &rank); err != nil {
			continue
		}

		r := map[string]interface{}{}
		if id.Valid {
			r["id"] = id.String
		}
		if userID.Valid {
			r["userId"] = userID.String
		}
		if score.Valid {
			r["score"] = score.Float64
		}
		if catchesCount.Valid {
			r["catchesCount"] = catchesCount.Int64
		}
		if currentLat.Valid {
			r["currentLat"] = currentLat.Float64
		}
		if currentLon.Valid {
			r["currentLon"] = currentLon.Float64
		}
		if username.Valid {
			r["username"] = username.String
		}
		if displayName.Valid {
			r["displayName"] = displayName.String
		}
		if kayakName.Valid {
			r["kayakName"] = kayakName.String
		}
		if rank.Valid {
			r["rank"] = rank.Int64
		}
		ranking = append(ranking, r)
	}

	return c.JSON(ranking)
}

func (s *TournamentService) JoinTournament(c *fiber.Ctx) error {
	tournamentID := c.Params("id")

	type Input struct {
		UserID  string `json:"userId"`
		KayakID string `json:"kayakId"`
	}

	var input Input
	if err := c.BodyParser(&input); err != nil {
		return c.Status(400).JSON(Error("Invalid request"))
	}

	_, err := s.db.ExecContext(context.Background(),
		`INSERT INTO tournament_participants (tournament_id, user_id, kayak_id, status, joined_at)
		 VALUES (@tournamentId, @userId, @kayakId, 'registered', GETDATE())`,
		sql.Named("tournamentId", tournamentID),
		sql.Named("userId", input.UserID),
		sql.Named("kayakId", input.KayakID),
	)

	if err != nil {
		return c.Status(500).JSON(Error("Failed to join"))
	}

	return c.Status(201).JSON(Success("Joined tournament"))
}

func (s *TournamentService) LeaveTournament(c *fiber.Ctx) error {
	tournamentID := c.Params("id")

	type Input struct {
		UserID string `json:"userId"`
	}

	var input Input
	if err := c.BodyParser(&input); err != nil {
		return c.Status(400).JSON(Error("Invalid request"))
	}

	_, err := s.db.ExecContext(context.Background(),
		`UPDATE tournament_participants SET status = 'left'
		 WHERE tournament_id = @tournamentId AND user_id = @userId`,
		sql.Named("tournamentId", tournamentID),
		sql.Named("userId", input.UserID),
	)

	if err != nil {
		return c.Status(500).JSON(Error("Failed to leave"))
	}

	return c.JSON(Success("Left tournament"))
}

func (s *TournamentService) RecordCatch(c *fiber.Ctx) error {
	tournamentID := c.Params("id")

	type Input struct {
		ParticipantID string  `json:"participantId"`
		Species       string  `json:"species"`
		Weight        float64 `json:"weight"`
		Latitude      float64 `json:"latitude"`
		Longitude     float64 `json:"longitude"`
		PhotoURL      string  `json:"photoUrl"`
	}

	var input Input
	if err := c.BodyParser(&input); err != nil {
		return c.Status(400).JSON(Error("Invalid request"))
	}

	tx, err := s.db.BeginTx(context.Background(), nil)
	if err != nil {
		return c.Status(500).JSON(Error("Transaction error"))
	}

	var catchID string
	err = tx.QueryRowContext(context.Background(),
		`INSERT INTO tournament_catches (tournament_id, participant_id, species, weight, latitude, longitude, photo_url)
		 OUTPUT CAST(INSERTED.id AS VARCHAR(36))
		 VALUES (@tournamentId, @participantId, @species, @weight, @lat, @lon, @photoUrl)`,
		sql.Named("tournamentId", tournamentID),
		sql.Named("participantId", input.ParticipantID),
		sql.Named("species", input.Species),
		sql.Named("weight", input.Weight),
		sql.Named("lat", input.Latitude),
		sql.Named("lon", input.Longitude),
		sql.Named("photoUrl", input.PhotoURL),
	).Scan(&catchID)

	if err != nil {
		tx.Rollback()
		return c.Status(500).JSON(Error("Failed to record catch"))
	}

	_, err = tx.ExecContext(context.Background(),
		`UPDATE tournament_participants SET 
		 catches_count = catches_count + 1,
		 score = score + @weight
		 WHERE id = @participantId`,
		sql.Named("participantId", input.ParticipantID),
		sql.Named("weight", input.Weight),
	)

	if err != nil {
		tx.Rollback()
		return c.Status(500).JSON(Error("Failed to update score"))
	}

	tx.Commit()
	return c.Status(201).JSON(fiber.Map{"id": catchID, "success": true})
}

func (s *TournamentService) GetCatches(c *fiber.Ctx) error {
	tournamentID := c.Params("id")

	rows, err := s.db.QueryContext(context.Background(),
		`SELECT CAST(c.id AS VARCHAR(36)), c.species, c.weight, c.latitude, c.longitude, c.photo_url, 
		        c.verified, c.verified_by, c.created_at,
		        u.username, CAST(p.id AS VARCHAR(36)) as participant_id
		 FROM tournament_catches c
		 JOIN tournament_participants p ON c.participant_id = p.id
		 JOIN users u ON p.user_id = u.id
		 WHERE c.tournament_id = @tournamentId
		 ORDER BY c.weight DESC`,
		sql.Named("tournamentId", tournamentID),
	)
	if err != nil {
		return c.Status(500).JSON(Error("Database error"))
	}
	defer rows.Close()

	catches := []map[string]interface{}{}
	for rows.Next() {
		var id, species, photoURL, verifiedBy, username, participantID sql.NullString
		var weight, latitude, longitude sql.NullFloat64
		var verified bool
		var createdAt sql.NullTime

		if err := rows.Scan(&id, &species, &weight, &latitude, &longitude, &photoURL,
			&verified, &verifiedBy, &createdAt, &username, &participantID); err != nil {
			continue
		}

		item := map[string]interface{}{
			"verified": verified,
		}
		if id.Valid {
			item["id"] = id.String
		}
		if species.Valid {
			item["species"] = species.String
		}
		if weight.Valid {
			item["weight"] = weight.Float64
		}
		if latitude.Valid {
			item["latitude"] = latitude.Float64
		}
		if longitude.Valid {
			item["longitude"] = longitude.Float64
		}
		if photoURL.Valid {
			item["photoUrl"] = photoURL.String
		}
		if verifiedBy.Valid {
			item["verifiedBy"] = verifiedBy.String
		}
		if createdAt.Valid {
			item["createdAt"] = createdAt.Time
		}
		if username.Valid {
			item["username"] = username.String
		}
		if participantID.Valid {
			item["participantId"] = participantID.String
		}
		catches = append(catches, item)
	}

	return c.JSON(catches)
}
