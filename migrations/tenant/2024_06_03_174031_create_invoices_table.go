package migrations

import (
	"github.com/goplaceapp/goplace-common/pkg/dbhelper"
	"github.com/jmoiron/sqlx"
)

func CreateInvoicesTable() dbhelper.SqlxMigration {
	return dbhelper.SqlxMigration{
		ID: "2024_06_03_174032_create_invoices_table",
		Migrate: func(tx *sqlx.Tx) error {
			query := `
            CREATE TABLE IF NOT EXISTS invoices (
                id SERIAL PRIMARY KEY,
                invoice_id TEXT NOT NULL,
                payment_request_id INTEGER NOT NULL,
                customer_id TEXT NOT NULL,
                card_number TEXT NOT NULL,
                exp_date TEXT NOT NULL,
                card_type TEXT NOT NULL,
                status TEXT NOT NULL DEFAULT 'unpaid',
                currency TEXT NOT NULL,

                CONSTRAINT fk_invoices_payment_requests FOREIGN KEY (payment_request_id) REFERENCES payment_requests(id)
            );
            `
			_, err := tx.Exec(query)
			return err
		},
		Seed: func(tx *sqlx.Tx) error {
			return nil
		},
		Rollback: func(tx *sqlx.Tx) error {
			query := "DROP TABLE IF EXISTS invoices;"
			_, err := tx.Exec(query)
			return err
		},
	}
}
