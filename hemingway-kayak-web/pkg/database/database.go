package database

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/microsoft/go-mssqldb"
	"github.com/redis/go-redis/v9"
)

var (
	server   = "10.10.1.21"
	port     = 1433
	user     = "sa"
	password = "[REDACTED]"
	dbname   = "Hemingway"
)

func Connect() (*sql.DB, error) {
	connString := fmt.Sprintf(
		"server=%s;port=%d;database=%s;user id=%s;password=%s;encrypt=disable",
		server, port, dbname, user, password,
	)

	db, err := sql.Open("sqlserver", connString)
	if err != nil {
		return nil, fmt.Errorf("failed to open connection: %w", err)
	}

	ctx := context.Background()
	if err := db.PingContext(ctx); err != nil {
		return nil, fmt.Errorf("failed to ping: %w", err)
	}

	log.Println("Connected to MS SQL Server")
	return db, nil
}

func EnsureDatabase() error {
	connString := fmt.Sprintf(
		"server=%s;port=%d;database=master;user id=%s;password=%s;encrypt=disable",
		server, port, user, password,
	)

	db, err := sql.Open("sqlserver", connString)
	if err != nil {
		return fmt.Errorf("failed to open master connection: %w", err)
	}
	defer db.Close()

	ctx := context.Background()
	if err := db.PingContext(ctx); err != nil {
		return fmt.Errorf("failed to ping master: %w", err)
	}

	var exists int
	err = db.QueryRowContext(ctx, fmt.Sprintf("SELECT COUNT(*) FROM sys.databases WHERE name = '%s'", dbname)).Scan(&exists)
	if err != nil {
		return fmt.Errorf("failed to check database existence: %w", err)
	}

	if exists == 0 {
		_, err = db.ExecContext(ctx, fmt.Sprintf("CREATE DATABASE %s", dbname))
		if err != nil {
			return fmt.Errorf("failed to create database: %w", err)
		}
		log.Printf("Database '%s' created successfully", dbname)
	} else {
		log.Printf("Database '%s' already exists", dbname)
	}

	return nil
}

func ConnectRedis() *redis.Client {
	redisAddr := os.Getenv("REDIS_ADDR")
	if redisAddr == "" {
		redisAddr = "localhost:6379"
	}

	client := redis.NewClient(&redis.Options{
		Addr:     redisAddr,
		Password: os.Getenv("REDIS_PASSWORD"),
		DB:       0,
	})

	ctx := context.Background()
	if _, err := client.Ping(ctx).Result(); err != nil {
		log.Printf("Warning: Redis connection failed: %v (continuing without Redis)", err)
		return nil
	}

	log.Println("Connected to Redis")
	return client
}

func RunMigrations() error {
	db, err := Connect()
	if err != nil {
		return fmt.Errorf("failed to connect: %w", err)
	}
	defer db.Close()

	ctx := context.Background()

	migrations := []string{
		`IF NOT EXISTS (SELECT * FROM sys.tables WHERE name = 'users')
		CREATE TABLE users (
			id UNIQUEIDENTIFIER PRIMARY KEY DEFAULT NEWID(),
			email NVARCHAR(255) UNIQUE NOT NULL,
			username NVARCHAR(100) UNIQUE NOT NULL,
			password_hash NVARCHAR(255) NOT NULL,
			display_name NVARCHAR(200),
			avatar_url NVARCHAR(500),
			bio NVARCHAR(1000),
			tracking_enabled BIT DEFAULT 0,
			created_at DATETIME2 DEFAULT GETDATE(),
			updated_at DATETIME2 DEFAULT GETDATE()
		)`,

		`IF NOT EXISTS (SELECT * FROM sys.tables WHERE name = 'kayaks')
		CREATE TABLE kayaks (
			id UNIQUEIDENTIFIER PRIMARY KEY DEFAULT NEWID(),
			user_id UNIQUEIDENTIFIER REFERENCES users(id),
			name NVARCHAR(200) NOT NULL,
			[kaytype] NVARCHAR(50),
			length FLOAT,
			width FLOAT,
			weight FLOAT,
			max_capacity FLOAT,
			material NVARCHAR(50),
			propulsion NVARCHAR(50),
			image_url NVARCHAR(500),
			created_at DATETIME2 DEFAULT GETDATE()
		)`,

		`IF NOT EXISTS (SELECT * FROM sys.tables WHERE name = 'trips')
		CREATE TABLE trips (
			id UNIQUEIDENTIFIER PRIMARY KEY DEFAULT NEWID(),
			user_id UNIQUEIDENTIFIER REFERENCES users(id),
			title NVARCHAR(300) NOT NULL,
			description NVARCHAR(MAX),
			[triptype] NVARCHAR(50) NOT NULL,
			status NVARCHAR(50) DEFAULT 'planned',
			start_time DATETIME2,
			end_time DATETIME2,
			start_lat FLOAT,
			start_lon FLOAT,
			end_lat FLOAT,
			end_lon FLOAT,
			distance FLOAT DEFAULT 0,
			created_at DATETIME2 DEFAULT GETDATE()
		)`,

		`IF NOT EXISTS (SELECT * FROM sys.tables WHERE name = 'trip_points')
		CREATE TABLE trip_points (
			id UNIQUEIDENTIFIER PRIMARY KEY DEFAULT NEWID(),
			trip_id UNIQUEIDENTIFIER REFERENCES trips(id),
			latitude FLOAT NOT NULL,
			longitude FLOAT NOT NULL,
			altitude FLOAT,
			speed FLOAT,
			heading FLOAT,
			timestamp DATETIME2 NOT NULL
		)`,

		`IF NOT EXISTS (SELECT * FROM sys.tables WHERE name = 'tournaments')
		CREATE TABLE tournaments (
			id UNIQUEIDENTIFIER PRIMARY KEY DEFAULT NEWID(),
			name NVARCHAR(300) NOT NULL,
			description NVARCHAR(MAX),
			organizer_id UNIQUEIDENTIFIER REFERENCES users(id),
			start_time DATETIME2 NOT NULL,
			end_time DATETIME2 NOT NULL,
			location_name NVARCHAR(500),
			start_lat FLOAT,
			start_lon FLOAT,
			status NVARCHAR(50) DEFAULT 'planned',
			[tournamenttype] NVARCHAR(50),
			max_participants INT,
			created_at DATETIME2 DEFAULT GETDATE()
		)`,

		`IF NOT EXISTS (SELECT * FROM sys.tables WHERE name = 'tournament_participants')
		CREATE TABLE tournament_participants (
			id UNIQUEIDENTIFIER PRIMARY KEY DEFAULT NEWID(),
			tournament_id UNIQUEIDENTIFIER REFERENCES tournaments(id),
			user_id UNIQUEIDENTIFIER REFERENCES users(id),
			kayak_id UNIQUEIDENTIFIER REFERENCES kayaks(id),
			status NVARCHAR(50) DEFAULT 'registered',
			current_lat FLOAT,
			current_lon FLOAT,
			score FLOAT DEFAULT 0,
			catches_count INT DEFAULT 0,
			last_update_time DATETIME2,
			joined_at DATETIME2 DEFAULT GETDATE()
		)`,

		`IF NOT EXISTS (SELECT * FROM sys.tables WHERE name = 'tournament_catches')
		CREATE TABLE tournament_catches (
			id UNIQUEIDENTIFIER PRIMARY KEY DEFAULT NEWID(),
			tournament_id UNIQUEIDENTIFIER REFERENCES tournaments(id),
			participant_id UNIQUEIDENTIFIER REFERENCES tournament_participants(id),
			species NVARCHAR(200) NOT NULL,
			weight FLOAT NOT NULL,
			latitude FLOAT,
			longitude FLOAT,
			photo_url NVARCHAR(500),
			verified BIT DEFAULT 0,
			verified_by NVARCHAR(200),
			created_at DATETIME2 DEFAULT GETDATE()
		)`,

		`IF NOT EXISTS (SELECT * FROM sys.tables WHERE name = 'tracking_sessions')
		CREATE TABLE tracking_sessions (
			id UNIQUEIDENTIFIER PRIMARY KEY DEFAULT NEWID(),
			user_id UNIQUEIDENTIFIER REFERENCES users(id),
			trip_id UNIQUEIDENTIFIER REFERENCES trips(id),
			start_time DATETIME2 DEFAULT GETDATE(),
			end_time DATETIME2,
			start_lat FLOAT,
			start_lon FLOAT,
			end_lat FLOAT,
			end_lon FLOAT,
			total_distance FLOAT DEFAULT 0,
			is_live BIT DEFAULT 0,
			viewer_count INT DEFAULT 0
		)`,

		`IF NOT EXISTS (SELECT * FROM sys.tables WHERE name = 'live_tracking')
		CREATE TABLE live_tracking (
			user_id UNIQUEIDENTIFIER PRIMARY KEY,
			trip_id UNIQUEIDENTIFIER REFERENCES trips(id),
			latitude FLOAT NOT NULL,
			longitude FLOAT NOT NULL,
			speed FLOAT,
			heading FLOAT,
			altitude FLOAT,
			timestamp DATETIME2 DEFAULT GETDATE(),
			is_active BIT DEFAULT 1
		)`,

		`IF NOT EXISTS (SELECT * FROM sys.tables WHERE name = 'history_entries')
		CREATE TABLE history_entries (
			id UNIQUEIDENTIFIER PRIMARY KEY DEFAULT NEWID(),
			user_id UNIQUEIDENTIFIER REFERENCES users(id),
			[historytype] NVARCHAR(50) NOT NULL,
			reference_id UNIQUEIDENTIFIER,
			title NVARCHAR(500) NOT NULL,
			description NVARCHAR(MAX),
			data NVARCHAR(MAX),
			occurred_at DATETIME2 NOT NULL,
			created_at DATETIME2 DEFAULT GETDATE()
		)`,

		`IF NOT EXISTS (SELECT * FROM sys.tables WHERE name = 'fishing_spots')
		CREATE TABLE fishing_spots (
			id UNIQUEIDENTIFIER PRIMARY KEY DEFAULT NEWID(),
			user_id UNIQUEIDENTIFIER REFERENCES users(id),
			name NVARCHAR(300) NOT NULL,
			description NVARCHAR(MAX),
			latitude FLOAT NOT NULL,
			longitude FLOAT NOT NULL,
			rating INT DEFAULT 0,
			fish_caught INT DEFAULT 0,
			is_shared BIT DEFAULT 0,
			created_at DATETIME2 DEFAULT GETDATE()
		)`,
	}

	for i, migration := range migrations {
		_, err := db.ExecContext(ctx, migration)
		if err != nil {
			log.Printf("Migration %d: %v", i+1, err)
		}
	}

	log.Println("Migrations completed")
	return nil
}

func SeedData() error {
	db, err := Connect()
	if err != nil {
		return fmt.Errorf("failed to connect: %w", err)
	}
	defer db.Close()

	ctx := context.Background()

	var count int
	err = db.QueryRowContext(ctx, "SELECT COUNT(*) FROM users").Scan(&count)
	if err != nil {
		return fmt.Errorf("failed to check users: %w", err)
	}

	if count > 0 {
		log.Println("Seed data already exists, skipping...")
		return nil
	}

	seeds := []string{
		`INSERT INTO users (id, email, username, password_hash, display_name, tracking_enabled)
		VALUES (NEWID(), 'peter@hemingway.com', 'peter_kayak', '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZRGdjGj/n3.QZjgM0pJ7OLLRzXOLa', 'Peter Garcia', 1)`,

		`INSERT INTO users (id, email, username, password_hash, display_name, tracking_enabled)
		VALUES (NEWID(), 'marina@hemingway.com', 'marina_kayak', '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZRGdjGj/n3.QZjgM0pJ7OLLRzXOLa', 'Marina Rodriguez', 1)`,

		`INSERT INTO users (id, email, username, password_hash, display_name, tracking_enabled)
		VALUES (NEWID(), 'carlos@hemingway.com', 'carlos_kayak', '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZRGdjGj/n3.QZjgM0pJ7OLLRzXOLa', 'Carlos Mendoza', 0)`,

		`INSERT INTO kayaks (id, user_id, name, kaytype, length, width, weight, max_capacity, material, propulsion)
		SELECT NEWID(), id, 'Explorer Pro', 'SIT_ON_TOP', 380, 85, 35, 250, 'FIBERGLASS', 'PADDLE' FROM users WHERE email = 'peter@hemingway.com'`,

		`INSERT INTO kayaks (id, user_id, name, kaytype, length, width, weight, max_capacity, material, propulsion)
		SELECT NEWID(), id, 'Trawler 12', 'SIT_INSIDE', 360, 75, 30, 200, 'PLASTIC', 'PEDAL' FROM users WHERE email = 'peter@hemingway.com'`,

		`INSERT INTO kayaks (id, user_id, name, kaytype, length, width, weight, max_capacity, material, propulsion)
		SELECT NEWID(), id, 'Sportster', 'SIT_ON_TOP', 340, 80, 28, 220, 'CARBON', 'PADDLE' FROM users WHERE email = 'marina@hemingway.com'`,

		`INSERT INTO trips (id, user_id, title, triptype, status, start_time, end_time, start_lat, start_lon, end_lat, end_lon, distance)
		SELECT NEWID(), id, 'Miami Everglades Expedition', 'FISHING', 'COMPLETED', '2026-03-15 06:00:00', '2026-03-15 14:30:00', 25.7617, -80.1918, 25.8900, -80.3600, 28.5 FROM users WHERE email = 'peter@hemingway.com'`,

		`INSERT INTO trips (id, user_id, title, triptype, status, start_time, end_time, start_lat, start_lon, distance)
		SELECT NEWID(), id, 'Keys Weekend Adventure', 'FISHING', 'PLANNED', '2026-03-28 05:00:00', '2026-03-30 18:00:00', 24.5554, -81.7800, 0 FROM users WHERE email = 'peter@hemingway.com'`,

		`INSERT INTO trips (id, user_id, title, triptype, status, start_time, end_time, start_lat, start_lon, end_lat, end_lon, distance)
		SELECT NEWID(), id, 'Biscayne Bay Morning', 'FISHING', 'COMPLETED', '2026-03-20 06:30:00', '2026-03-20 11:00:00', 25.5100, -80.2000, 25.6500, -80.2500, 15.2 FROM users WHERE email = 'marina@hemingway.com'`,

		`INSERT INTO tournaments (id, name, organizer_id, start_time, end_time, location_name, start_lat, start_lon, status, tournamenttype, max_participants)
		SELECT NEWID(), 'Spring Kingfish Championship 2026', id, '2026-04-05 07:00:00', '2026-04-05 17:00:00', 'Miami Marina', 25.7617, -80.1918, 'PLANNED', 'FISHING', 50 FROM users WHERE email = 'peter@hemingway.com'`,

		`INSERT INTO tournaments (id, name, organizer_id, start_time, end_time, location_name, start_lat, start_lon, status, tournamenttype, max_participants)
		SELECT NEWID(), 'Kayak Racing League - Miami', id, '2026-03-22 08:00:00', '2026-03-22 14:00:00', 'Crandon Park', 25.5100, -80.2000, 'LIVE', 'RACING', 30 FROM users WHERE email = 'marina@hemingway.com'`,

		`INSERT INTO tournament_participants (id, tournament_id, user_id, kayak_id, status, current_lat, current_lon, score, catches_count)
		SELECT NEWID(), t.id, u.id, k.id, 'ACTIVE', 25.5200, -80.2100, 85.5, 3
		FROM tournaments t, users u, kayaks k
		WHERE t.name = 'Kayak Racing League - Miami' AND u.email = 'peter@hemingway.com' AND k.name = 'Explorer Pro'`,

		`INSERT INTO tournament_participants (id, tournament_id, user_id, kayak_id, status, current_lat, current_lon, score, catches_count)
		SELECT NEWID(), t.id, u.id, k.id, 'ACTIVE', 25.5300, -80.2200, 72.0, 2
		FROM tournaments t, users u, kayaks k
		WHERE t.name = 'Kayak Racing League - Miami' AND u.email = 'marina@hemingway.com' AND k.name = 'Sportster'`,

		`INSERT INTO live_tracking (user_id, trip_id, latitude, longitude, speed, is_active)
		SELECT u.id, t.id, 25.5250, -80.2050, 8.5, 1
		FROM users u, trips t
		WHERE u.email = 'peter@hemingway.com' AND t.title = 'Biscayne Bay Morning'`,

		`INSERT INTO live_tracking (user_id, trip_id, latitude, longitude, speed, is_active)
		SELECT u.id, t.id, 25.5350, -80.2150, 6.2, 1
		FROM users u, trips t
		WHERE u.email = 'marina@hemingway.com' AND t.title = 'Biscayne Bay Morning'`,

		`INSERT INTO fishing_spots (id, user_id, name, description, latitude, longitude, rating, fish_caught, is_shared)
		SELECT NEWID(), id, 'Secret Point', 'Great for snook at dawn', 25.8500, -80.3500, 5, 47, 1 FROM users WHERE email = 'peter@hemingway.com'`,

		`INSERT INTO fishing_spots (id, user_id, name, description, latitude, longitude, rating, fish_caught, is_shared)
		SELECT NEWID(), id, 'North Bay Spots', 'Multiple species, deep water', 25.9500, -80.2800, 4, 32, 1 FROM users WHERE email = 'peter@hemingway.com'`,

		`INSERT INTO fishing_spots (id, user_id, name, description, latitude, longitude, rating, fish_caught, is_shared)
		SELECT NEWID(), id, 'Biscayne Inlet', 'Perfect for tarpon', 25.5100, -80.2000, 5, 28, 1 FROM users WHERE email = 'marina@hemingway.com'`,
	}

	for i, seed := range seeds {
		_, err := db.ExecContext(ctx, seed)
		if err != nil {
			log.Printf("Seed %d: %v", i+1, err)
		}
	}

	log.Println("Seed data inserted successfully")
	return nil
}
