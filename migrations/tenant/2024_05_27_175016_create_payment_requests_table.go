package migrations

import (
	"github.com/goplaceapp/goplace-common/pkg/dbhelper"
	"github.com/jmoiron/sqlx"
)

func CreatePaymentRequestsTable() dbhelper.SqlxMigration {
	return dbhelper.SqlxMigration{
		ID: "2024_05_27_175021_create_payment_requests_table",
		Migrate: func(tx *sqlx.Tx) error {
			query := `
            CREATE TABLE IF NOT EXISTS payment_requests (
                id SERIAL PRIMARY KEY,
                branch_id INTEGER NOT NULL,
                reservation_id INTEGER NOT NULL,
                delivery TEXT NOT NULL,
				date DATE NOT NULL,
				special_request_id INTEGER NULL,
				note TEXT NULL,
                created_at TIMESTAMPTZ DEFAULT NOW(),
                updated_at TIMESTAMPTZ DEFAULT NOW(),

                CONSTRAINT fk_payment_requests_branch FOREIGN KEY (branch_id) REFERENCES branches(id),
                CONSTRAINT fk_payment_requests_reservation FOREIGN KEY (reservation_id) REFERENCES reservations(id)
            );
            `
			_, err := tx.Exec(query)
			return err
		},
		Rollback: func(tx *sqlx.Tx) error {
			query := "DROP TABLE IF EXISTS payment_requests;"
			_, err := tx.Exec(query)
			return err
		},
	}
}
