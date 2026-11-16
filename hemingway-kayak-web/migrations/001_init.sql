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
