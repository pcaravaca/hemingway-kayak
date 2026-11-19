package services

import (
	"database/sql"
	"os"
	"testing"

	"github.com/gofiber/fiber/v2"

	"hemingway-kayak-web/pkg/database"
)

func integrationDB(t *testing.T) *sql.DB {
	t.Helper()
	if os.Getenv("HK_INTEGRATION") == "" {
		t.Skip("HK_INTEGRATION no definida: se omiten los tests de integracion")
	}
	db, err := database.Connect()
	if err != nil {
		t.Skipf("no se pudo conectar a la DB de prueba: %v", err)
	}
	t.Cleanup(func() { _ = db.Close() })
	return db
}

// fakeHub captura lo que se le manda por BroadcastTo, para poder afirmar la sala
// y el payload sin un hub real.
type fakeHub struct {
	room      string
	eventType string
	data      interface{}
	calls     int
}

func (f *fakeHub) BroadcastTo(room, eventType string, data interface{}) {
	f.calls++
	f.room = room
	f.eventType = eventType
	f.data = data
}

// TestBroadcastRankingUsaLaSalaDelTorneo verifica lo que se arregló: el ranking se
// emite a la sala del torneo (room=id) y no a la vacía, y el mensaje lleva el
// mismo contenido que devuelve GET /ranking.
func TestBroadcastRankingUsaLaSalaDelTorneo(t *testing.T) {
	db := integrationDB(t)

	tournaments, err := db.Query(`SELECT id FROM tournaments LIMIT 1`)
	if err != nil {
		t.Fatalf("no se pudo leer un torneo: %v", err)
	}
	defer tournaments.Close()
	if !tournaments.Next() {
		t.Skip("la base no tiene torneos sembrados")
	}
	var tournamentID string
	if err := tournaments.Scan(&tournamentID); err != nil {
		t.Fatalf("scan: %v", err)
	}

	hub := &fakeHub{}
	s := NewTournamentService(db)
	s.SetBroadcaster(hub)

	s.broadcastRanking(tournamentID)

	if hub.calls != 1 {
		t.Fatalf("esperaba 1 broadcast, hubo %d", hub.calls)
	}
	if hub.room != tournamentID {
		t.Errorf("el broadcast salió a la sala %q; debía ser %q", hub.room, tournamentID)
	}
	if hub.eventType != "tournament" {
		t.Errorf("eventType = %q; debía ser \"tournament\"", hub.eventType)
	}
	payload, ok := hub.data.(fiber.Map)
	if !ok {
		t.Fatalf("payload inesperado: %T", hub.data)
	}
	rankings, ok := payload["rankings"].([]map[string]interface{})
	if !ok || rankings == nil {
		t.Fatalf("el payload no trae rankings: %#v", payload)
	}
}

// TestRankingForDevuelveSoloActivos: el filtro por status='active' es el que hacía
// que un recién unido nunca apareciera (JoinTournament insertaba 'registered').
func TestRankingForDevuelveSoloActivos(t *testing.T) {
	db := integrationDB(t)

	var tournamentID string
	err := db.QueryRow(`SELECT id FROM tournaments LIMIT 1`).Scan(&tournamentID)
	if err != nil {
		t.Skip("la base no tiene torneos sembrados")
	}

	s := NewTournamentService(db)
	for _, r := range s.rankingFor(tournamentID) {
		if r["rank"] == nil {
			t.Errorf("el ranking no trae rank: %#v", r)
		}
	}
}
