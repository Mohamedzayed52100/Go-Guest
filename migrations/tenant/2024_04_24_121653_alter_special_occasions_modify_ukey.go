package migrations

import (
	"github.com/goplaceapp/goplace-common/pkg/dbhelper"
	"github.com/jmoiron/sqlx"
)

func AlterSpecialOccasionsModifyUkey() dbhelper.SqlxMigration {
	return dbhelper.SqlxMigration{
		ID: "2024_04_24_121653_alter_special_occasions_modify_ukey",
		Migrate: func(tx *sqlx.Tx) error {
			query := `
            ALTER TABLE special_occasions
			DROP CONSTRAINT IF EXISTS special_occasions_name_unique,
			DROP CONSTRAINT IF EXISTS special_occasions_name_ukey,
			ADD CONSTRAINT special_occasions_name_ukey UNIQUE (branch_id, name);
			`

			_, err := tx.Exec(query)
			return err
		},
	}
}
