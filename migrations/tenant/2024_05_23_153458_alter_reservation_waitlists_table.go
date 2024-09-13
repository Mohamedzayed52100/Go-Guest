package migrations

import (
	"github.com/goplaceapp/goplace-common/pkg/dbhelper"
	"github.com/jmoiron/sqlx"
)

func AlterReservationWaitlistsTable() dbhelper.SqlxMigration {
	return dbhelper.SqlxMigration{
		ID: "2024_05_23_153458_alter_reservation_waitlists_table",
		Migrate: func(tx *sqlx.Tx) error {
			query := `
			ALTER TABLE reservation_waitlists
			ADD COLUMN IF NOT EXISTS "date" DATE;
			`

			if _, err := tx.Exec(query); err != nil {
				return err
			}

			// update date column with created_at value
			query = `UPDATE reservation_waitlists SET date = created_at;`

			if _, err := tx.Exec(query); err != nil {
				return err
			}

			return nil
		},
	}
}
