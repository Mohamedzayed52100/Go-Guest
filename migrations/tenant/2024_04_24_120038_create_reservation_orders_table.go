package migrations

import (
	"github.com/goplaceapp/goplace-common/pkg/dbhelper"
	"github.com/jmoiron/sqlx"
)

func CreateReservationOrdersTable() dbhelper.SqlxMigration {
	return dbhelper.SqlxMigration{
		ID: "2024_04_24_120038_create_reservation_orders_table",
		Migrate: func(tx *sqlx.Tx) error {
			query := `
            CREATE TABLE IF NOT EXISTS reservation_orders (
                id SERIAL PRIMARY KEY,
                reservation_id INTEGER NOT NULL,
                item_name TEXT,
                cost NUMERIC(10, 2),
                quantity INTEGER,
				created_at TIMESTAMPTZ DEFAULT NOW(),
				updated_at TIMESTAMPTZ DEFAULT NOW(),

				CONSTRAINT fk_reservations_orders FOREIGN KEY (reservation_id) REFERENCES reservations(id)
            )
            `
			_, err := tx.Exec(query)
			return err
		},
		Seed: func(tx *sqlx.Tx) error {
			return nil
		},
		Rollback: func(tx *sqlx.Tx) error {
			query := `DROP TABLE IF EXISTS reservation_orders`
			_, err := tx.Exec(query)
			return err
		},
	}
}
