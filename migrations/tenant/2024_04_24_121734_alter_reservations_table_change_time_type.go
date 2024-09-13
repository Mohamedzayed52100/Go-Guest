package migrations

import (
	"github.com/goplaceapp/goplace-common/pkg/dbhelper"
	"github.com/jmoiron/sqlx"
)

func AlterReservationsTableChangeTimeType() dbhelper.SqlxMigration {
	return dbhelper.SqlxMigration{
		ID: "2024_04_24_121734_alter_reservations_table_change_time_type",
		Migrate: func(tx *sqlx.Tx) error {
			_, err := tx.Exec(`
				ALTER TABLE reservations
				ALTER COLUMN time TYPE TIME
			`)

			return err
		},
	}
}
