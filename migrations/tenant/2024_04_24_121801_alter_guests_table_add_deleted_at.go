package migrations

import (
	"github.com/goplaceapp/goplace-common/pkg/dbhelper"
	"github.com/jmoiron/sqlx"
)

func AlterGuestsTableAddDeletedAt() dbhelper.SqlxMigration {
	return dbhelper.SqlxMigration{
		ID: "2024_04_24_121801_alter_guests_table_add_deleted_at",
		Migrate: func(tx *sqlx.Tx) error {
			query := `
			ALTER TABLE guests
			ADD COLUMN IF NOT EXISTS deleted_at timestamp NULL;
			`

			_, err := tx.Exec(query)
			return err
		},
	}
}
