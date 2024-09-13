package migrations

import (
	"github.com/goplaceapp/goplace-common/pkg/dbhelper"
	"github.com/jmoiron/sqlx"
)

func AlterReservationsAddDateIndex() dbhelper.SqlxMigration {
	return dbhelper.SqlxMigration{
		ID: "2024_04_24_120748_alter_reservations_add_date_index",
		Migrate: func(tx *sqlx.Tx) error {
			query := `
				CREATE INDEX IF NOT EXISTS reservations_date_index ON reservations (date);
			`
			_, err := tx.Exec(query)
			return err
		},
	}
}
