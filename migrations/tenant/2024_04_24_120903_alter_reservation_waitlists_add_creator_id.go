package migrations

import (
	"github.com/goplaceapp/goplace-common/pkg/dbhelper"
	"github.com/jmoiron/sqlx"
)

func AlterReservationWaitlistsAddCreatorId() dbhelper.SqlxMigration {
	return dbhelper.SqlxMigration{
		ID: "2024_04_24_120903_alter_reservation_waitlists_add_creator_id",
		Migrate: func(tx *sqlx.Tx) error {
			query := `
			    ALTER TABLE reservation_waitlists
				ADD COLUMN IF NOT EXISTS creator_id INTEGER
			`

			_, err := tx.Exec(query)
			return err
		},
	}
}
