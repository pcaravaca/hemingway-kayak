package services

import (
	"context"
	"database/sql"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"

	"hemingway-kayak-web/internal/middleware"
	"hemingway-kayak-web/pkg/auth"
	"hemingway-kayak-web/pkg/config"
)

type UserService struct {
	db        *sql.DB
	jwtSecret []byte
	revoked   *auth.Revocation
	ttl       time.Duration
}

func NewUserService(db *sql.DB) *UserService {
	return &UserService{
		db:        db,
		jwtSecret: []byte(config.Load().JWTSecret),
		revoked:   auth.NewRevocation(),
		ttl:       24 * time.Hour * 7,
	}
}

type FiberContext interface {
	Params(name string) string
	Query(name string) string
	BodyParser(out interface{}) error
	Status(code int) *Response
	JSON(data interface{}) error
}

type Response struct {
	c FiberContext
}

func (r *Response) JSON(data interface{}) error {
	return r.c.JSON(data)
}

type FiberMap map[string]interface{}

func Error(message string) FiberMap {
	return FiberMap{"success": false, "error": message}
}

func Success(message string) FiberMap {
	return FiberMap{"success": true, "message": message}
}

func (s *UserService) Register(c *fiber.Ctx) error {
	type RegisterInput struct {
		Email       string `json:"email"`
		Username    string `json:"username"`
		Password    string `json:"password"`
		DisplayName string `json:"displayName"`
	}

	var input RegisterInput
	if err := c.BodyParser(&input); err != nil {
		return c.Status(400).JSON(Error("Invalid request body"))
	}

	if input.Email == "" || input.Username == "" || input.Password == "" {
		return c.Status(400).JSON(Error("email, username and password are required"))
	}
	if len(input.Password) < 8 {
		return c.Status(400).JSON(Error("password must be at least 8 characters"))
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		return c.Status(500).JSON(Error("Failed to hash password"))
	}

	_, err = s.db.ExecContext(context.Background(),
		`INSERT INTO users (email, username, password_hash, display_name) 
		 VALUES ($1, $2, $3, $4)`,
		input.Email,
		input.Username,
		string(hashedPassword),
		input.DisplayName,
	)
	if err != nil {
		return c.Status(400).JSON(Error("Email or username already exists"))
	}

	return c.Status(201).JSON(Success("User registered successfully"))
}

func (s *UserService) Login(c *fiber.Ctx) error {
	type LoginInput struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	var input LoginInput
	if err := c.BodyParser(&input); err != nil {
		return c.Status(400).JSON(Error("Invalid request body"))
	}
	// Normalizar: si el mail llega con espacios o en mayúsculas (copy-paste,
	// autocompletado) el SELECT no encuentra la fila y el usuario recibe un
	// "credenciales inválidas" que no le dice nada.
	input.Email = strings.ToLower(strings.TrimSpace(input.Email))

	var userID, passwordHash string
	err := s.db.QueryRowContext(context.Background(),
		"SELECT id, password_hash FROM users WHERE email = $1",
		input.Email,
	).Scan(&userID, &passwordHash)

	if err == sql.ErrNoRows {
		return c.Status(401).JSON(Error("Ese correo no está registrado"))
	}
	if err != nil {
		return c.Status(500).JSON(Error("Database error"))
	}

	if err := bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(input.Password)); err != nil {
		return c.Status(401).JSON(Error("Contraseña incorrecta"))
	}

	token, err := s.generateToken(userID, input.Email)
	if err != nil {
		return c.Status(500).JSON(Error("Failed to generate token"))
	}

	user, err := s.fetchUser(userID)
	if err != nil {
		return c.Status(500).JSON(Error("Database error"))
	}

	return c.JSON(fiber.Map{
		"success": true,
		"token":   token,
		"userId":  userID,
		"user":    user,
	})
}

func (s *UserService) generateToken(userID, email string) (string, error) {
	return auth.Generate(s.jwtSecret, userID, email, s.ttl)
}

func (s *UserService) Logout(c *fiber.Ctx) error {
	tokenString := strings.TrimSpace(strings.TrimPrefix(c.Get("Authorization"), "Bearer "))
	if tokenString != "" {
		if claims, err := auth.Parse(s.jwtSecret, tokenString); err == nil && claims.ExpiresAt != nil {
			s.revoked.Revoke(tokenString, claims.ExpiresAt.Time)
		}
	}
	return c.JSON(Success("Logged out successfully"))
}

func (s *UserService) GetUser(c *fiber.Ctx) error {
	return s.getUserByID(c, c.Params("id"))
}

// GetProfile devuelve el perfil de un usuario por id.
func (s *UserService) GetProfile(c *fiber.Ctx) error {
	return s.getUserByID(c, c.Params("id"))
}

// GetMe devuelve el usuario autenticado.
func (s *UserService) GetMe(c *fiber.Ctx) error {
	return s.getUserByID(c, middleware.UserID(c))
}

// GetMyStats devuelve contadores del usuario autenticado.
func (s *UserService) GetMyStats(c *fiber.Ctx) error {
	id := middleware.UserID(c)

	query := `SELECT
	  (SELECT COUNT(*) FROM trips WHERE user_id = $1) AS trips,
	  (SELECT COUNT(*) FROM fishing_spots WHERE user_id = $1) AS spots,
	  (SELECT COUNT(*) FROM tournament_participants WHERE user_id = $1) AS tournaments,
	  (SELECT COUNT(*) FROM tracking_sessions WHERE user_id = $1) AS sessions`

	var trips, spots, tournaments, sessions int
	if err := s.db.QueryRowContext(context.Background(), query, id).
		Scan(&trips, &spots, &tournaments, &sessions); err != nil {
		return c.Status(500).JSON(Error("Database error"))
	}

	return c.JSON(fiber.Map{
		"trips":       trips,
		"spots":       spots,
		"tournaments": tournaments,
		"sessions":    sessions,
	})
}

func (s *UserService) getUserByID(c *fiber.Ctx, id string) error {
	user, err := s.fetchUser(id)
	if err == sql.ErrNoRows {
		return c.Status(404).JSON(Error("User not found"))
	}
	if err != nil {
		return c.Status(500).JSON(Error("Database error"))
	}
	return c.JSON(user)
}

// fetchUser devuelve el perfil público de un usuario.
func (s *UserService) fetchUser(id string) (fiber.Map, error) {
	var user struct {
		ID              string
		Email           string
		Username        string
		DisplayName     sql.NullString
		AvatarURL       sql.NullString
		Bio             sql.NullString
		TrackingEnabled bool
		CreatedAt       time.Time
	}

	err := s.db.QueryRowContext(context.Background(),
		`SELECT id, email, username, display_name, avatar_url, bio, tracking_enabled, created_at 
		 FROM users WHERE id = $1`,
		id,
	).Scan(&user.ID, &user.Email, &user.Username, &user.DisplayName, &user.AvatarURL, &user.Bio, &user.TrackingEnabled, &user.CreatedAt)
	if err != nil {
		return nil, err
	}

	return fiber.Map{
		"id":              user.ID,
		"email":           user.Email,
		"username":        user.Username,
		"displayName":     user.DisplayName.String,
		"avatarUrl":       user.AvatarURL.String,
		"bio":             user.Bio.String,
		"trackingEnabled": user.TrackingEnabled,
		"createdAt":       user.CreatedAt,
	}, nil
}

func (s *UserService) UpdateUser(c *fiber.Ctx) error {
	id := c.Params("id")

	type UpdateInput struct {
		DisplayName     *string `json:"displayName"`
		Bio             *string `json:"bio"`
		AvatarURL       *string `json:"avatarUrl"`
		TrackingEnabled *bool   `json:"trackingEnabled"`
	}

	var input UpdateInput
	if err := c.BodyParser(&input); err != nil {
		return c.Status(400).JSON(Error("Invalid request body"))
	}

	_, err := s.db.ExecContext(context.Background(),
		`UPDATE users SET 
		 display_name = COALESCE($1, display_name),
		 bio = COALESCE($2, bio), 
		 avatar_url = COALESCE($3, avatar_url),
		 tracking_enabled = COALESCE($4, tracking_enabled),
		 updated_at = now()
		 WHERE id = $5`,
		input.DisplayName,
		input.Bio,
		input.AvatarURL,
		input.TrackingEnabled,
		id,
	)

	if err != nil {
		return c.Status(500).JSON(Error("Failed to update user"))
	}

	return c.JSON(Success("User updated successfully"))
}

func (s *UserService) GetKayaks(c *fiber.Ctx) error {
	userID := middleware.UserID(c)
	if userID == "" {
		return c.Status(401).JSON(Error("unauthenticated"))
	}

	rows, err := s.db.QueryContext(context.Background(),
		`SELECT id, user_id, name, kaytype, length, width, weight, max_capacity, material, propulsion, image_url 
		 FROM kayaks WHERE user_id = $1`,
		userID,
	)
	if err != nil {
		return c.Status(500).JSON(Error("Database error"))
	}
	defer rows.Close()

	kayaks := []map[string]interface{}{}
	for rows.Next() {
		var id, userID, name string
		var type_, material, propulsion, imageURL sql.NullString
		var length, width, weight, maxCapacity sql.NullFloat64

		if err := rows.Scan(&id, &userID, &name, &type_, &length, &width, &weight, &maxCapacity, &material, &propulsion, &imageURL); err != nil {
			continue
		}
		kayak := map[string]interface{}{
			"id":     id,
			"userId": userID,
			"name":   name,
		}
		if type_.Valid {
			kayak["type"] = type_.String
		}
		if length.Valid {
			kayak["length"] = length.Float64
		}
		if material.Valid {
			kayak["material"] = material.String
		}
		if propulsion.Valid {
			kayak["propulsion"] = propulsion.String
		}
		if imageURL.Valid {
			kayak["imageUrl"] = imageURL.String
		}
		kayaks = append(kayaks, kayak)
	}

	return c.JSON(kayaks)
}

func (s *UserService) CreateKayak(c *fiber.Ctx) error {
	type KayakInput struct {
		UserID      string  `json:"userId"`
		Name        string  `json:"name"`
		Type        string  `json:"type"`
		Length      float64 `json:"length"`
		Width       float64 `json:"width"`
		Weight      float64 `json:"weight"`
		MaxCapacity float64 `json:"maxCapacity"`
		Material    string  `json:"material"`
		Propulsion  string  `json:"propulsion"`
		ImageURL    string  `json:"imageUrl"`
	}

	var input KayakInput
	if err := c.BodyParser(&input); err != nil {
		return c.Status(400).JSON(Error("Invalid request body"))
	}
	input.UserID = middleware.UserID(c)

	var id string
	err := s.db.QueryRowContext(context.Background(),
		`INSERT INTO kayaks (user_id, name, kaytype, length, width, weight, max_capacity, material, propulsion, image_url)
		 VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
		 RETURNING id::text`,
		input.UserID,
		input.Name,
		input.Type,
		input.Length,
		input.Width,
		input.Weight,
		input.MaxCapacity,
		input.Material,
		input.Propulsion,
		input.ImageURL,
	).Scan(&id)

	if err != nil {
		return c.Status(500).JSON(Error("Failed to create kayak"))
	}

	return c.Status(201).JSON(fiber.Map{"id": id, "success": true})
}

func (s *UserService) UpdateKayak(c *fiber.Ctx) error {
	id := c.Params("id")

	type KayakInput struct {
		Name        *string  `json:"name"`
		Type        *string  `json:"type"`
		Length      *float64 `json:"length"`
		Width       *float64 `json:"width"`
		Weight      *float64 `json:"weight"`
		MaxCapacity *float64 `json:"maxCapacity"`
		Material    *string  `json:"material"`
		Propulsion  *string  `json:"propulsion"`
		ImageURL    *string  `json:"imageUrl"`
	}

	var input KayakInput
	if err := c.BodyParser(&input); err != nil {
		return c.Status(400).JSON(Error("Invalid request body"))
	}

	_, err := s.db.ExecContext(context.Background(),
		`UPDATE kayaks SET name = COALESCE($1, name), kaytype = COALESCE($2, kaytype),
		 length = COALESCE($3, length), width = COALESCE($4, width),
		 weight = COALESCE($5, weight), max_capacity = COALESCE($6, max_capacity),
		 material = COALESCE($7, material), propulsion = COALESCE($8, propulsion),
		 image_url = COALESCE($9, image_url)
		 WHERE id = $10`,
		input.Name,
		input.Type,
		input.Length,
		input.Width,
		input.Weight,
		input.MaxCapacity,
		input.Material,
		input.Propulsion,
		input.ImageURL,
		id,
	)

	if err != nil {
		return c.Status(500).JSON(Error("Failed to update kayak"))
	}

	return c.JSON(Success("Kayak updated"))
}

func (s *UserService) DeleteKayak(c *fiber.Ctx) error {
	id := c.Params("id")

	_, err := s.db.ExecContext(context.Background(),
		"DELETE FROM kayaks WHERE id = $1",
		id,
	)

	if err != nil {
		return c.Status(500).JSON(Error("Failed to delete kayak"))
	}

	return c.JSON(Success("Kayak deleted"))
}
