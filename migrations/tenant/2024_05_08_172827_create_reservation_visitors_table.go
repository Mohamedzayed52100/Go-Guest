package migrations

import (
	"github.com/goplaceapp/goplace-common/pkg/dbhelper"
	"github.com/jmoiron/sqlx"
)

func CreateReservationVisitorsTable() dbhelper.SqlxMigration {
	return dbhelper.SqlxMigration{
		ID: "2024_05_08_172827_create_reservation_visitors_table",
		Migrate: func(tx *sqlx.Tx) error {
			query := `
            CREATE TABLE IF NOT EXISTS reservation_visitors (
                id SERIAL PRIMARY KEY,
                guest_id INTEGER NOT NULL,
                reservation_id INTEGER NOT NULL,
                created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
                updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,

                FOREIGN KEY (guest_id) REFERENCES guests(id) ON DELETE CASCADE,
                FOREIGN KEY (reservation_id) REFERENCES reservations(id) ON DELETE CASCADE
            );
            `
			_, err := tx.Exec(query)
			return err
		},
		Rollback: func(tx *sqlx.Tx) error {
			query := "DROP TABLE IF EXISTS reservation_visitors;"
			_, err := tx.Exec(query)
			return err
		},
	}
}
