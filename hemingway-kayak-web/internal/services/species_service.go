package services

import (
	"context"
	"database/sql"
	"encoding/json"

	"github.com/gofiber/fiber/v2"
)

// SpeciesService sirve la DB de peces de Costa Rica (migración 004_species.sql).
type SpeciesService struct {
	db *sql.DB
}

func NewSpeciesService(db *sql.DB) *SpeciesService {
	return &SpeciesService{db: db}
}

const speciesSelect = `SELECT CAST(id AS VARCHAR(36)), common_cr, scientific, family, taxonomic_order,
	       habitat, size_cm, gear, edible, invasive, lookalikes, distinguishing, incopesca
	FROM species`

func (s *SpeciesService) scanSpecies(rows *sql.Rows) ([]map[string]interface{}, error) {
	out := []map[string]interface{}{}
	for rows.Next() {
		var id, commonCR, scientific, family, orderName, habitat, sizeCm, gear sql.NullString
		var edible, invasive bool
		var lookalikes, distinguishing, incopesca []byte

		if err := rows.Scan(&id, &commonCR, &scientific, &family, &orderName, &habitat,
			&sizeCm, &gear, &edible, &invasive, &lookalikes, &distinguishing, &incopesca); err != nil {
			continue
		}

		likes := []string{}
		if len(lookalikes) > 0 {
			_ = json.Unmarshal(lookalikes, &likes)
		}

		out = append(out, map[string]interface{}{
			"id":             id.String,
			"commonName":     commonCR.String,
			"scientific":     scientific.String,
			"family":         family.String,
			"order":          orderName.String,
			"habitat":        habitat.String,
			"sizeCm":         sizeCm.String,
			"gear":           gear.String,
			"edible":         edible,
			"invasive":       invasive,
			"lookalikes":     likes,
			"distinguishing": string(distinguishing),
			"incopesca":      string(incopesca),
		})
	}
	return out, rows.Err()
}

// GetSpecies devuelve todas las especies, opcionalmente filtradas por habitat
// (rio | estuario | mar | manglar). Público: lo consume el formulario de capturas.
func (s *SpeciesService) GetSpecies(c *fiber.Ctx) error {
	query := speciesSelect
	args := []interface{}{}
	if habitat := c.Query("habitat"); habitat != "" {
		query += " WHERE habitat = $1"
		args = append(args, habitat)
	}
	query += " ORDER BY habitat, common_cr"

	rows, err := s.db.QueryContext(context.Background(), query, args...)
	if err != nil {
		return c.Status(500).JSON(Error("Database error"))
	}
	defer rows.Close()

	out, err := s.scanSpecies(rows)
	if err != nil {
		return c.Status(500).JSON(Error("Database error"))
	}
	return c.JSON(out)
}

// GetSpeciesByID devuelve una especie por id.
func (s *SpeciesService) GetSpeciesByID(c *fiber.Ctx) error {
	rows, err := s.db.QueryContext(context.Background(), speciesSelect+" WHERE id = $1", c.Params("id"))
	if err != nil {
		return c.Status(500).JSON(Error("Database error"))
	}
	defer rows.Close()

	out, err := s.scanSpecies(rows)
	if err != nil {
		return c.Status(500).JSON(Error("Database error"))
	}
	if len(out) == 0 {
		return c.Status(404).JSON(Error("Species not found"))
	}
	return c.JSON(out[0])
}
