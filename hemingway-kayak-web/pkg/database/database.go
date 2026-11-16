// Package database concentra la conexión a PostgreSQL, Redis y la aplicación de
// migraciones/seed. Las credenciales salen de pkg/config (variables de entorno).
package database

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"sort"
	"strings"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/redis/go-redis/v9"
	"golang.org/x/crypto/bcrypt"

	"hemingway-kayak-web/migrations"
	"hemingway-kayak-web/pkg/config"
)

var cfg = config.Load()

// Connect abre una conexión a la base de datos de la aplicación.
func Connect() (*sql.DB, error) {
	db, err := sql.Open("pgx", cfg.DSN(cfg.DBName))
	if err != nil {
		return nil, fmt.Errorf("failed to open connection: %w", err)
	}

	ctx := context.Background()
	if err := db.PingContext(ctx); err != nil {
		return nil, fmt.Errorf("failed to ping: %w", err)
	}

	log.Printf("Connected to PostgreSQL (%s:%d/%s)", cfg.DBServer, cfg.DBPort, cfg.DBName)
	return db, nil
}

// EnsureDatabase crea la base si no existe. Requiere conectarse a la base 'postgres'.
func EnsureDatabase() error {
	db, err := sql.Open("pgx", cfg.DSN("postgres"))
	if err != nil {
		return fmt.Errorf("failed to open admin connection: %w", err)
	}
	defer db.Close()

	ctx := context.Background()
	if err := db.PingContext(ctx); err != nil {
		return fmt.Errorf("failed to ping admin: %w", err)
	}

	var exists int
	if err := db.QueryRowContext(ctx,
		"SELECT COUNT(*) FROM pg_database WHERE datname = $1", cfg.DBName,
	).Scan(&exists); err != nil {
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
		filename TEXT PRIMARY KEY,
		applied_at TIMESTAMPTZ DEFAULT now()
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
			"INSERT INTO schema_migrations (filename) VALUES ($1)", name,
		); err != nil {
			return fmt.Errorf("failed to record migration %s: %w", name, err)
		}

		log.Printf("Migration applied: %s", name)
	}

	log.Println("Migrations completed")
	return nil
}

// seedPassword es la contraseña de las cuentas de demostración que inserta SeedData.
// Todas las cuentas de ejemplo la comparten: es para probar la app, no para producción.
const seedPassword = "HemingwayKayak2026!"

// SeedData inserta un conjunto de demostración completo (usuarios, kayaks, viajes
// con track, torneos con participantes y capturas, historial y spots) solo si la
// tabla users está vacía.
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

	seeds := []struct {
		query string
		args  []interface{}
	}{
		// ---------- Usuarios ----------
		{
			`INSERT INTO users (id, email, username, password_hash, display_name, bio, tracking_enabled)
			 VALUES (gen_random_uuid(), 'peter@hemingway.com', 'peter_kayak', $1, 'Peter Garcia',
			         'Pescador de kayak desde 2015. Snook y tarpon.', TRUE)`,
			[]interface{}{string(hash)},
		},
		{
			`INSERT INTO users (id, email, username, password_hash, display_name, bio, tracking_enabled)
			 VALUES (gen_random_uuid(), 'marina@hemingway.com', 'marina_kayak', $1, 'Marina Rodriguez',
			         'Campeona de kayak racing en Florida.', TRUE)`,
			[]interface{}{string(hash)},
		},
		{
			`INSERT INTO users (id, email, username, password_hash, display_name, bio, tracking_enabled)
			 VALUES (gen_random_uuid(), 'carlos@hemingway.com', 'carlos_kayak', $1, 'Carlos Mendoza',
			         'Pesca de costa y kayak de travesia.', FALSE)`,
			[]interface{}{string(hash)},
		},
		{
			`INSERT INTO users (id, email, username, password_hash, display_name, bio, tracking_enabled)
			 VALUES (gen_random_uuid(), 'lucia@hemingway.com', 'lucia_kayak', $1, 'Lucia Fernandez',
			         'Fotografa de naturaleza y pescadora.', TRUE)`,
			[]interface{}{string(hash)},
		},
		{
			`INSERT INTO users (id, email, username, password_hash, display_name, bio, tracking_enabled)
			 VALUES (gen_random_uuid(), 'diego@hemingway.com', 'diego_kayak', $1, 'Diego Torres',
			         'Entrenador de paddle surf.', TRUE)`,
			[]interface{}{string(hash)},
		},
		{
			`INSERT INTO users (id, email, username, password_hash, display_name, bio, tracking_enabled)
			 VALUES (gen_random_uuid(), 'ana@hemingway.com', 'ana_kayak', $1, 'Ana Morales',
			         'Pesca en familia los fines de semana.', FALSE)`,
			[]interface{}{string(hash)},
		},

		// ---------- Kayaks ----------
		{
			`INSERT INTO kayaks (id, user_id, name, kaytype, length, width, weight, max_capacity, material, propulsion, image_url)
			 SELECT gen_random_uuid(), id, 'Explorer Pro', 'SIT_ON_TOP', 380, 85, 35, 250, 'FIBERGLASS', 'PADDLE', NULL
			 FROM users WHERE email = 'peter@hemingway.com'`,
			nil,
		},
		{
			`INSERT INTO kayaks (id, user_id, name, kaytype, length, width, weight, max_capacity, material, propulsion, image_url)
			 SELECT gen_random_uuid(), id, 'Trawler 12', 'SIT_INSIDE', 360, 75, 30, 200, 'PLASTIC', 'PEDAL', NULL
			 FROM users WHERE email = 'peter@hemingway.com'`,
			nil,
		},
		{
			`INSERT INTO kayaks (id, user_id, name, kaytype, length, width, weight, max_capacity, material, propulsion, image_url)
			 SELECT gen_random_uuid(), id, 'Sportster', 'SIT_ON_TOP', 340, 80, 28, 220, 'CARBON', 'PADDLE', NULL
			 FROM users WHERE email = 'marina@hemingway.com'`,
			nil,
		},
		{
			`INSERT INTO kayaks (id, user_id, name, kaytype, length, width, weight, max_capacity, material, propulsion, image_url)
			 SELECT gen_random_uuid(), id, 'Carbon Razor', 'RACING', 520, 52, 18, 130, 'CARBON', 'PADDLE', NULL
			 FROM users WHERE email = 'marina@hemingway.com'`,
			nil,
		},
		{
			`INSERT INTO kayaks (id, user_id, name, kaytype, length, width, weight, max_capacity, material, propulsion, image_url)
			 SELECT gen_random_uuid(), id, 'Bay Stealth', 'SIT_INSIDE', 410, 78, 33, 230, 'FIBERGLASS', 'PADDLE', NULL
			 FROM users WHERE email = 'lucia@hemingway.com'`,
			nil,
		},
		{
			`INSERT INTO kayaks (id, user_id, name, kaytype, length, width, weight, max_capacity, material, propulsion, image_url)
			 SELECT gen_random_uuid(), id, 'Pedal Duke', 'PEDAL', 390, 88, 42, 260, 'PLASTIC', 'PEDAL', NULL
			 FROM users WHERE email = 'diego@hemingway.com'`,
			nil,
		},
		{
			`INSERT INTO kayaks (id, user_id, name, kaytype, length, width, weight, max_capacity, material, propulsion, image_url)
			 SELECT gen_random_uuid(), id, 'Family Tandem', 'SIT_ON_TOP', 450, 95, 48, 320, 'PLASTIC', 'PADDLE', NULL
			 FROM users WHERE email = 'ana@hemingway.com'`,
			nil,
		},

		// ---------- Viajes ----------
		{
			`INSERT INTO trips (id, user_id, title, description, triptype, status, start_time, end_time, start_lat, start_lon, end_lat, end_lon, distance)
			 SELECT gen_random_uuid(), id, 'Miami Everglades Expedition', 'Travesia por los canales del Everglades buscando snook.', 'FISHING', 'COMPLETED',
			        '2026-03-15 06:00:00', '2026-03-15 14:30:00', 25.7617, -80.1918, 25.8900, -80.3600, 28.5
			 FROM users WHERE email = 'peter@hemingway.com'`,
			nil,
		},
		{
			`INSERT INTO trips (id, user_id, title, description, triptype, status, start_time, end_time, start_lat, start_lon, distance)
			 SELECT gen_random_uuid(), id, 'Keys Weekend Adventure', 'Fin de semana en los cayos con campamento.', 'CAMPING', 'PLANNED',
			        '2026-03-28 05:00:00', '2026-03-30 18:00:00', 24.5554, -81.7800, 0
			 FROM users WHERE email = 'peter@hemingway.com'`,
			nil,
		},
		{
			`INSERT INTO trips (id, user_id, title, description, triptype, status, start_time, end_time, start_lat, start_lon, end_lat, end_lon, distance)
			 SELECT gen_random_uuid(), id, 'Biscayne Bay Morning', 'Sesion corta de pesca en la bahia.', 'FISHING', 'COMPLETED',
			        '2026-03-20 06:30:00', '2026-03-20 11:00:00', 25.5100, -80.2000, 25.6500, -80.2500, 15.2
			 FROM users WHERE email = 'marina@hemingway.com'`,
			nil,
		},
		{
			`INSERT INTO trips (id, user_id, title, description, triptype, status, start_time, end_time, start_lat, start_lon, end_lat, end_lon, distance)
			 SELECT gen_random_uuid(), id, 'Crandon Park Sprint', 'Carrera de velocidad en el parque.', 'RACING', 'COMPLETED',
			        '2026-03-22 08:00:00', '2026-03-22 10:00:00', 25.5100, -80.2000, 25.5300, -80.2200, 12.4
			 FROM users WHERE email = 'marina@hemingway.com'`,
			nil,
		},
		{
			`INSERT INTO trips (id, user_id, title, description, triptype, status, start_time, end_time, start_lat, start_lon, end_lat, end_lon, distance)
			 SELECT gen_random_uuid(), id, 'Night Paddle Under Stars', 'Remada nocturna por la laguna.', 'TOURING', 'COMPLETED',
			        '2026-04-02 20:00:00', '2026-04-02 23:30:00', 25.9500, -80.2800, 25.9700, -80.3000, 9.8
			 FROM users WHERE email = 'lucia@hemingway.com'`,
			nil,
		},
		{
			`INSERT INTO trips (id, user_id, title, description, triptype, status, start_time, start_lat, start_lon)
			 SELECT gen_random_uuid(), id, 'Live: Florida Bay Crossing', 'Travesia en vivo por Florida Bay.', 'FISHING', 'LIVE',
			        '2026-11-16 07:00:00', 25.5250, -80.2050
			 FROM users WHERE email = 'peter@hemingway.com'`,
			nil,
		},

		// ---------- Track de los viajes (12 puntos por viaje) ----------
		{
			`INSERT INTO trip_points (id, trip_id, latitude, longitude, altitude, speed, heading, timestamp)
			 SELECT gen_random_uuid(), t.id,
			        25.7617 + (g * 0.0120), -80.1918 + (g * 0.0090),
			        2 + (g % 5), 6 + (g % 4), (g * 45) % 360,
			        t.start_time + (g * interval '8 minutes')
			 FROM trips t CROSS JOIN generate_series(1, 12) g
			 WHERE t.title = 'Miami Everglades Expedition'`,
			nil,
		},
		{
			`INSERT INTO trip_points (id, trip_id, latitude, longitude, altitude, speed, heading, timestamp)
			 SELECT gen_random_uuid(), t.id,
			        25.5100 + (g * 0.0130), -80.2000 + (g * 0.0045),
			        1 + (g % 3), 7 + (g % 3), (g * 30) % 360,
			        t.start_time + (g * interval '4 minutes')
			 FROM trips t CROSS JOIN generate_series(1, 12) g
			 WHERE t.title = 'Biscayne Bay Morning'`,
			nil,
		},
		{
			`INSERT INTO trip_points (id, trip_id, latitude, longitude, altitude, speed, heading, timestamp)
			 SELECT gen_random_uuid(), t.id,
			        25.5100 + (g * 0.0020), -80.2000 - (g * 0.0015),
			        1, 14 + (g % 6), (g * 15) % 360,
			        t.start_time + (g * interval '1 minute')
			 FROM trips t CROSS JOIN generate_series(1, 12) g
			 WHERE t.title = 'Crandon Park Sprint'`,
			nil,
		},
		{
			`INSERT INTO trip_points (id, trip_id, latitude, longitude, altitude, speed, heading, timestamp)
			 SELECT gen_random_uuid(), t.id,
			        25.9500 + (g * 0.0010), -80.2800 - (g * 0.0012),
			        1, 5 + (g % 3), (g * 20) % 360,
			        t.start_time + (g * interval '18 minutes')
			 FROM trips t CROSS JOIN generate_series(1, 12) g
			 WHERE t.title = 'Night Paddle Under Stars'`,
			nil,
		},
		{
			`INSERT INTO trip_points (id, trip_id, latitude, longitude, altitude, speed, heading, timestamp)
			 SELECT gen_random_uuid(), t.id,
			        25.5250, -80.2050, 1, 8.5, 90, t.start_time
			 FROM trips t WHERE t.title = 'Live: Florida Bay Crossing'`,
			nil,
		},

		// ---------- Torneos ----------
		{
			`INSERT INTO tournaments (id, name, description, organizer_id, start_time, end_time, location_name, start_lat, start_lon, status, tournamenttype, max_participants)
			 SELECT gen_random_uuid(), 'Spring Kingfish Championship 2026', 'El clasico campeonato de kingfish de la costa de Florida.', id,
			        '2026-04-05 07:00:00', '2026-04-05 17:00:00', 'Miami Marina', 25.7617, -80.1918, 'PLANNED', 'FISHING', 50
			 FROM users WHERE email = 'peter@hemingway.com'`,
			nil,
		},
		{
			`INSERT INTO tournaments (id, name, description, organizer_id, start_time, end_time, location_name, start_lat, start_lon, status, tournamenttype, max_participants)
			 SELECT gen_random_uuid(), 'Kayak Racing League - Miami', 'Liga de carreras de kayak por equipos en Crandon Park.', id,
			        '2026-03-22 08:00:00', '2026-03-22 14:00:00', 'Crandon Park', 25.5100, -80.2000, 'LIVE', 'RACING', 30
			 FROM users WHERE email = 'marina@hemingway.com'`,
			nil,
		},
		{
			`INSERT INTO tournaments (id, name, description, organizer_id, start_time, end_time, location_name, start_lat, start_lon, status, tournamenttype, max_participants)
			 SELECT gen_random_uuid(), 'Fall Tarpon Shootout', 'Torneo de otono enfocado en tarpon.', id,
			        '2026-09-12 06:00:00', '2026-09-12 18:00:00', 'Biscayne Inlet', 25.5100, -80.2000, 'FINISHED', 'FISHING', 40
			 FROM users WHERE email = 'lucia@hemingway.com'`,
			nil,
		},

		// ---------- Participantes ----------
		{
			`INSERT INTO tournament_participants (id, tournament_id, user_id, kayak_id, status, current_lat, current_lon, score, catches_count, last_update_time, joined_at)
			 SELECT gen_random_uuid(), t.id, u.id, k.id, 'active', 25.5200, -80.2100, 85.5, 3, now(), '2026-03-20 09:00:00'
			 FROM tournaments t, users u, kayaks k
			 WHERE t.name = 'Kayak Racing League - Miami' AND u.email = 'peter@hemingway.com' AND k.name = 'Explorer Pro'`,
			nil,
		},
		{
			`INSERT INTO tournament_participants (id, tournament_id, user_id, kayak_id, status, current_lat, current_lon, score, catches_count, last_update_time, joined_at)
			 SELECT gen_random_uuid(), t.id, u.id, k.id, 'active', 25.5300, -80.2200, 72.0, 2, now(), '2026-03-20 09:05:00'
			 FROM tournaments t, users u, kayaks k
			 WHERE t.name = 'Kayak Racing League - Miami' AND u.email = 'marina@hemingway.com' AND k.name = 'Sportster'`,
			nil,
		},
		{
			`INSERT INTO tournament_participants (id, tournament_id, user_id, kayak_id, status, current_lat, current_lon, score, catches_count, last_update_time, joined_at)
			 SELECT gen_random_uuid(), t.id, u.id, k.id, 'active', 25.5260, -80.2150, 64.5, 4, now(), '2026-03-20 09:10:00'
			 FROM tournaments t, users u, kayaks k
			 WHERE t.name = 'Kayak Racing League - Miami' AND u.email = 'lucia@hemingway.com' AND k.name = 'Bay Stealth'`,
			nil,
		},
		{
			`INSERT INTO tournament_participants (id, tournament_id, user_id, kayak_id, status, current_lat, current_lon, score, catches_count, last_update_time, joined_at)
			 SELECT gen_random_uuid(), t.id, u.id, k.id, 'active', 25.5180, -80.2080, 55.0, 2, now(), '2026-03-20 09:15:00'
			 FROM tournaments t, users u, kayaks k
			 WHERE t.name = 'Kayak Racing League - Miami' AND u.email = 'diego@hemingway.com' AND k.name = 'Pedal Duke'`,
			nil,
		},
		{
			`INSERT INTO tournament_participants (id, tournament_id, user_id, kayak_id, status, current_lat, current_lon, score, catches_count, last_update_time, joined_at)
			 SELECT gen_random_uuid(), t.id, u.id, k.id, 'active', 25.5340, -80.2250, 41.0, 1, now(), '2026-03-20 09:20:00'
			 FROM tournaments t, users u, kayaks k
			 WHERE t.name = 'Kayak Racing League - Miami' AND u.email = 'ana@hemingway.com' AND k.name = 'Family Tandem'`,
			nil,
		},
		{
			`INSERT INTO tournament_participants (id, tournament_id, user_id, kayak_id, status, score, catches_count, joined_at)
			 SELECT gen_random_uuid(), t.id, u.id, k.id, 'active', 96.0, 5, '2026-09-10 08:00:00'
			 FROM tournaments t, users u, kayaks k
			 WHERE t.name = 'Fall Tarpon Shootout' AND u.email = 'peter@hemingway.com' AND k.name = 'Explorer Pro'`,
			nil,
		},
		{
			`INSERT INTO tournament_participants (id, tournament_id, user_id, kayak_id, status, score, catches_count, joined_at)
			 SELECT gen_random_uuid(), t.id, u.id, k.id, 'left', 30.0, 1, '2026-09-10 08:10:00'
			 FROM tournaments t, users u, kayaks k
			 WHERE t.name = 'Fall Tarpon Shootout' AND u.email = 'carlos@hemingway.com' AND k.name = 'Trawler 12'`,
			nil,
		},

		// ---------- Capturas ----------
		{
			`INSERT INTO tournament_catches (id, tournament_id, participant_id, species, weight, latitude, longitude, photo_url, verified, verified_by, created_at)
			 SELECT gen_random_uuid(), p.tournament_id, p.id, 'King Mackerel', 12.4, 25.5205, -80.2105, NULL, TRUE, 'marina_kayak', '2026-03-22 09:30:00'
			 FROM tournament_participants p JOIN users u ON p.user_id = u.id
			 WHERE u.email = 'peter@hemingway.com' AND p.tournament_id IN (SELECT id FROM tournaments WHERE name = 'Kayak Racing League - Miami')`,
			nil,
		},
		{
			`INSERT INTO tournament_catches (id, tournament_id, participant_id, species, weight, latitude, longitude, photo_url, verified, verified_by, created_at)
			 SELECT gen_random_uuid(), p.tournament_id, p.id, 'Spanish Mackerel', 8.1, 25.5210, -80.2110, NULL, TRUE, 'marina_kayak', '2026-03-22 10:05:00'
			 FROM tournament_participants p JOIN users u ON p.user_id = u.id
			 WHERE u.email = 'peter@hemingway.com' AND p.tournament_id IN (SELECT id FROM tournaments WHERE name = 'Kayak Racing League - Miami')`,
			nil,
		},
		{
			`INSERT INTO tournament_catches (id, tournament_id, participant_id, species, weight, latitude, longitude, photo_url, verified, verified_by, created_at)
			 SELECT gen_random_uuid(), p.tournament_id, p.id, 'Tarpon', 42.0, 25.5305, -80.2205, NULL, FALSE, NULL, '2026-03-22 11:20:00'
			 FROM tournament_participants p JOIN users u ON p.user_id = u.id
			 WHERE u.email = 'marina@hemingway.com' AND p.tournament_id IN (SELECT id FROM tournaments WHERE name = 'Kayak Racing League - Miami')`,
			nil,
		},
		{
			`INSERT INTO tournament_catches (id, tournament_id, participant_id, species, weight, latitude, longitude, photo_url, verified, verified_by, created_at)
			 SELECT gen_random_uuid(), p.tournament_id, p.id, 'Snapper', 5.6, 25.5265, -80.2155, NULL, TRUE, 'peter_kayak', '2026-03-22 11:40:00'
			 FROM tournament_participants p JOIN users u ON p.user_id = u.id
			 WHERE u.email = 'lucia@hemingway.com' AND p.tournament_id IN (SELECT id FROM tournaments WHERE name = 'Kayak Racing League - Miami')`,
			nil,
		},
		{
			`INSERT INTO tournament_catches (id, tournament_id, participant_id, species, weight, latitude, longitude, photo_url, verified, verified_by, created_at)
			 SELECT gen_random_uuid(), p.tournament_id, p.id, 'Barracuda', 9.8, 25.5185, -80.2085, NULL, FALSE, NULL, '2026-03-22 12:10:00'
			 FROM tournament_participants p JOIN users u ON p.user_id = u.id
			 WHERE u.email = 'diego@hemingway.com' AND p.tournament_id IN (SELECT id FROM tournaments WHERE name = 'Kayak Racing League - Miami')`,
			nil,
		},
		{
			`INSERT INTO tournament_catches (id, tournament_id, participant_id, species, weight, latitude, longitude, photo_url, verified, verified_by, created_at)
			 SELECT gen_random_uuid(), p.tournament_id, p.id, 'Tarpon', 58.5, 25.5150, -80.2000, NULL, TRUE, 'lucia_kayak', '2026-09-12 10:00:00'
			 FROM tournament_participants p JOIN users u ON p.user_id = u.id
			 WHERE u.email = 'peter@hemingway.com' AND p.tournament_id IN (SELECT id FROM tournaments WHERE name = 'Fall Tarpon Shootout')`,
			nil,
		},
		{
			`INSERT INTO tournament_catches (id, tournament_id, participant_id, species, weight, latitude, longitude, photo_url, verified, verified_by, created_at)
			 SELECT gen_random_uuid(), p.tournament_id, p.id, 'Snook', 7.2, 25.5155, -80.2005, NULL, TRUE, 'lucia_kayak', '2026-09-12 11:30:00'
			 FROM tournament_participants p JOIN users u ON p.user_id = u.id
			 WHERE u.email = 'peter@hemingway.com' AND p.tournament_id IN (SELECT id FROM tournaments WHERE name = 'Fall Tarpon Shootout')`,
			nil,
		},
		{
			`INSERT INTO tournament_catches (id, tournament_id, participant_id, species, weight, latitude, longitude, photo_url, verified, verified_by, created_at)
			 SELECT gen_random_uuid(), p.tournament_id, p.id, 'Jack Crevalle', 4.4, 25.8505, -80.3505, NULL, FALSE, NULL, '2026-03-15 08:00:00'
			 FROM tournament_participants p JOIN users u ON p.user_id = u.id
			 WHERE u.email = 'carlos@hemingway.com' AND p.tournament_id IN (SELECT id FROM tournaments WHERE name = 'Fall Tarpon Shootout')`,
			nil,
		},

		// ---------- Tracking en vivo ----------
		{
			`INSERT INTO live_tracking (user_id, trip_id, latitude, longitude, speed, heading, altitude, is_active)
			 SELECT u.id, t.id, 25.5250, -80.2050, 8.5, 90, 1, TRUE
			 FROM users u, trips t
			 WHERE u.email = 'peter@hemingway.com' AND t.title = 'Live: Florida Bay Crossing'`,
			nil,
		},
		{
			`INSERT INTO live_tracking (user_id, trip_id, latitude, longitude, speed, heading, altitude, is_active)
			 SELECT u.id, NULL, 25.5350, -80.2150, 6.2, 120, 1, TRUE
			 FROM users u WHERE u.email = 'marina@hemingway.com'`,
			nil,
		},
		{
			`INSERT INTO live_tracking (user_id, trip_id, latitude, longitude, speed, heading, altitude, is_active)
			 SELECT u.id, NULL, 25.9505, -80.2805, 5.0, 200, 1, TRUE
			 FROM users u WHERE u.email = 'lucia@hemingway.com'`,
			nil,
		},
		{
			`INSERT INTO live_tracking (user_id, trip_id, latitude, longitude, speed, is_active)
			 SELECT u.id, t.id, 25.7617, -80.1918, 0, FALSE
			 FROM users u, trips t
			 WHERE u.email = 'carlos@hemingway.com' AND t.title = 'Miami Everglades Expedition'`,
			nil,
		},

		// ---------- Sesiones de tracking ----------
		{
			`INSERT INTO tracking_sessions (id, user_id, trip_id, start_time, start_lat, start_lon, total_distance, is_live, viewer_count)
			 SELECT gen_random_uuid(), u.id, t.id, '2026-11-16 07:00:00', 25.5250, -80.2050, 0, TRUE, 12
			 FROM users u, trips t
			 WHERE u.email = 'peter@hemingway.com' AND t.title = 'Live: Florida Bay Crossing'`,
			nil,
		},
		{
			`INSERT INTO tracking_sessions (id, user_id, trip_id, start_time, end_time, start_lat, start_lon, end_lat, end_lon, total_distance, is_live, viewer_count)
			 SELECT gen_random_uuid(), u.id, t.id, '2026-03-15 06:00:00', '2026-03-15 14:30:00', 25.7617, -80.1918, 25.8900, -80.3600, 28.5, FALSE, 43
			 FROM users u, trips t
			 WHERE u.email = 'peter@hemingway.com' AND t.title = 'Miami Everglades Expedition'`,
			nil,
		},
		{
			`INSERT INTO tracking_sessions (id, user_id, trip_id, start_time, end_time, start_lat, start_lon, end_lat, end_lon, total_distance, is_live, viewer_count)
			 SELECT gen_random_uuid(), u.id, t.id, '2026-03-22 08:00:00', '2026-03-22 10:00:00', 25.5100, -80.2000, 25.5300, -80.2200, 12.4, FALSE, 18
			 FROM users u, trips t
			 WHERE u.email = 'marina@hemingway.com' AND t.title = 'Crandon Park Sprint'`,
			nil,
		},
		{
			`INSERT INTO tracking_sessions (id, user_id, trip_id, start_time, end_time, start_lat, start_lon, end_lat, end_lon, total_distance, is_live, viewer_count)
			 SELECT gen_random_uuid(), u.id, t.id, '2026-03-20 06:30:00', '2026-03-20 11:00:00', 25.5100, -80.2000, 25.6500, -80.2500, 15.2, FALSE, 9
			 FROM users u, trips t
			 WHERE u.email = 'marina@hemingway.com' AND t.title = 'Biscayne Bay Morning'`,
			nil,
		},

		// ---------- Spots de pesca ----------
		{
			`INSERT INTO fishing_spots (id, user_id, name, description, latitude, longitude, rating, fish_caught, is_shared)
			 SELECT gen_random_uuid(), id, 'Secret Point', 'Great for snook at dawn', 25.8500, -80.3500, 5, 47, TRUE
			 FROM users WHERE email = 'peter@hemingway.com'`,
			nil,
		},
		{
			`INSERT INTO fishing_spots (id, user_id, name, description, latitude, longitude, rating, fish_caught, is_shared)
			 SELECT gen_random_uuid(), id, 'North Bay Spots', 'Multiple species, deep water', 25.9500, -80.2800, 4, 32, TRUE
			 FROM users WHERE email = 'peter@hemingway.com'`,
			nil,
		},
		{
			`INSERT INTO fishing_spots (id, user_id, name, description, latitude, longitude, rating, fish_caught, is_shared)
			 SELECT gen_random_uuid(), id, 'Biscayne Inlet', 'Perfect for tarpon', 25.5100, -80.2000, 5, 28, TRUE
			 FROM users WHERE email = 'marina@hemingway.com'`,
			nil,
		},
		{
			`INSERT INTO fishing_spots (id, user_id, name, description, latitude, longitude, rating, fish_caught, is_shared)
			 SELECT gen_random_uuid(), id, 'Crandon Flats', 'Aguas bajas, ideal para fotografia', 25.5300, -80.2200, 4, 15, TRUE
			 FROM users WHERE email = 'lucia@hemingway.com'`,
			nil,
		},
		{
			`INSERT INTO fishing_spots (id, user_id, name, description, latitude, longitude, rating, fish_caught, is_shared)
			 SELECT gen_random_uuid(), id, 'Pedal Cove', 'Cala tranquila para practicar pedal', 25.9700, -80.3000, 3, 8, FALSE
			 FROM users WHERE email = 'diego@hemingway.com'`,
			nil,
		},
		{
			`INSERT INTO fishing_spots (id, user_id, name, description, latitude, longitude, rating, fish_caught, is_shared)
			 SELECT gen_random_uuid(), id, 'Family Cove', 'Zona segura para ninos', 25.8900, -80.3600, 4, 12, FALSE
			 FROM users WHERE email = 'ana@hemingway.com'`,
			nil,
		},

		// ---------- Historial ----------
		{
			`INSERT INTO history_entries (id, user_id, historytype, reference_id, title, description, data, occurred_at)
			 SELECT gen_random_uuid(), u.id, 'trip_created', t.id, 'Viaje creado', t.title,
			         '{"tripId":' || to_json(t.id::text) || '}', t.created_at
			 FROM users u JOIN trips t ON t.user_id = u.id`,
			nil,
		},
		{
			`INSERT INTO history_entries (id, user_id, historytype, reference_id, title, description, data, occurred_at)
			 SELECT gen_random_uuid(), u.id, 'kayak_added', k.id, 'Kayak registrado', k.name,
			         '{"kayakId":' || to_json(k.id::text) || ',"type":"' || coalesce(k.kaytype, '') || '"}', k.created_at
			 FROM users u JOIN kayaks k ON k.user_id = u.id`,
			nil,
		},
		{
			`INSERT INTO history_entries (id, user_id, historytype, reference_id, title, description, data, occurred_at)
			 SELECT gen_random_uuid(), u.id, 'tournament_joined', p.id, 'Se unio a un torneo', t.name,
			         '{"tournamentId":' || to_json(p.tournament_id::text) || '}', p.joined_at
			 FROM users u JOIN tournament_participants p ON p.user_id = u.id JOIN tournaments t ON t.id = p.tournament_id`,
			nil,
		},
		{
			`INSERT INTO history_entries (id, user_id, historytype, reference_id, title, description, data, occurred_at)
			 SELECT gen_random_uuid(), u.id, 'catch_recorded', c.id, 'Captura registrada', c.species,
			         '{"weight":' || c.weight || ',"species":"' || c.species || '"}', c.created_at
			 FROM users u JOIN tournament_participants p ON p.user_id = u.id
			          JOIN tournament_catches c ON c.participant_id = p.id`,
			nil,
		},
		{
			`INSERT INTO history_entries (id, user_id, historytype, reference_id, title, description, data, occurred_at)
			 SELECT gen_random_uuid(), u.id, 'spot_shared', s.id, 'Spot compartido', s.name,
			         '{"spotId":' || to_json(s.id::text) || '}', s.created_at
			 FROM users u JOIN fishing_spots s ON s.user_id = u.id AND s.is_shared = TRUE`,
			nil,
		},
		{
			`INSERT INTO history_entries (id, user_id, historytype, reference_id, title, description, data, occurred_at)
			 SELECT gen_random_uuid(), u.id, 'tracking_started', ss.id, 'Tracking iniciado', 'Sesion en vivo',
			         '{"sessionId":' || to_json(ss.id::text) || '}', ss.start_time
			 FROM users u JOIN tracking_sessions ss ON ss.user_id = u.id`,
			nil,
		},
	}

	for i, seed := range seeds {
		if _, err := db.ExecContext(ctx, seed.query, seed.args...); err != nil {
			return fmt.Errorf("seed %d failed: %w", i+1, err)
		}
	}

	log.Println("Seed data inserted successfully")
	return nil
}

// quoteIdent envuelve un identificador SQL entre comillas dobles, escapando las internas.
func quoteIdent(name string) string {
	return `"` + strings.ReplaceAll(name, `"`, `""`) + `"`
}
