package migrations

import (
	"os"

	"github.com/goplaceapp/goplace-common/pkg/dbhelper"
	"github.com/goplaceapp/goplace-common/pkg/meta"
	"github.com/goplaceapp/goplace-guest/internal/services/reservation-waitlist/domain"
	"github.com/jmoiron/sqlx"
)

func CreateReservationWaitlistLogsTable() dbhelper.SqlxMigration {
	return dbhelper.SqlxMigration{
		ID: "2024_04_24_121054_create_reservation_waitlist_logs_table",
		Migrate: func(tx *sqlx.Tx) error {
			query := `
            CREATE TABLE IF NOT EXISTS reservation_waitlist_logs (
                id SERIAL PRIMARY KEY,
				reservation_waitlist_id INTEGER NOT NULL,
				creator_id INTEGER,
				made_by TEXT NOT NULL DEFAULT 'System',
				field_name TEXT NOT NULL,
				old_value TEXT NULL,
				new_value TEXT NULL,
                action TEXT NULL,
				created_at TIMESTAMPTZ DEFAULT NOW(),
				updated_at TIMESTAMPTZ DEFAULT NOW(),

				CONSTRAINT fk_logs_reservation FOREIGN KEY (reservation_waitlist_id) REFERENCES reservation_waitlists(id)
            )
            `
			_, err := tx.Exec(query)
			return err
		},
		Seed: func(tx *sqlx.Tx) error {
			environment := os.Getenv("ENVIRONMENT")
			if environment != meta.ProdEnvironment && environment != meta.StagingEnvironment {
				var reservations []domain.ReservationWaitlist
				if err := tx.Select(&reservations, `SELECT id FROM reservation_waitlists`); err != nil {
					return err
				}
				for _, r := range reservations {
					log := domain.ReservationWaitlistLog{
						ReservationWaitlistID: r.ID,
						CreatorID:             1,
						Action:                "create",
						FieldName:             "reservation-waitlist",
					}

					if _, err := tx.NamedExec(`INSERT INTO reservation_waitlist_logs (reservation_waitlist_id, creator_id, action, field_name, created_at, updated_at) VALUES (:reservation_waitlist_id, :creator_id, :action, :field_name, NOW(), NOW())`, log); err != nil {
						return err
					}
				}
			}
			return nil
		},
		Rollback: func(tx *sqlx.Tx) error {
			query := `DROP TABLE IF EXISTS reservation_waitlist_logs`
			_, err := tx.Exec(query)
			return err
		},
	}
}
