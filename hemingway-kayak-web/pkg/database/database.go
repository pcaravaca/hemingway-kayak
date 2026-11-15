// Package database concentra la conexión a MS SQL Server, Redis y la aplicación
// de migraciones/seed. Las credenciales salen de pkg/config (variables de entorno).
package database

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"sort"
	"strings"

	_ "github.com/microsoft/go-mssqldb"
	"github.com/redis/go-redis/v9"
	"golang.org/x/crypto/bcrypt"

	"hemingway-kayak-web/migrations"
	"hemingway-kayak-web/pkg/config"
)

var cfg = config.Load()

// Connect abre una conexión a la base de datos de la aplicación.
func Connect() (*sql.DB, error) {
	db, err := sql.Open("sqlserver", cfg.DSN(cfg.DBName))
	if err != nil {
		return nil, fmt.Errorf("failed to open connection: %w", err)
	}

	ctx := context.Background()
	if err := db.PingContext(ctx); err != nil {
		return nil, fmt.Errorf("failed to ping: %w", err)
	}

	log.Printf("Connected to MS SQL Server (%s:%d/%s)", cfg.DBServer, cfg.DBPort, cfg.DBName)
	return db, nil
}

// EnsureDatabase crea la base si no existe. Requiere conectarse a master.
func EnsureDatabase() error {
	db, err := sql.Open("sqlserver", cfg.DSN("master"))
	if err != nil {
		return fmt.Errorf("failed to open master connection: %w", err)
	}
	defer db.Close()

	ctx := context.Background()
	if err := db.PingContext(ctx); err != nil {
		return fmt.Errorf("failed to ping master: %w", err)
	}

	var exists int
	err = db.QueryRowContext(ctx,
		"SELECT COUNT(*) FROM sys.databases WHERE name = @name",
		sql.Named("name", cfg.DBName),
	).Scan(&exists)
	if err != nil {
		return fmt.Errorf("failed to check database existence: %w", err)
	}

	if exists == 0 {
		if _, err := db.ExecContext(ctx, "CREATE DATABASE "+quoteIdent(cfg.DBName)); err != nil {
			return fmt.Errorf("failed to create database: %w", err)
		}
		log.Printf("Database '%s' created successfully", cfg.DBName)
	} else {
		log.Printf("Database '%s' already exists", cfg.DBName)
	}

	return nil
}

// ConnectRedis devuelve el cliente Redis o nil si no está disponible.
// El backend debe verificar el nil antes de usarlo.
func ConnectRedis() *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr:     cfg.RedisAddr,
		Password: cfg.RedisPassword,
		DB:       0,
	})

	ctx := context.Background()
	if _, err := client.Ping(ctx).Result(); err != nil {
		log.Printf("Warning: Redis connection failed (%v); continuing without Redis", err)
		return nil
	}

	log.Printf("Connected to Redis (%s)", cfg.RedisAddr)
	return client
}

// RunMigrations aplica los archivos *.sql de migrations/ en orden, una sola vez,
// registrándolos en schema_migrations. Falla en la primera migración con error.
func RunMigrations() error {
	db, err := Connect()
	if err != nil {
		return fmt.Errorf("failed to connect: %w", err)
	}
	defer db.Close()

	ctx := context.Background()

	if _, err := db.ExecContext(ctx, `CREATE TABLE IF NOT EXISTS schema_migrations (
		filename NVARCHAR(255) PRIMARY KEY,
		applied_at DATETIME2 DEFAULT GETDATE()
	)`); err != nil {
		return fmt.Errorf("failed to create schema_migrations: %w", err)
	}

	entries, err := migrations.FS.ReadDir(".")
	if err != nil {
		return fmt.Errorf("failed to read embedded migrations: %w", err)
	}

	var names []string
	for _, e := range entries {
		if !e.IsDir() && strings.HasSuffix(e.Name(), ".sql") {
			names = append(names, e.Name())
		}
	}
	sort.Strings(names)

	applied := map[string]bool{}
	rows, err := db.QueryContext(ctx, "SELECT filename FROM schema_migrations")
	if err != nil {
		return fmt.Errorf("failed to read schema_migrations: %w", err)
	}
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err == nil {
			applied[name] = true
		}
	}
	rows.Close()

	for _, name := range names {
		if applied[name] {
			continue
		}

		content, err := migrations.FS.ReadFile(name)
		if err != nil {
			return fmt.Errorf("failed to read migration %s: %w", name, err)
		}

		if _, err := db.ExecContext(ctx, string(content)); err != nil {
			return fmt.Errorf("migration %s failed: %w", name, err)
		}

		if _, err := db.ExecContext(ctx,
			"INSERT INTO schema_migrations (filename) VALUES (@name)",
			sql.Named("name", name),
		); err != nil {
			return fmt.Errorf("failed to record migration %s: %w", name, err)
		}

		log.Printf("Migration applied: %s", name)
	}

	log.Println("Migrations completed")
	return nil
}

// seedPassword es la contraseña de las cuentas de demostración que inserta SeedData.
const seedPassword = "HemingwayKayak2026!"

// SeedData inserta datos de demostración solo si la tabla users está vacía.
func SeedData() error {
	db, err := Connect()
	if err != nil {
		return fmt.Errorf("failed to connect: %w", err)
	}
	defer db.Close()

	ctx := context.Background()

	var count int
	if err := db.QueryRowContext(ctx, "SELECT COUNT(*) FROM users").Scan(&count); err != nil {
		return fmt.Errorf("failed to check users: %w", err)
	}
	if count > 0 {
		log.Println("Seed data already exists, skipping...")
		return nil
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(seedPassword), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("failed to hash seed password: %w", err)
	}

	seeds := []string{
		`INSERT INTO users (id, email, username, password_hash, display_name, tracking_enabled)
		 VALUES (NEWID(), 'peter@hemingway.com', 'peter_kayak', @hash, 'Peter Garcia', 1)`,

		`INSERT INTO users (id, email, username, password_hash, display_name, tracking_enabled)
		 VALUES (NEWID(), 'marina@hemingway.com', 'marina_kayak', @hash, 'Marina Rodriguez', 1)`,

		`INSERT INTO users (id, email, username, password_hash, display_name, tracking_enabled)
		 VALUES (NEWID(), 'carlos@hemingway.com', 'carlos_kayak', @hash, 'Carlos Mendoza', 0)`,

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
		 WHERE u.email = 'peter@hemingway.com' AND t.title = 'Miami Everglades Expedition'`,

		`INSERT INTO live_tracking (user_id, trip_id, latitude, longitude, speed, is_active)
		 SELECT u.id, t.id, 25.5350, -80.2150, 6.2, 1
		 FROM users u, trips t
		 WHERE u.email = 'marina@hemingway.com' AND t.title = 'Biscayne Bay Morning'`,

		`INSERT INTO tracking_sessions (id, user_id, trip_id, start_time, start_lat, start_lon, total_distance, is_live, viewer_count)
		 SELECT NEWID(), u.id, t.id, '2026-03-15 06:00:00', 25.7617, -80.1918, 28.5, 1, 3
		 FROM users u, trips t
		 WHERE u.email = 'peter@hemingway.com' AND t.title = 'Miami Everglades Expedition'`,

		`INSERT INTO fishing_spots (id, user_id, name, description, latitude, longitude, rating, fish_caught, is_shared)
		 SELECT NEWID(), id, 'Secret Point', 'Great for snook at dawn', 25.8500, -80.3500, 5, 47, 1 FROM users WHERE email = 'peter@hemingway.com'`,

		`INSERT INTO fishing_spots (id, user_id, name, description, latitude, longitude, rating, fish_caught, is_shared)
		 SELECT NEWID(), id, 'North Bay Spots', 'Multiple species, deep water', 25.9500, -80.2800, 4, 32, 1 FROM users WHERE email = 'peter@hemingway.com'`,

		`INSERT INTO fishing_spots (id, user_id, name, description, latitude, longitude, rating, fish_caught, is_shared)
		 SELECT NEWID(), id, 'Biscayne Inlet', 'Perfect for tarpon', 25.5100, -80.2000, 5, 28, 1 FROM users WHERE email = 'marina@hemingway.com'`,
	}

	for i, seed := range seeds {
		if _, err := db.ExecContext(ctx, seed, sql.Named("hash", string(hash))); err != nil {
			return fmt.Errorf("seed %d failed: %w", i+1, err)
		}
	}

	log.Println("Seed data inserted successfully")
	return nil
}

// quoteIdent envuelve un identificador SQL entre corchetes, escapando los ].
func quoteIdent(name string) string {
	return "[" + strings.ReplaceAll(name, "]", "]]") + "]"
}
