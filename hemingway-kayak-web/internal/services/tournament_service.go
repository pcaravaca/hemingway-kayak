package services

import (
	"context"
	"database/sql"
	"strconv"

	"github.com/gofiber/fiber/v2"

	"hemingway-kayak-web/internal/middleware"
)

type TournamentService struct {
	db  *sql.DB
	hub wsBroadcaster
}

func NewTournamentService(db *sql.DB) *TournamentService {
	return &TournamentService{db: db}
}

// SetBroadcaster engancha el hub de WebSocket para el ranking en vivo.
func (s *TournamentService) SetBroadcaster(hub wsBroadcaster) {
	s.hub = hub
}

// rankingQuery es la única fuente de verdad del ranking: la usa GetRanking y
// también el broadcast en vivo.
const rankingQuery = `SELECT CAST(p.id AS VARCHAR(36)), CAST(p.user_id AS VARCHAR(36)), p.score, p.catches_count,
	        p.current_lat, p.current_lon,
	        u.username, u.display_name, k.name as kayak_name,
	        ROW_NUMBER() OVER (ORDER BY p.score DESC, p.catches_count DESC) as rank
	 FROM tournament_participants p
	 JOIN users u ON p.user_id = u.id
	 LEFT JOIN kayaks k ON p.kayak_id = k.id
	 WHERE p.tournament_id = $1 AND p.status = 'active'`

func (s *TournamentService) rankingFor(tournamentID string) []map[string]interface{} {
	rows, err := s.db.QueryContext(context.Background(), rankingQuery, tournamentID)
	if err != nil {
		return nil
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

	return ranking
}

// broadcastRanking empuja el ranking a la sala del torneo.
func (s *TournamentService) broadcastRanking(tournamentID string) {
	if s.hub == nil {
		return
	}
	s.hub.BroadcastTo(tournamentID, "tournament", fiber.Map{
		"tournamentId": tournamentID,
		"rankings":     s.rankingFor(tournamentID),
	})
}

func (s *TournamentService) GetTournaments(c *fiber.Ctx) error {
	status := c.Query("status")

	query := `SELECT CAST(t.id AS VARCHAR(36)), CAST(t.organizer_id AS VARCHAR(36)), t.name, t.description, t.start_time, t.end_time,
	          t.location_name, t.start_lat, t.start_lon, UPPER(t.status) AS status, t.tournamenttype,
	          t.max_participants, t.created_at,
	          u.username as organizer_username,
	          (SELECT COUNT(*) FROM tournament_participants WHERE tournament_id = t.id) as participant_count
	          FROM tournaments t
	          JOIN users u ON t.organizer_id = u.id WHERE 1=1`

	var args []interface{}
	if status != "" {
		query += " AND UPPER(t.status) = UPPER($" + strconv.Itoa(len(args)+1) + ")"
		args = append(args, status)
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
		var maxParticipants, participantCount sql.NullInt64

		if err := rows.Scan(&id, &organizerID, &name, &description, &startTime, &endTime,
			&locationName, &startLat, &startLon, &status, &tournamentType,
			&maxParticipants, &createdAt, &organizerUsername, &participantCount); err != nil {
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
		if participantCount.Valid {
			t["participantsCount"] = participantCount.Int64
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
		ID, Name                   string
		Description, OrganizerID   sql.NullString
		LocationName, Status, Type sql.NullString
		StartTime, EndTime         sql.NullTime
		StartLat, StartLon         sql.NullFloat64
		MaxParticipants            sql.NullInt64
	}

	err := s.db.QueryRowContext(context.Background(),
		`SELECT id, name, description, organizer_id, start_time, end_time,
		        location_name, start_lat, start_lon, status, tournamenttype, max_participants
		 FROM tournaments WHERE id = $1`,
		id,
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
		"description":     tournament.Description.String,
		"organizerId":     tournament.OrganizerID.String,
		"startTime":       tournament.StartTime.Time,
		"endTime":         tournament.EndTime.Time,
		"locationName":    tournament.LocationName.String,
		"startLat":        tournament.StartLat.Float64,
		"startLon":        tournament.StartLon.Float64,
		"status":          tournament.Status.String,
		"type":            tournament.Type.String,
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
	// El organizador es siempre el usuario autenticado.
	input.OrganizerID = middleware.UserID(c)

	var id string
	err := s.db.QueryRowContext(context.Background(),
		`INSERT INTO tournaments (name, description, organizer_id, start_time, end_time, 
		 location_name, start_lat, start_lon, tournamenttype, max_participants, status)
		 VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, 'PLANNED')
		 RETURNING id::text`,
		input.Name,
		input.Description,
		input.OrganizerID,
		input.StartTime,
		input.EndTime,
		input.LocationName,
		input.StartLat,
		input.StartLon,
		input.Type,
		input.MaxParticipants,
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
		 name = COALESCE($1, name),
		 description = COALESCE($2, description),
		 status = COALESCE($3, status),
		 start_time = COALESCE($4, start_time),
		 end_time = COALESCE($5, end_time)
		 WHERE id = $6`,
		input.Name,
		input.Description,
		input.Status,
		input.StartTime,
		input.EndTime,
		id,
	)

	if err != nil {
		return c.Status(500).JSON(Error("Failed to update"))
	}

	// Si el torneo pasa a LIVE (o termina), quien esté mirando la tabla tiene que
	// enterarse sin recargar.
	s.broadcastRanking(id)
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
		 WHERE p.tournament_id = $1
		 ORDER BY p.score DESC, p.catches_count DESC`,
		tournamentID,
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
	return c.JSON(s.rankingFor(c.Params("id")))
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
	input.UserID = middleware.UserID(c)

	// kayak_id es opcional: un string vacío no es un UUID válido y Postgres
	// rechaza el INSERT con un 500 pelado. NullString deja la columna en NULL.
	var kayakID sql.NullString
	if input.KayakID != "" {
		kayakID = sql.NullString{String: input.KayakID, Valid: true}
	}

	// 'active' y no 'registered': el ranking filtra por status='active', así que
	// con 'registered' el recién unido nunca aparecía en la tabla del torneo.
	// Idempotente: volver a unirse reactiva la fila en lugar de duplicarla.
	res, err := s.db.ExecContext(context.Background(),
		`UPDATE tournament_participants
		    SET status = 'active', kayak_id = $3, joined_at = now()
		  WHERE tournament_id = $1 AND user_id = $2`,
		tournamentID,
		input.UserID,
		kayakID,
	)
	if err != nil {
		return c.Status(500).JSON(Error("Failed to join"))
	}

	if n, err := res.RowsAffected(); err == nil && n == 0 {
		_, err = s.db.ExecContext(context.Background(),
			`INSERT INTO tournament_participants (tournament_id, user_id, kayak_id, status, joined_at)
			 VALUES ($1, $2, $3, 'active', now())`,
			tournamentID,
			input.UserID,
			kayakID,
		)
		if err != nil {
			return c.Status(500).JSON(Error("Failed to join"))
		}
	}

	s.broadcastRanking(tournamentID)
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
	input.UserID = middleware.UserID(c)

	_, err := s.db.ExecContext(context.Background(),
		`UPDATE tournament_participants SET status = 'left'
		 WHERE tournament_id = $1 AND user_id = $2`,
		tournamentID,
		input.UserID,
	)

	if err != nil {
		return c.Status(500).JSON(Error("Failed to leave"))
	}

	s.broadcastRanking(tournamentID)
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
		 VALUES ($1, $2, $3, $4, $5, $6, $7)
		 RETURNING id::text`,
		tournamentID,
		input.ParticipantID,
		input.Species,
		input.Weight,
		input.Latitude,
		input.Longitude,
		input.PhotoURL,
	).Scan(&catchID)

	if err != nil {
		tx.Rollback()
		return c.Status(500).JSON(Error("Failed to record catch"))
	}

	tx.Commit()
	// La captura entra sin verificar: NO suma al puntaje todavía. Lo hace
	// SetCatchVerified cuando el organizador la verifica (y lo quita al desverificar).
	// Tampoco hay ranking que empujar acá.
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
		 WHERE c.tournament_id = $1
		 ORDER BY c.weight DESC`,
		tournamentID,
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

// SetCatchVerified marca (o desmarca) una captura como verificada. Sólo el
// organizador del torneo puede hacerlo: el middleware RequireResourceOwnership
// compara tournaments.organizer_id con el JWT y responde 403 si no coincide.
//
// El puntaje cuenta SÓLO las capturas verificadas: verificar suma el peso y una
// captura al participante; desverificar lo resta (con GREATEST para no ir negativo).
// Por eso la transición es la que dispara el broadcast del ranking.
func (s *TournamentService) SetCatchVerified(c *fiber.Ctx) error {
	tournamentID := c.Params("id")
	catchID := c.Params("catchId")
	if catchID == "" {
		return c.Status(400).JSON(Error("missing catch id"))
	}

	// Sin cuerpo (o cuerpo sin `verified`) se entiende como "verificar".
	verified := true
	var body struct {
		Verified *bool `json:"verified"`
	}
	if err := c.BodyParser(&body); err == nil && body.Verified != nil {
		verified = *body.Verified
	}

	ctx := context.Background()
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return c.Status(500).JSON(Error("Transaction error"))
	}
	defer tx.Rollback()

	var before bool
	var participantID string
	var weight sql.NullFloat64
	if err := tx.QueryRowContext(ctx,
		`SELECT verified, CAST(participant_id AS VARCHAR(36)), weight
		   FROM tournament_catches
		  WHERE id = $1 AND tournament_id = $2`,
		catchID, tournamentID,
	).Scan(&before, &participantID, &weight); err != nil {
		if err == sql.ErrNoRows {
			return c.Status(404).JSON(Error("Catch not found in this tournament"))
		}
		return c.Status(500).JSON(Error("Database error"))
	}

	var verifiedBy sql.NullString
	if verified {
		if userID, _ := c.Locals("userId").(string); userID != "" {
			var username string
			if err := tx.QueryRowContext(ctx,
				`SELECT username FROM users WHERE id = $1`, userID,
			).Scan(&username); err == nil {
				verifiedBy = sql.NullString{String: username, Valid: true}
			}
		}
	}

	if _, err := tx.ExecContext(ctx,
		`UPDATE tournament_catches
		    SET verified = $1, verified_by = $2
		  WHERE id = $3 AND tournament_id = $4`,
		verified, verifiedBy, catchID, tournamentID,
	); err != nil {
		return c.Status(500).JSON(Error("Database error"))
	}

	if before != verified {
		w := weight.Float64
		if verified {
			_, err = tx.ExecContext(ctx,
				`UPDATE tournament_participants
				    SET score = score + $1, catches_count = catches_count + 1
				  WHERE id = $2`, w, participantID)
		} else {
			_, err = tx.ExecContext(ctx,
				`UPDATE tournament_participants
				    SET score = GREATEST(score - $1, 0),
				        catches_count = GREATEST(catches_count - 1, 0)
				  WHERE id = $2`, w, participantID)
		}
		if err != nil {
			return c.Status(500).JSON(Error("Failed to update score"))
		}
	}

	if err := tx.Commit(); err != nil {
		return c.Status(500).JSON(Error("Transaction error"))
	}

	if before != verified {
		s.broadcastRanking(tournamentID)
	}

	return c.JSON(fiber.Map{"success": true, "verified": verified})
}
