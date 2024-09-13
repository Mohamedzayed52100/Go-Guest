package migrations

import (
	"github.com/goplaceapp/goplace-common/pkg/dbhelper"
	"github.com/jmoiron/sqlx"
)

func AlterGuestsAddAddress() dbhelper.SqlxMigration {
	return dbhelper.SqlxMigration{
		ID: "2024_06_13_162619_alter_guests_add_address",
		Migrate: func(tx *sqlx.Tx) error {
			query := `
            ALTER TABLE guests
            ADD COLUMN IF NOT EXISTS address TEXT;
            `
			_, err := tx.Exec(query)
			return err
		},
	}
}
