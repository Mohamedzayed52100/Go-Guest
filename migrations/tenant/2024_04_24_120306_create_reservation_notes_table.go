package migrations

import (
	"os"

	"github.com/goplaceapp/goplace-common/pkg/dbhelper"
	"github.com/goplaceapp/goplace-common/pkg/meta"
	"github.com/goplaceapp/goplace-guest/pkg/reservationservice/domain"
	"github.com/jmoiron/sqlx"
	"github.com/tinygg/gofaker"
)

func CreateReservationNotesTable() dbhelper.SqlxMigration {
	return dbhelper.SqlxMigration{
		ID: "2024_04_24_120306_create_reservation_notes_table",
		Migrate: func(tx *sqlx.Tx) error {
			query := `
			CREATE TABLE IF NOT EXISTS reservation_notes (
				id SERIAL PRIMARY KEY,
				creator_id INTEGER NOT NULL,
				description TEXT,
				created_at TIMESTAMPTZ DEFAULT NOW(),
				updated_at TIMESTAMPTZ DEFAULT NOW()
			)
			`
			_, err := tx.Exec(query)
			return err
		},
		Seed: func(tx *sqlx.Tx) error {
			environment := os.Getenv("ENVIRONMENT")
			if environment != meta.ProdEnvironment && environment != meta.StagingEnvironment {
				notes := []domain.ReservationNote{
					{
						CreatorID:   1,
						Description: gofaker.Sentence(10),
					},
					{
						CreatorID:   1,
						Description: gofaker.Sentence(10),
					},
				}

				for _, note := range notes {
					_, err := tx.NamedExec(`INSERT INTO reservation_notes (creator_id, description) VALUES ( :creator_id, :description)`, note)
					if err != nil {
						return err
					}
				}
			}
			return nil
		},
		Rollback: func(tx *sqlx.Tx) error {
			query := `DROP TABLE IF EXISTS reservation_notes`
			_, err := tx.Exec(query)
			return err
		},
	}
}
