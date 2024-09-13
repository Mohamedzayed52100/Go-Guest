package migrations

import (
	"github.com/goplaceapp/goplace-common/pkg/dbhelper"
	"github.com/jmoiron/sqlx"
)

func AlterSpecialOccasionsAddBranchId() dbhelper.SqlxMigration {
	return dbhelper.SqlxMigration{
		ID: "2024_04_24_121612_alter_special_occasions_add_branch_id",
		Migrate: func(tx *sqlx.Tx) error {
			query := `
			ALTER TABLE special_occasions
			ADD COLUMN IF NOT EXISTS branch_id INTEGER,
			DROP CONSTRAINT IF EXISTS special_occasions_branch_id_fk,
			ADD CONSTRAINT special_occasions_branch_id_fk
			FOREIGN KEY (branch_id)
			REFERENCES branches (id)
			ON DELETE CASCADE;
			`

			_, err := tx.Exec(query)
			return err
		},
	}
}
