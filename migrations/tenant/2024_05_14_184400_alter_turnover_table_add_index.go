package migrations

import (
	"github.com/goplaceapp/goplace-common/pkg/dbhelper"
	"github.com/jmoiron/sqlx"
)

func AlterTurnoverTableAddIndex() dbhelper.SqlxMigration {
	return dbhelper.SqlxMigration{
		ID: "2024_05_14_184400_alter_turnover_table_add_index",
		Migrate: func(tx *sqlx.Tx) error {
			query := `
            CREATE INDEX IF NOT EXISTS turnover_shift_index ON turnover (shift_id);
            `
			_, err := tx.Exec(query)
			return err
		},
	}
}
