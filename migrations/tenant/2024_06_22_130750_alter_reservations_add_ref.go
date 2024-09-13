package migrations

import (
	"github.com/goplaceapp/goplace-common/pkg/dbhelper"
	"github.com/jmoiron/sqlx"
)

func AlterReservationsAddRef() dbhelper.SqlxMigration {
	return dbhelper.SqlxMigration{
		ID: "2024_06_22_130750_alter_reservations_add_ref",
		Migrate: func(tx *sqlx.Tx) error {
			query := `
            ALTER TABLE reservations
            ADD COLUMN IF NOT EXISTS reservation_ref TEXT NOT NULL DEFAULT '';
            `
			_, err := tx.Exec(query)
			if err != nil {
				return err
			}

			query = `
            UPDATE reservations SET reservation_ref = id::TEXT;
            `

			_, err = tx.Exec(query)
			return err
		},
	}
}
