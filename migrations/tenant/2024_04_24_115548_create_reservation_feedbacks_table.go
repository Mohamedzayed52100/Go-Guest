package migrations

import (
	"os"

	"github.com/goplaceapp/goplace-common/pkg/dbhelper"
	"github.com/goplaceapp/goplace-common/pkg/meta"
	feedbackDomain "github.com/goplaceapp/goplace-guest/internal/services/reservation-feedback/domain"
	"github.com/goplaceapp/goplace-guest/pkg/reservationservice/domain"
	"github.com/jmoiron/sqlx"
	"github.com/tinygg/gofaker"
)

func CreateReservationFeedbacksTable() dbhelper.SqlxMigration {
	return dbhelper.SqlxMigration{
		ID: "2024_04_24_115548_create_reservation_feedbacks_table",
		Migrate: func(tx *sqlx.Tx) error {
			query := `
			CREATE TABLE IF NOT EXISTS reservation_feedbacks(
				id SERIAL PRIMARY KEY,
				reservation_id INTEGER UNIQUE NOT NULL,
				status_id INTEGER NULL,
				solution_id INTEGER NULL,
				rate INTEGER NOT NULL,
				description TEXT NULL,
				created_at TIMESTAMPTZ DEFAULT NOW(),
				updated_at TIMESTAMPTZ DEFAULT NOW(),

				CONSTRAINT fk_reservation_feedback FOREIGN KEY (reservation_id) REFERENCES reservations(id)
			)
			`

			_, err := tx.Exec(query)
			return err
		},
		Seed: func(tx *sqlx.Tx) error {
			environment := os.Getenv("ENVIRONMENT")
			if environment != meta.ProdEnvironment && environment != meta.StagingEnvironment {
				var reservations []domain.Reservation
				if err := tx.Select(&reservations, `SELECT id FROM reservations LIMIT 2`); err != nil {
					return err
				}

				for _, r := range reservations {
					feedback := feedbackDomain.ReservationFeedback{
						ReservationID: r.ID,
						Rate:          4,
						Description:   gofaker.Sentence(10),
					}

					var count int
					if err := tx.Get(&count, `SELECT COUNT(*) FROM reservation_feedbacks WHERE reservation_id = $1`, r.ID); err != nil {
						return err
					}

					if count == 0 {
						if _, err := tx.NamedExec(`INSERT INTO reservation_feedbacks (reservation_id, rate, description, created_at, updated_at) VALUES (:reservation_id, :rate, :description, now(), now())`, feedback); err != nil {
							return err
						}
					}
				}

			}
			return nil
		},
		Rollback: func(tx *sqlx.Tx) error {
			query := `DROP TABLE IF EXISTS reservation_feedbacks`
			_, err := tx.Exec(query)
			return err
		},
	}
}
