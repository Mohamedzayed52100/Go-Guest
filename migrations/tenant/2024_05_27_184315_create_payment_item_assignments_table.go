package migrations

import (
	"github.com/goplaceapp/goplace-common/pkg/dbhelper"
	"github.com/jmoiron/sqlx"
)

func CreatePaymentItemAssignmentsTable() dbhelper.SqlxMigration {
	return dbhelper.SqlxMigration{
		ID: "2024_05_27_184317_create_payment_item_assignments_table",
		Migrate: func(tx *sqlx.Tx) error {
			query := `
            CREATE TABLE IF NOT EXISTS payment_item_assignments (
                id SERIAL PRIMARY KEY,
                payment_id INTEGER NOT NULL,
                item_id INTEGER NOT NULL,
                price DECIMAL(10, 2) NOT NULL,
                quantity INTEGER NOT NULL,

                CONSTRAINT fk_payment_item_assignments_payment FOREIGN KEY (payment_id) REFERENCES payment_requests(id),
                CONSTRAINT fk_invoice_item_assignments_item FOREIGN KEY (item_id) REFERENCES restaurant_items(id)
            );
            `
			_, err := tx.Exec(query)
			return err
		},
		Rollback: func(tx *sqlx.Tx) error {
			query := "DROP TABLE IF EXISTS payment_item_assignments;"
			_, err := tx.Exec(query)
			return err
		},
	}
}
