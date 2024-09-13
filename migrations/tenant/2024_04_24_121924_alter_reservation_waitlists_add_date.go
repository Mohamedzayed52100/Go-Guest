package migrations

import (
	"github.com/goplaceapp/goplace-common/pkg/dbhelper"
	"github.com/jmoiron/sqlx"
)

func AlterReservationWaitlistsAddDate() dbhelper.SqlxMigration {
	return dbhelper.SqlxMigration{
		ID: "2024_04_24_121924_alter_reservation_waitlists_add_date",
		Migrate: func(tx *sqlx.Tx) error {
			query := `
			ALTER TABLE reservation_waitlists
			ADD COLUMN IF NOT EXISTS "date" DATE;
			`

			if _, err := tx.Exec(query); err != nil {
				return err
			}

			query = `
			UPDATE reservation_waitlists
			SET "date" = "created_at";
			`

			if _, err := tx.Exec(query); err != nil {
				return err
			}

			return nil
		},
	}
}
