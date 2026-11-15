-- Hemingway Kayak - esquema inicial (MS SQL Server)
-- Se aplica una sola vez; el backend registra el archivo en schema_migrations.

IF NOT EXISTS (SELECT * FROM sys.tables WHERE name = 'users')
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

IF NOT EXISTS (SELECT * FROM sys.tables WHERE name = 'kayaks')
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
);

IF NOT EXISTS (SELECT * FROM sys.tables WHERE name = 'trips')
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
);

IF NOT EXISTS (SELECT * FROM sys.tables WHERE name = 'trip_points')
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

IF NOT EXISTS (SELECT * FROM sys.tables WHERE name = 'tournaments')
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
);

IF NOT EXISTS (SELECT * FROM sys.tables WHERE name = 'tournament_participants')
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

IF NOT EXISTS (SELECT * FROM sys.tables WHERE name = 'tournament_catches')
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

IF NOT EXISTS (SELECT * FROM sys.tables WHERE name = 'tracking_sessions')
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

IF NOT EXISTS (SELECT * FROM sys.tables WHERE name = 'live_tracking')
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

IF NOT EXISTS (SELECT * FROM sys.tables WHERE name = 'history_entries')
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
);

IF NOT EXISTS (SELECT * FROM sys.tables WHERE name = 'fishing_spots')
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
