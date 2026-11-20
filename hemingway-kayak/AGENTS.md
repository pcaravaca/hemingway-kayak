# Hemingway Kayak - Agent Guidelines

## Project Overview

Hemingway Kayak is a full-stack solution for kayak fishing enthusiasts:

| Component | Technology | Location |
|-----------|------------|----------|
| **Android App** | Kotlin, MVVM, OSMDroid | `/hemingway-kayak/` |
| **Web Backend** | Go, Fiber, PostgreSQL | `/hemingway-kayak-web/` |
| **Web Frontend** | Vue 3, Pinia, Leaflet | `/hemingway-kayak-web/web/` |

### Android App Features
- Real-time weather & barometric pressure trends
- Moon phase tracking
- Fishing recommendations
- Live tracking (permission-based) con salida propia, distancia y mapa
- Spots compartidos (listar y compartir)
- Tournament management with live rankings
- Captura de torneo con foto
- Trip history & planning
- Emergency notifications with SMS (SOS + monitor de estabilidad)
- Comandos de voz (tracking, SOS, marcar punto, distancia/velocidad)
- Sonar: **no implementado** (requiere SDKs propietarios)
- Suscripciones: **no implementado** (requiere Play Billing)

### Web App Features
- Live tracking of kayakers (with permission)
- Tournament management with live rankings
- Trip history viewing
- Trip planning

## Build Commands

### Android App

```bash
cd /mnt/c/Users/Peter/Documents/Main\ Repo/hemingway-kayak

# Build
./gradlew assembleDebug          # Debug build
./gradlew assembleRelease        # Release build
./gradlew build                  # Full build with checks

# Testing
./gradlew test                   # Run unit tests
./gradlew test --tests "com.hemingway.kayak.data.*"  # Run specific test package

# Code analysis
./gradlew lint                   # Run lint
./gradlew ktlintCheck           # Kotlin style check

# Clean
./gradlew clean
```

### Web Backend (Go)

```bash
cd /mnt/c/Users/Peter/Documents/Main\ Repo/hemingway-kayak-web

# Build & Run
go build ./...                    # Verify compilation
go mod download                  # Download dependencies
go run cmd/server/main.go        # Start server on :8080
```

### Web Frontend (Vue)

```bash
cd /mnt/c/Users/Peter/Documents/Main\ Repo/hemingway-kayak-web/web

# Install & Run
npm install                       # Install dependencies
npm run dev                       # Dev server on :5173
npm run build                     # Production build
```

### Full Stack

```bash
# Terminal 1 - Backend
cd /mnt/c/Users/Peter/Documents/Main\ Repo/hemingway-kayak-web
go run cmd/server/main.go        # Server on :8080

# Terminal 2 - Frontend
cd /mnt/c/Users/Peter/Documents/Main\ Repo/hemingway-kayak-web/web
npm run dev                      # Frontend on :5173
```

## Architecture

### Android App Structure

```
hemingway-kayak/
├── app/src/main/java/com/hemingway/kayak/
│   ├── data/           # Repository, Room DB, API services, models
│   ├── di/             # Hilt dependency injection modules
│   ├── ui/             # Fragments, ViewModels, MVVM
│   ├── location/       # GPS tracking
│   ├── sensors/        # Device sensors
│   ├── emergency/      # Emergency notifications
│   ├── voice/          # Voice commands
│   ├── sonar/          # Sonar integration
│   ├── places/         # Places API
│   ├── receivers/      # Broadcast receivers
│   └── utils/          # Utilities
├── app/build/outputs/apk/debug/app-debug.apk  # Built APK
└── AGENTS.md           # This file
```

### Web App Structure

```
hemingway-kayak-web/
├── cmd/server/           # Go entrypoint
│   └── main.go          # Server startup, migrations, seed
├── internal/
│   ├── api/             # HTTP handlers
│   │   └── server.go    # Fiber server setup, routes
│   ├── models/          # Data models
│   │   └── models.go
│   └── services/        # Business logic
│       ├── user_service.go
│       ├── trip_service.go
│       ├── tournament_service.go
│       ├── tracking_service.go
│       └── history_service.go
├── pkg/
│   └── database/        # DB connections & migrations
│       └── database.go  # PostgreSQL, migrations, seed data
├── web/                 # Vue 3 frontend
│   ├── src/
│   │   ├── assets/
│   │   │   └── hemingway.css  # Shared design system
│   │   ├── components/
│   │   ├── views/
│   │   │   ├── HomeView.vue
│   │   │   ├── TrackingView.vue
│   │   │   ├── TournamentsView.vue
│   │   │   ├── TournamentDetailView.vue
│   │   │   ├── HistoryView.vue
│   │   │   ├── LoginView.vue
│   │   │   └── ProfileView.vue
│   │   ├── stores/
│   │   │   └── auth.js       # Pinia auth store
│   │   ├── composables/
│   │   │   └── useTrackingWebSocket.js
│   │   ├── router/
│   │   │   └── index.js
│   │   ├── App.vue
│   │   └── main.js
│   └── public/
├── go.mod
└── AGENTS.md
```

## Database

### PostgreSQL Configuration
- Host: `100.89.167.67` (servidor Ubuntu, contenedor `fcfeed-pgdev`)
- Port: `5499`
- Database: `hemingway_dev`
- User: `hemingway`
- Password: en `.env` (`DB_PASSWORD`). Si `DATABASE_URL` está definida, tiene prioridad.
- Driver: `github.com/jackc/pgx/v5/stdlib` (nombre de driver `pgx`)

### Database Tables (Auto-created via migrations)

```sql
-- Hemingway Kayak - esquema inicial (PostgreSQL)
-- Se aplica una sola vez; el backend registra el archivo en schema_migrations.

CREATE TABLE IF NOT EXISTS users (
	id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
	email TEXT UNIQUE NOT NULL,
	username TEXT UNIQUE NOT NULL,
	password_hash TEXT NOT NULL,
	display_name TEXT,
	avatar_url TEXT,
	bio TEXT,
	tracking_enabled BOOLEAN DEFAULT FALSE,
	created_at TIMESTAMPTZ DEFAULT now(),
	updated_at TIMESTAMPTZ DEFAULT now()
);

CREATE TABLE IF NOT EXISTS kayaks (
	id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
	user_id UUID REFERENCES users(id),
	name TEXT NOT NULL,
	kaytype VARCHAR(50),
	length DOUBLE PRECISION,
	width DOUBLE PRECISION,
	weight DOUBLE PRECISION,
	max_capacity DOUBLE PRECISION,
	material VARCHAR(50),
	propulsion VARCHAR(50),
	image_url TEXT,
	created_at TIMESTAMPTZ DEFAULT now()
);

CREATE TABLE IF NOT EXISTS trips (
	id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
	user_id UUID REFERENCES users(id),
	title TEXT NOT NULL,
	description TEXT,
	triptype VARCHAR(50) NOT NULL,
	status VARCHAR(50) DEFAULT 'planned',
	start_time TIMESTAMPTZ,
	end_time TIMESTAMPTZ,
	start_lat DOUBLE PRECISION,
	start_lon DOUBLE PRECISION,
	end_lat DOUBLE PRECISION,
	end_lon DOUBLE PRECISION,
	distance DOUBLE PRECISION DEFAULT 0,
	created_at TIMESTAMPTZ DEFAULT now()
);

CREATE TABLE IF NOT EXISTS trip_points (
	id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
	trip_id UUID REFERENCES trips(id),
	latitude DOUBLE PRECISION NOT NULL,
	longitude DOUBLE PRECISION NOT NULL,
	altitude DOUBLE PRECISION,
	speed DOUBLE PRECISION,
	heading DOUBLE PRECISION,
	timestamp TIMESTAMPTZ NOT NULL
);

CREATE TABLE IF NOT EXISTS tournaments (
	id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
	name TEXT NOT NULL,
	description TEXT,
	organizer_id UUID REFERENCES users(id),
	start_time TIMESTAMPTZ NOT NULL,
	end_time TIMESTAMPTZ NOT NULL,
	location_name TEXT,
	start_lat DOUBLE PRECISION,
	start_lon DOUBLE PRECISION,
	status VARCHAR(50) DEFAULT 'planned',
	tournamenttype VARCHAR(50),
	max_participants INTEGER,
	created_at TIMESTAMPTZ DEFAULT now()
);

CREATE TABLE IF NOT EXISTS tournament_participants (
	id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
	tournament_id UUID REFERENCES tournaments(id),
	user_id UUID REFERENCES users(id),
	kayak_id UUID REFERENCES kayaks(id),
	status VARCHAR(50) DEFAULT 'registered',
	current_lat DOUBLE PRECISION,
	current_lon DOUBLE PRECISION,
	score DOUBLE PRECISION DEFAULT 0,
	catches_count INTEGER DEFAULT 0,
	last_update_time TIMESTAMPTZ,
	joined_at TIMESTAMPTZ DEFAULT now()
);

CREATE TABLE IF NOT EXISTS tournament_catches (
	id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
	tournament_id UUID REFERENCES tournaments(id),
	participant_id UUID REFERENCES tournament_participants(id),
	species TEXT NOT NULL,
	weight DOUBLE PRECISION NOT NULL,
	latitude DOUBLE PRECISION,
	longitude DOUBLE PRECISION,
	photo_url TEXT,
	verified BOOLEAN DEFAULT FALSE,
	verified_by VARCHAR(200),
	created_at TIMESTAMPTZ DEFAULT now()
);

CREATE TABLE IF NOT EXISTS tracking_sessions (
	id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
	user_id UUID REFERENCES users(id),
	trip_id UUID REFERENCES trips(id),
	start_time TIMESTAMPTZ DEFAULT now(),
	end_time TIMESTAMPTZ,
	start_lat DOUBLE PRECISION,
	start_lon DOUBLE PRECISION,
	end_lat DOUBLE PRECISION,
	end_lon DOUBLE PRECISION,
	total_distance DOUBLE PRECISION DEFAULT 0,
	is_live BOOLEAN DEFAULT FALSE,
	viewer_count INTEGER DEFAULT 0
);

CREATE TABLE IF NOT EXISTS live_tracking (
	user_id UUID PRIMARY KEY,
	trip_id UUID REFERENCES trips(id),
	latitude DOUBLE PRECISION NOT NULL,
	longitude DOUBLE PRECISION NOT NULL,
	speed DOUBLE PRECISION,
	heading DOUBLE PRECISION,
	altitude DOUBLE PRECISION,
	timestamp TIMESTAMPTZ DEFAULT now(),
	is_active BOOLEAN DEFAULT TRUE
);

CREATE TABLE IF NOT EXISTS history_entries (
	id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
	user_id UUID REFERENCES users(id),
	historytype VARCHAR(50) NOT NULL,
	reference_id UUID,
	title TEXT NOT NULL,
	description TEXT,
	data TEXT,
	occurred_at TIMESTAMPTZ NOT NULL,
	created_at TIMESTAMPTZ DEFAULT now()
);

CREATE TABLE IF NOT EXISTS fishing_spots (
	id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
	user_id UUID REFERENCES users(id),
	name TEXT NOT NULL,
	description TEXT,
	latitude DOUBLE PRECISION NOT NULL,
	longitude DOUBLE PRECISION NOT NULL,
	rating INTEGER DEFAULT 0,
	fish_caught INTEGER DEFAULT 0,
	is_shared BOOLEAN DEFAULT FALSE,
	created_at TIMESTAMPTZ DEFAULT now()
);
```

### Seed Data

The server automatically seeds test data on first run:

**Users:**
| Email | Username | Display Name |
|-------|----------|--------------|
| peter@hemingway.com | peter_kayak | Peter Garcia |
| marina@hemingway.com | marina_kayak | Marina Rodriguez |
| carlos@hemingway.com | carlos_kayak | Carlos Mendoza |

**Password for all users:** la define la semilla y se hashea con bcrypt (no es `password123`). En el ambiente de desarrollo la cuenta de Peter es `peter@hemingway.com`.

**Sample Data:**
- 3 kayaks (Explorer Pro, Sportster, Trawler 12)
- 3 trips (Miami Everglades, Keys Weekend, Biscayne Bay)
- 2 tournaments (Spring Kingfish Championship, Kayak Racing League)
- 2 active tracking sessions
- 3 fishing spots

## API Endpoints

### Authentication
| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/api/v1/auth/register` | Register new user |
| POST | `/api/v1/auth/login` | Login user |
| POST | `/api/v1/auth/logout` | Logout user |

### Users
| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/api/v1/users/:id` | Get user profile |
| PUT | `/api/v1/users/:id` | Update user |
| GET | `/api/v1/users/:id/profile` | Get user profile details |

### Kayaks
| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/api/v1/kayaks?userId=` | List user's kayaks |
| POST | `/api/v1/kayaks` | Create kayak |
| PUT | `/api/v1/kayaks/:id` | Update kayak |
| DELETE | `/api/v1/kayaks/:id` | Delete kayak |

### Trips
| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/api/v1/trips` | List trips |
| GET | `/api/v1/trips/:id` | Get trip details |
| POST | `/api/v1/trips` | Create trip |
| PUT | `/api/v1/trips/:id` | Update trip |
| DELETE | `/api/v1/trips/:id` | Delete trip |
| GET | `/api/v1/trips/:id/points` | Get trip GPS points |
| POST | `/api/v1/trips/:id/points` | Add trip point |

### Tournaments
| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/api/v1/tournaments` | List tournaments |
| GET | `/api/v1/tournaments/live` | Live tournaments |
| GET | `/api/v1/tournaments/:id` | Tournament details |
| POST | `/api/v1/tournaments` | Create tournament |
| PUT | `/api/v1/tournaments/:id` | Update tournament |
| GET | `/api/v1/tournaments/:id/participants` | Get participants |
| GET | `/api/v1/tournaments/:id/ranking` | Get rankings |
| POST | `/api/v1/tournaments/:id/join` | Join tournament |
| POST | `/api/v1/tournaments/:id/leave` | Leave tournament |
| POST | `/api/v1/tournaments/:id/catches` | Record catch |
| GET | `/api/v1/tournaments/:id/catches` | Get catches |

### Live Tracking
| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/api/v1/tracking/session/:userId` | Get active session |
| POST | `/api/v1/tracking/start` | Start tracking |
| POST | `/api/v1/tracking/stop` | Stop tracking |
| POST | `/api/v1/tracking/point` | Update location |
| GET | `/api/v1/tracking/follow/:userId` | Follow user |
| GET | `/api/v1/tracking/followers` | Get all followers |
| WS | `/ws/tracking` | WebSocket for live updates |

### History
| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/api/v1/history/user/:userId` | Get user history |
| GET | `/api/v1/history/user/:userId/trips` | Get trip history |
| GET | `/api/v1/history/user/:userId/tournaments` | Get tournament history |

### Fishing Spots
| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/api/v1/spots` | List shared spots |
| POST | `/api/v1/spots` | Create spot |
| PUT | `/api/v1/spots/:id` | Update spot |
| DELETE | `/api/v1/spots/:id` | Delete spot |

### Health
| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/api/v1/health` | Health check |

### WebSocket Messages

```json
// Client -> Server
{"action": "subscribe", "tournamentId": "uuid"}
{"action": "unsubscribe", "tournamentId": "uuid"}

// Server -> Client
{"type": "tracking", "data": {"userId": "...", "latitude": 40.7, "longitude": -74.0}}
{"type": "tournament", "data": {"rankings": [...]}}
```

## Design System

### Color Palette

| Color | Hex | CSS Variable | Usage |
|-------|-----|--------------|-------|
| Neon | `#4DFFA9` | `--hemmingway-neon` | Primary actions, highlights |
| Turquoise | `#40E0D0` | `--hemmingway-turquoise` | Secondary elements |
| Dark | `#1A1B1E` | `--hemmingway-dark` | Backgrounds |
| Background Secondary | `#2A2B2E` | `--hemmingway-bg-secondary` | Cards, panels |

### CSS Variables Available
```css
--hemmingway-neon: #4DFFA9;
--hemmingway-turquoise: #40E0D0;
--hemmingway-dark: #1A1B1E;
--hemmingway-bg-primary: #1A1B1E;
--hemmingway-bg-secondary: #2A2B2E;
--hemmingway-bg-tertiary: #3A3B3E;
--hemmingway-text-primary: #FFFFFF;
--hemmingway-text-secondary: #B0B0B0;
--hemmingway-text-muted: #707070;
--hemmingway-accent-success: #4DFFA9;
--hemmingway-accent-warning: #FFD700;
--hemmingway-accent-error: #FF6B6B;
--hemmingway-accent-info: #40E0D0;
--hemmingway-gradient: linear-gradient(135deg, #4DFFA9 0%, #40E0D0 100%);
```

### Usage in Vue Components
```vue
<style scoped>
.primary-btn {
  background: var(--hemmingway-neon);
  color: var(--hemmingway-dark);
}
.card {
  background: var(--hemmingway-bg-secondary);
  border-radius: var(--hemmingway-radius-xl);
}
</style>
```

## Android Dependencies

| Library | Version | Purpose |
|---------|---------|---------|
| Kotlin | 1.8.0 | Language |
| Hilt | 2.46.1 | Dependency Injection |
| Room | 2.5.2 | Local Database |
| Retrofit | 2.9.0 | HTTP Client |
| Moshi | 1.14.0 | JSON Parsing |
| OSMDroid | 6.8.0 | Maps (no API key) |
| Navigation | 2.7.5 | Navigation Component |
| WorkManager | 2.8.1 | Background Tasks |
| Coroutines | 1.7.3 | Async Operations |

## Backend Dependencies

| Library | Version | Purpose |
|---------|---------|---------|
| github.com/gofiber/fiber/v2 | v2.52.0 | HTTP framework |
| github.com/jackc/pgx/v5 | v5.7.6 | PostgreSQL driver |
| github.com/google/uuid | v1.6.0 | UUID generation |
| github.com/golang-jwt/jwt/v5 | v5.2.0 | JWT authentication |
| golang.org/x/crypto/bcrypt | v0.21.0 | Password hashing |
| github.com/redis/go-redis/v9 | v9.4.0 | Redis client (optional) |

## Frontend Dependencies

| Library | Version | Purpose |
|---------|---------|---------|
| vue | ^3.4.0 | Framework |
| pinia | ^2.1.0 | State management |
| vue-router | ^4.2.0 | Routing |
| axios | ^1.6.0 | HTTP client |
| leaflet | ^1.9.0 | Maps |
| @vue-leaflet/vue-leaflet | ^0.10.0 | Vue Leaflet wrapper |

## Environment Variables

### Backend
```bash
PORT=8080
DB_SERVER=100.89.167.67
DB_PORT=5499
DB_NAME=hemingway_dev
DB_USER=hemingway
DB_PASSWORD=
# DATABASE_URL=postgres://hemingway:PASSWORD@100.89.167.67:5499/hemingway_dev?sslmode=disable
CORS_ORIGINS=http://localhost:5173
JWT_SECRET=your-secret-key-here
REDIS_ADDR=localhost:6379  # Optional
```

### Frontend
```bash
VITE_API_URL=http://localhost:8080
VITE_WS_URL=ws://localhost:8080
```

## Android Integration

The Android app can send tracking data to this backend:

```kotlin
// When user enables live tracking
POST /api/v1/tracking/start
Body: {"userId": "...", "tripId": "...", "latitude": 40.7, "longitude": -74.0}

// Periodic location updates
POST /api/v1/tracking/point
Body: {"userId": "...", "latitude": 40.8, "longitude": -74.1, "speed": 5.0}
```

## Code Style

### Kotlin Conventions
- Use `val` by default; `var` only when mutability needed
- Nullable types: `T?` with proper null handling
- Use `when` instead of `if-else` chains
- Use `data class` for models, `object` for singletons
- Enums for constrained sets: `KayakType`, `KayakMaterial`, `PropulsionType`

### Go Conventions
- Use `camelCase` for function names
- Use `PascalCase` for type names
- Error handling: return `nil` on success, `error` on failure
- Context passed as first parameter: `func DoThing(ctx context.Context, ...)`
- PostgreSQL: las columnas `UUID` se escanean directo a `string` con `pgx`, sin `CAST`.

### Vue Conventions
- Use `<script setup>` syntax
- Use Composition API with `ref` and `computed`
- Store state in Pinia stores
- Follow Vue style guide for naming

## Testing

### Android Tests
```bash
./gradlew test --tests "com.hemingway.kayak.data.WeatherRepositoryTest"
```

### API Tests
```bash
# Health check
curl http://localhost:8080/api/v1/health

# Get all tournaments
curl http://localhost:8080/api/v1/tournaments

# Get live tournaments
curl http://localhost:8080/api/v1/tournaments/live

# Get all trips
curl http://localhost:8080/api/v1/trips

# Get active trackers
curl http://localhost:8080/api/v1/tracking/followers
```

## TODO

- [ ] Redis for WebSocket pub/sub scaling
- [ ] JWT middleware for protected routes
- [ ] File upload for catch photos
- [ ] Push notifications
- [ ] Real-time tournament updates via WebSocket
- [ ] Tournament catch submission API
- [ ] Photo upload for catches
- [ ] Connect Android app to web backend
- [ ] Real password hashing verification
