package migrations

import (
	"github.com/goplaceapp/goplace-common/pkg/dbhelper"
	"github.com/jmoiron/sqlx"
)

func CreatePaymentRequestContactsTable() dbhelper.SqlxMigration {
	return dbhelper.SqlxMigration{
		ID: "2024_06_16_000919_create_payment_request_contacts_table",
		Migrate: func(tx *sqlx.Tx) error {
			query := `
            CREATE TABLE IF NOT EXISTS payment_request_contacts (
                id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
                payment_request_id INTEGER NOT NULL,
                contact_id INTEGER NOT NULL,
                created_at TIMESTAMP WITH TIME ZONE DEFAULT now(),
                updated_at TIMESTAMP WITH TIME ZONE DEFAULT now(),

                CONSTRAINT fk_payment_request_contacts_payment_request_id FOREIGN KEY (payment_request_id) REFERENCES payment_requests(id),
                CONSTRAINT fk_payment_request_contacts_contact_id FOREIGN KEY (contact_id) REFERENCES guests(id),
                CONSTRAINT uq_payment_request_contacts UNIQUE (payment_request_id, contact_id)
            );
            `
			_, err := tx.Exec(query)
			return err
		},
		Seed: func(tx *sqlx.Tx) error {
			return nil
		},
		Rollback: func(tx *sqlx.Tx) error {
			query := "DROP TABLE IF EXISTS payment_request_contacts;"
			_, err := tx.Exec(query)
			return err
		},
	}
}
