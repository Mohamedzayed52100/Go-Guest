package migrations

import (
	"os"

	"github.com/goplaceapp/goplace-common/pkg/dbhelper"
	"github.com/goplaceapp/goplace-common/pkg/meta"
	"github.com/goplaceapp/goplace-guest/pkg/reservationservice/domain"
	logDomain "github.com/goplaceapp/goplace-guest/internal/services/reservation-log/domain"
	"github.com/jmoiron/sqlx"
)

func CreateReservationLogsTable() dbhelper.SqlxMigration {
	return dbhelper.SqlxMigration{
		ID: "2024_04_24_115839_create_reservation_logs_table",
		Migrate: func(tx *sqlx.Tx) error {
			query := `
            CREATE TABLE IF NOT EXISTS reservation_logs (
                id SERIAL PRIMARY KEY,
				reservation_id INTEGER NOT NULL,
				creator_id INTEGER,
				made_by TEXT NOT NULL DEFAULT 'System',
				field_name TEXT NOT NULL,
				old_value TEXT NULL,
				new_value TEXT NULL,
                action TEXT NULL,
				created_at TIMESTAMPTZ DEFAULT NOW(),
				updated_at TIMESTAMPTZ DEFAULT NOW(),

				CONSTRAINT fk_logs_reservation FOREIGN KEY (reservation_id) REFERENCES reservations(id)
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
					log := logDomain.ReservationLog{
						ReservationID: r.ID,
						CreatorID:     1,
						Action:        "create",
						FieldName:     "reservation",
					}

					if _, err := tx.NamedExec(`INSERT INTO reservation_logs (reservation_id, creator_id, action, field_name, created_at, updated_at) VALUES (:reservation_id, :creator_id, :action, :field_name, NOW(), NOW())`, log); err != nil {
						return err
					}
				}
			}
			return nil
		},
		Rollback: func(tx *sqlx.Tx) error {
			query := `DROP TABLE IF EXISTS reservation_logs`
			_, err := tx.Exec(query)
			return err
		},
	}
}
