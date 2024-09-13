package migrations

import (
	"github.com/goplaceapp/goplace-common/pkg/dbhelper"
	"github.com/jmoiron/sqlx"
)

func AlterReservationsAddDeletedAt() dbhelper.SqlxMigration {
	return dbhelper.SqlxMigration{
		ID: "2024_04_24_120842_alter_reservations_add_deleted_at",
		Migrate: func(tx *sqlx.Tx) error {
			query := `
				ALTER TABLE reservations
				ADD COLUMN IF NOT EXISTS deleted_at timestamp NULL
			`
			_, err := tx.Exec(query)
			return err
		},
	}
}
