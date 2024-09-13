package migrations

import (
	"os"

	"github.com/goplaceapp/goplace-common/pkg/dbhelper"
	"github.com/goplaceapp/goplace-common/pkg/meta"
	"github.com/goplaceapp/goplace-guest/internal/services/reservation-feedback/domain"
	"github.com/jmoiron/sqlx"
)

func CreateReservationFeedbackSectionsTable() dbhelper.SqlxMigration {
	return dbhelper.SqlxMigration{
		ID: "2024_04_24_120457_create_reservation_feedback_sections_table",
		Migrate: func(tx *sqlx.Tx) error {
			query := `
            CREATE TABLE IF NOT EXISTS reservation_feedback_sections(
                id SERIAL PRIMARY KEY,
                name TEXT NOT NULL,
                branch_id INTEGER NOT NULL,
                created_at TIMESTAMPTZ DEFAULT NOW(),
                updated_at TIMESTAMPTZ DEFAULT NOW(),
                
                FOREIGN KEY (branch_id) REFERENCES branches(id) ON DELETE CASCADE,
                CONSTRAINT unique_name_branch_id UNIQUE (name, branch_id)
            )
            `
			_, err := tx.Exec(query)
			return err
		},
		Seed: func(tx *sqlx.Tx) error {
			environment := os.Getenv("ENVIRONMENT")
			if environment != meta.ProdEnvironment && environment != meta.StagingEnvironment {
				sections := []*domain.ReservationFeedbackSection{
					{
						Name:     "General",
						BranchID: 1,
					},
					{
						Name:     "Kitchen",
						BranchID: 1,
					},
				}

				var count int
				if err := tx.Get(&count, `SELECT COUNT(*) FROM reservation_feedback_sections`); err != nil {
					return err
				}

				if count == 0 {
					for _, s := range sections {
						if _, err := tx.NamedExec(`INSERT INTO reservation_feedback_sections (name, branch_id, created_at, updated_at) VALUES (:name, :branch_id, NOW(), NOW())`, s); err != nil {
							return err
						}
					}
				}
			}

			return nil
		},
		Rollback: func(tx *sqlx.Tx) error {
			_, err := tx.Exec(`DROP TABLE IF EXISTS reservation_feedback_sections`)
			return err
		},
	}
}
