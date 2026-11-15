package models

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID              uuid.UUID `json:"id"`
	Email           string    `json:"email"`
	Username        string    `json:"username"`
	PasswordHash    string    `json:"-"`
	DisplayName     string    `json:"displayName"`
	AvatarURL       string    `json:"avatarUrl,omitempty"`
	Bio             string    `json:"bio,omitempty"`
	TrackingEnabled bool      `json:"trackingEnabled"`
	CreatedAt       time.Time `json:"createdAt"`
	UpdatedAt       time.Time `json:"updatedAt"`
}

type Kayak struct {
	ID          uuid.UUID `json:"id"`
	UserID      uuid.UUID `json:"userId"`
	Name        string    `json:"name"`
	Type        string    `json:"type"`
	Length      float64   `json:"length"`
	Width       float64   `json:"width"`
	Weight      float64   `json:"weight"`
	MaxCapacity float64   `json:"maxCapacity"`
	Material    string    `json:"material"`
	Propulsion  string    `json:"propulsion"`
	ImageURL    string    `json:"imageUrl,omitempty"`
	CreatedAt   time.Time `json:"createdAt"`
}

type Trip struct {
	ID          uuid.UUID  `json:"id"`
	UserID      uuid.UUID  `json:"userId"`
	Title       string     `json:"title"`
	Description string     `json:"description,omitempty"`
	Type        TripType   `json:"type"`
	Status      TripStatus `json:"status"`
	StartTime   time.Time  `json:"startTime"`
	EndTime     *time.Time `json:"endTime,omitempty"`
	StartLat    float64    `json:"startLat"`
	StartLon    float64    `json:"startLon"`
	EndLat      *float64   `json:"endLat,omitempty"`
	EndLon      *float64   `json:"endLon,omitempty"`
	Distance    float64    `json:"distance"`
	CreatedAt   time.Time  `json:"createdAt"`
}

type TripType string

const (
	TripTypeFishing TripType = "fishing"
	TripTypeCamping TripType = "camping"
	TripTypeTour    TripType = "tour"
	TripTypeRacing  TripType = "racing"
)

type TripStatus string

const (
	TripStatusPlanned   TripStatus = "planned"
	TripStatusActive    TripStatus = "active"
	TripStatusCompleted TripStatus = "completed"
	TripStatusCancelled TripStatus = "cancelled"
)

type TripPoint struct {
	ID        uuid.UUID `json:"id"`
	TripID    uuid.UUID `json:"tripId"`
	Latitude  float64   `json:"latitude"`
	Longitude float64   `json:"longitude"`
	Altitude  float64   `json:"altitude,omitempty"`
	Speed     float64   `json:"speed,omitempty"`
	Heading   float64   `json:"heading,omitempty"`
	Timestamp time.Time `json:"timestamp"`
}

type Tournament struct {
	ID              uuid.UUID        `json:"id"`
	Name            string           `json:"name"`
	Description     string           `json:"description,omitempty"`
	OrganizerID     uuid.UUID        `json:"organizerId"`
	StartTime       time.Time        `json:"startTime"`
	EndTime         time.Time        `json:"endTime"`
	LocationName    string           `json:"locationName"`
	StartLat        float64          `json:"startLat"`
	StartLon        float64          `json:"startLon"`
	Status          TournamentStatus `json:"status"`
	Type            TournamentType   `json:"type"`
	MaxParticipants int              `json:"maxParticipants"`
	CreatedAt       time.Time        `json:"createdAt"`
}

type TournamentStatus string

const (
	TournamentStatusPlanned  TournamentStatus = "planned"
	TournamentStatusLive     TournamentStatus = "live"
	TournamentStatusFinished TournamentStatus = "finished"
)

type TournamentType string

const (
	TournamentTypeFishing  TournamentType = "fishing"
	TournamentTypeRacing   TournamentType = "racing"
	TournamentTypeDistance TournamentType = "distance"
)

type TournamentParticipant struct {
	ID             uuid.UUID `json:"id"`
	TournamentID   uuid.UUID `json:"tournamentId"`
	UserID         uuid.UUID `json:"userId"`
	KayakID        uuid.UUID `json:"kayakId"`
	Status         string    `json:"status"`
	CurrentLat     float64   `json:"currentLat,omitempty"`
	CurrentLon     float64   `json:"currentLon,omitempty"`
	Score          float64   `json:"score"`
	CatchesCount   int       `json:"catchesCount"`
	LastUpdateTime time.Time `json:"lastUpdateTime,omitempty"`
	JoinedAt       time.Time `json:"joinedAt"`
}

type TournamentCatch struct {
	ID            uuid.UUID `json:"id"`
	TournamentID  uuid.UUID `json:"tournamentId"`
	ParticipantID uuid.UUID `json:"participantId"`
	Species       string    `json:"species"`
	Weight        float64   `json:"weight"`
	Latitude      float64   `json:"latitude"`
	Longitude     float64   `json:"longitude"`
	PhotoURL      string    `json:"photoUrl,omitempty"`
	Verified      bool      `json:"verified"`
	VerifiedBy    *string   `json:"verifiedBy,omitempty"`
	CreatedAt     time.Time `json:"createdAt"`
}

type LiveTracking struct {
	UserID    uuid.UUID `json:"userId"`
	TripID    uuid.UUID `json:"tripId,omitempty"`
	Latitude  float64   `json:"latitude"`
	Longitude float64   `json:"longitude"`
	Speed     float64   `json:"speed,omitempty"`
	Heading   float64   `json:"heading,omitempty"`
	Altitude  float64   `json:"altitude,omitempty"`
	Timestamp time.Time `json:"timestamp"`
	IsActive  bool      `json:"isActive"`
}

type TrackingSession struct {
	ID            uuid.UUID  `json:"id"`
	UserID        uuid.UUID  `json:"userId"`
	TripID        uuid.UUID  `json:"tripId,omitempty"`
	StartTime     time.Time  `json:"startTime"`
	EndTime       *time.Time `json:"endTime,omitempty"`
	StartLat      float64    `json:"startLat"`
	StartLon      float64    `json:"startLon"`
	EndLat        *float64   `json:"endLat,omitempty"`
	EndLon        *float64   `json:"endLon,omitempty"`
	TotalDistance float64    `json:"totalDistance"`
	IsLive        bool       `json:"isLive"`
	ViewerCount   int        `json:"viewerCount"`
}

type HistoryEntry struct {
	ID          uuid.UUID   `json:"id"`
	UserID      uuid.UUID   `json:"userId"`
	Type        HistoryType `json:"type"`
	ReferenceID uuid.UUID   `json:"referenceId"`
	Title       string      `json:"title"`
	Description string      `json:"description,omitempty"`
	Data        string      `json:"data,omitempty"`
	OccurredAt  time.Time   `json:"occurredAt"`
	CreatedAt   time.Time   `json:"createdAt"`
}

type HistoryType string

const (
	HistoryTypeTrip       HistoryType = "trip"
	HistoryTypeTournament HistoryType = "tournament"
	HistoryTypeCatch      HistoryType = "catch"
	HistoryTypeMilestone  HistoryType = "milestone"
)
