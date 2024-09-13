package migrations

import (
	"github.com/goplaceapp/goplace-common/pkg/dbhelper"
	"github.com/jmoiron/sqlx"
)

func CreateReservationFeedbackCommentsTable() dbhelper.SqlxMigration {
	return dbhelper.SqlxMigration{
		ID: "2024_04_24_120551_create_reservation_feedback_comments_table",
		Migrate: func(tx *sqlx.Tx) error {
			query := `
            CREATE TABLE IF NOT EXISTS reservation_feedback_comments(
                id SERIAL PRIMARY KEY,
                reservation_feedback_id INTEGER NOT NULL,
				creator_id INTEGER NOT NULL,
				comment TEXT NOT NULL,
                created_at TIMESTAMPTZ DEFAULT NOW(),
                updated_at TIMESTAMPTZ DEFAULT NOW(),

				CONSTRAINT fk_inquiry_comments_reservation_feedback_id FOREIGN KEY(reservation_feedback_id) REFERENCES reservation_feedbacks(id)
            )
            `
			_, err := tx.Exec(query)
			return err
		},
		Seed: func(tx *sqlx.Tx) error {
			return nil
		},
		Rollback: func(tx *sqlx.Tx) error {
			_, err := tx.Exec(`DROP TABLE IF EXISTS reservation_feedback_comments`)
			return err
		},
	}
}
