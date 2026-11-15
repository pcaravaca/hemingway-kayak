# Hemingway Kayak Web - Agent Guidelines

## Project Overview

Web companion app for Hemingway Kayak Android app.

### Features
- **Live Tracking**: Follow kayakers in real-time (with permission)
- **Tournaments**: Live rankings, locations, catches
- **History**: View past trips, tournaments, catches
- **Trip Planning**: Plan fishing trips

### Tech Stack
- **Backend**: Go 1.21+ with Fiber v2 framework
- **Frontend**: Vue 3 with Composition API, Pinia, Vue Router
- **Database**: MS SQL Server 2022
- **Maps**: Leaflet (OSM)
- **WebSocket**: Fiber WebSocket for real-time updates

## Project Structure

```
hemingway-kayak-web/
├── cmd/server/           # Go entrypoint
│   └── main.go          # Server startup, migrations, seed
├── internal/
│   ├── api/             # HTTP handlers
│   │   └── server.go    # Fiber server setup, routes
│   ├── models/         # Data models
│   │   └── models.go
│   ├── services/        # Business logic
│   │   ├── user_service.go
│   │   ├── trip_service.go
│   │   ├── tournament_service.go
│   │   ├── tracking_service.go
│   │   └── history_service.go
│   └── ws/             # WebSocket hub
│       └── hub.go
├── pkg/
│   └── database/       # DB connections & migrations
│       └── database.go # MS SQL, migrations, seed data
├── web/                # Vue 3 frontend
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

## Build Commands

### Backend

```bash
cd /mnt/c/Users/Peter/Documents/Main\ Repo/hemingway-kayak-web

# Download dependencies
go mod download

# Verify compilation
go build ./...

# Run server (port 8080)
go run cmd/server/main.go
```

### Frontend

```bash
cd /mnt/c/Users/Peter/Documents/Main\ Repo/hemingway-kayak-web/web

# Install dependencies
npm install

# Run dev server (port 5173)
npm run dev

# Production build
npm run build
```

### Full Stack

```bash
# Terminal 1 - Backend
cd /mnt/c/Users/Peter/Documents/Main\ Repo/hemingway-kayak-web
go run cmd/server/main.go  # Server on :8080

# Terminal 2 - Frontend
cd /mnt/c/Users/Peter/Documents/Main\ Repo/hemingway-kayak-web/web
npm run dev  # Frontend on :5173
```

## Database

### MS SQL Server Configuration
- Host: `10.10.1.21`
- Port: `1433`
- Database: `Hemingway`
- User: `sa`
- Password: `[REDACTED]`

### Database Tables (Auto-created via migrations)

```sql
-- Users
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
);

-- Kayaks (type column renamed to kaytype - 'type' is reserved)
CREATE TABLE kayaks (
    id UNIQUEIDENTIFIER PRIMARY KEY DEFAULT NEWID(),
    user_id UNIQUEIDENTIFIER REFERENCES users(id),
    name NVARCHAR(200) NOT NULL,
    kaytype NVARCHAR(50),
    length FLOAT,
    width FLOAT,
    weight FLOAT,
    max_capacity FLOAT,
    material NVARCHAR(50),
    propulsion NVARCHAR(50),
    image_url NVARCHAR(500),
    created_at DATETIME2 DEFAULT GETDATE()
);

-- Trips (type column renamed to triptype)
CREATE TABLE trips (
    id UNIQUEIDENTIFIER PRIMARY KEY DEFAULT NEWID(),
    user_id UNIQUEIDENTIFIER REFERENCES users(id),
    title NVARCHAR(300) NOT NULL,
    description NVARCHAR(MAX),
    triptype NVARCHAR(50) NOT NULL,
    status NVARCHAR(50) DEFAULT 'planned',
    start_time DATETIME2,
    end_time DATETIME2,
    start_lat FLOAT,
    start_lon FLOAT,
    end_lat FLOAT,
    end_lon FLOAT,
    distance FLOAT DEFAULT 0,
    created_at DATETIME2 DEFAULT GETDATE()
);

-- Trip Points
CREATE TABLE trip_points (
    id UNIQUEIDENTIFIER PRIMARY KEY DEFAULT NEWID(),
    trip_id UNIQUEIDENTIFIER REFERENCES trips(id),
    latitude FLOAT NOT NULL,
    longitude FLOAT NOT NULL,
    altitude FLOAT,
    speed FLOAT,
    heading FLOAT,
    timestamp DATETIME2 NOT NULL
);

-- Tournaments (type column renamed to tournamenttype)
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
    tournamenttype NVARCHAR(50),
    max_participants INT,
    created_at DATETIME2 DEFAULT GETDATE()
);

-- Tournament Participants
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
);

-- Tournament Catches
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
);

-- Tracking Sessions
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
);

-- Live Tracking
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
);

-- History Entries (type renamed to historytype)
CREATE TABLE history_entries (
    id UNIQUEIDENTIFIER PRIMARY KEY DEFAULT NEWID(),
    user_id UNIQUEIDENTIFIER REFERENCES users(id),
    historytype NVARCHAR(50) NOT NULL,
    reference_id UNIQUEIDENTIFIER,
    title NVARCHAR(500) NOT NULL,
    description NVARCHAR(MAX),
    data NVARCHAR(MAX),
    occurred_at DATETIME2 NOT NULL,
    created_at DATETIME2 DEFAULT GETDATE()
);

-- Fishing Spots
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

**Password for all users:** `password123`

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

## Backend Dependencies

| Library | Version | Purpose |
|---------|---------|---------|
| github.com/gofiber/fiber/v2 | v2.52.0 | HTTP framework |
| github.com/microsoft/go-mssqldb | v1.6.0 | MS SQL driver |
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
DATABASE_HOST=10.10.1.21
DATABASE_PORT=1433
DATABASE_NAME=Hemingway
DATABASE_USER=sa
DATABASE_PASSWORD=[REDACTED]
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

### Go Conventions
- Use `camelCase` for function names
- Use `PascalCase` for type names
- Error handling: return `nil` on success, `error` on failure
- Context passed as first parameter: `func DoThing(ctx context.Context, ...)`
- **MS SQL**: UNIQUEIDENTIFIER columns must use `CAST(col AS VARCHAR(36))` in SELECT statements

### Vue Conventions
- Use `<script setup>` syntax
- Use Composition API with `ref` and `computed`
- Store state in Pinia stores
- Follow Vue style guide for naming

### Important: Reserved SQL Keywords

MS SQL reserved words like `type` cannot be used as column names. The following renames are used:
- `kayaks.type` → `kayaks.kaytype`
- `trips.type` → `trips.triptype`
- `tournaments.type` → `tournaments.tournamenttype`
- `history_entries.type` → `history_entries.historytype`

## Testing

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
