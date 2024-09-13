package migrations

import (
	"github.com/goplaceapp/goplace-common/pkg/dbhelper"
	"github.com/jmoiron/sqlx"
)

func AlterReservationOrdersAddOrderId() dbhelper.SqlxMigration {
	return dbhelper.SqlxMigration{
		ID: "2024_04_24_121526_alter_reservation_orders_add_order_id",
		Migrate: func(tx *sqlx.Tx) error {
			query := `
                ALTER TABLE reservation_orders
                ADD COLUMN IF NOT EXISTS order_id TEXT;
            `
			_, err := tx.Exec(query)
			return err
		},
		Rollback: func(tx *sqlx.Tx) error {
			query := `
                ALTER TABLE reservation_orders
                DROP COLUMN order_id;
            `
			_, err := tx.Exec(query)
			return err
		},
	}
}
