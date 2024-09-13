package migrations

import (
	"os"

	"github.com/goplaceapp/goplace-common/pkg/dbhelper"
	"github.com/goplaceapp/goplace-common/pkg/meta"
	"github.com/goplaceapp/goplace-guest/pkg/guestservice/domain"
	"github.com/jmoiron/sqlx"
	"github.com/tinygg/gofaker"
)

func CreateGuestNotesTable() dbhelper.SqlxMigration {
	return dbhelper.SqlxMigration{
		ID: "2024_04_24_120209_create_guest_notes_table",
		Migrate: func(tx *sqlx.Tx) error {
			query := `
			CREATE TABLE IF NOT EXISTS guest_notes (
				id SERIAL PRIMARY KEY,
				guest_id INTEGER NOT NULL,
				creator_id INTEGER NOT NULL,
				description TEXT,
				created_at TIMESTAMPTZ DEFAULT NOW(),
				updated_at TIMESTAMPTZ DEFAULT NOW(),

				CONSTRAINT fk_guests_notes FOREIGN KEY(guest_id) REFERENCES guests(id)
			)
			`
			_, err := tx.Exec(query)
			return err
		},
		Seed: func(tx *sqlx.Tx) error {
			environment := os.Getenv("ENVIRONMENT")
			if environment != meta.ProdEnvironment && environment != meta.StagingEnvironment {
				notes := []domain.GuestNote{
					{
						GuestID:     1,
						CreatorID:   1,
						Description: gofaker.LoremIpsumParagraph(1, 3, 6, " "),
					},
					{
						GuestID:     1,
						CreatorID:   1,
						Description: gofaker.LoremIpsumParagraph(1, 3, 6, " "),
					},
				}

				for _, note := range notes {
					_, err := tx.NamedExec(`
                    INSERT INTO guest_notes (guest_id, creator_id, description, created_at, updated_at) VALUES (:guest_id, :creator_id, :description, NOW(), NOW())
                `, note)
					if err != nil {
						return err
					}
				}
			}
			return nil
		},
		Rollback: func(tx *sqlx.Tx) error {
			query := `DROP TABLE IF EXISTS guest_notes`
			_, err := tx.Exec(query)
			return err
		},
	}
}
