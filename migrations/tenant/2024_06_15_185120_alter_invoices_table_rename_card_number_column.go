package migrations

import (
	"github.com/goplaceapp/goplace-common/pkg/dbhelper"
	"github.com/jmoiron/sqlx"
)

func AlterInvoicesTableRenameCardNumberColumn() dbhelper.SqlxMigration {
	return dbhelper.SqlxMigration{
		ID: "2024_06_15_185120_alter_invoices_table_rename_card_number_column",
		Migrate: func(tx *sqlx.Tx) error {
			query := `
            ALTER TABLE invoices
            RENAME COLUMN card_number TO last_four_digits;
            `
			_, err := tx.Exec(query)
			return err
		},
		Seed: func(tx *sqlx.Tx) error {
			return nil
		},
		Rollback: func(tx *sqlx.Tx) error {
			query := ""
			_, err := tx.Exec(query)
			return err
		},
	}
}
