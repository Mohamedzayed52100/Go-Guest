package migrations

import (
	"github.com/goplaceapp/goplace-common/pkg/dbhelper"
	"github.com/jmoiron/sqlx"
)

func CreateReservationFeedbackSectionAssignmentsTable() dbhelper.SqlxMigration {
	return dbhelper.SqlxMigration{
		ID: "2024_04_24_120525_create_reservation_feedback_section_assignments_table",
		Migrate: func(tx *sqlx.Tx) error {
			query := `
            CREATE TABLE IF NOT EXISTS reservation_feedback_section_assignments(
                id SERIAL PRIMARY KEY,
                feedback_id INTEGER NOT NULL,
                section_id INTEGER NOT NULL,
                created_at TIMESTAMPTZ DEFAULT NOW(),
                updated_at TIMESTAMPTZ DEFAULT NOW(),

				CONSTRAINT fk_reservation_feedbacks FOREIGN KEY (feedback_id) REFERENCES reservation_feedbacks(id) ON DELETE CASCADE,
                CONSTRAINT fk_reservation_sections FOREIGN KEY (section_id) REFERENCES reservation_feedback_sections(id) ON DELETE CASCADE
            )
            `
			_, err := tx.Exec(query)
			return err
		},
		Seed: func(tx *sqlx.Tx) error {
			return nil
		},
		Rollback: func(tx *sqlx.Tx) error {
			_, err := tx.Exec(`DROP TABLE IF EXISTS reservation_feedback_section_assignments`)
			return err
		},
	}
}
