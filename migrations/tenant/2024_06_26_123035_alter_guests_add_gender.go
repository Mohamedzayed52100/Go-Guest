package migrations

import (
	"github.com/goplaceapp/goplace-common/pkg/dbhelper"
	"github.com/jmoiron/sqlx"
)

func AlterGuestsAddGender() dbhelper.SqlxMigration {
	return dbhelper.SqlxMigration{
		ID: "2024_06_26_123035_alter_guests_add_gender",
		Migrate: func(tx *sqlx.Tx) error {
			query := `
            ALTER TABLE guests
            ADD COLUMN IF NOT EXISTS gender TEXT DEFAULT 'male';
            `
			_, err := tx.Exec(query)
			return err
		},
	}
}
