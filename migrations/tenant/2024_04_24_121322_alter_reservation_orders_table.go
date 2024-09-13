package migrations

import (
	"github.com/goplaceapp/goplace-common/pkg/dbhelper"
	"github.com/jmoiron/sqlx"
)

func AlterOrdersTable() dbhelper.SqlxMigration {
	return dbhelper.SqlxMigration{
		ID: "2024_04_24_121322_alter_reservation_orders_table",
		Migrate: func(tx *sqlx.Tx) error {
			query := `
            -- Drop old columns
            ALTER TABLE IF EXISTS reservation_orders DROP COLUMN IF EXISTS item_name;
            ALTER TABLE IF EXISTS reservation_orders DROP COLUMN IF EXISTS cost;
            ALTER TABLE IF EXISTS reservation_orders DROP COLUMN IF EXISTS quantity;

            -- Add new columns
            ALTER TABLE IF EXISTS reservation_orders ADD COLUMN IF NOT EXISTS discount_amount NUMERIC(10, 2);
            ALTER TABLE IF EXISTS reservation_orders ADD COLUMN IF NOT EXISTS discount_reason TEXT;
            ALTER TABLE IF EXISTS reservation_orders ADD COLUMN IF NOT EXISTS prevailing_tax NUMERIC(10, 2);
            ALTER TABLE IF EXISTS reservation_orders ADD COLUMN IF NOT EXISTS tax NUMERIC(10, 2);
            ALTER TABLE IF EXISTS reservation_orders ADD COLUMN IF NOT EXISTS subtotal NUMERIC(10, 2);
            ALTER TABLE IF EXISTS reservation_orders ADD COLUMN IF NOT EXISTS final_total NUMERIC(10, 2);

            -- Drop the foreign key constraint
            ALTER TABLE IF EXISTS reservation_orders DROP CONSTRAINT IF EXISTS fk_reservations_orders;

            -- Build a new foreign key constraint
			ALTER TABLE IF EXISTS reservation_orders 
			DROP CONSTRAINT IF EXISTS fk_reservations_order_reservations,
            ADD CONSTRAINT fk_reservations_order_reservations FOREIGN KEY (reservation_id) REFERENCES reservations(id) ON DELETE CASCADE;

            -- Build a new unique constraint
            ALTER TABLE IF EXISTS reservation_orders 
			DROP CONSTRAINT IF EXISTS unique_reservation_order,
			ADD CONSTRAINT unique_reservation_order UNIQUE (reservation_id);
            `
			_, err := tx.Exec(query)
			return err
		},
		Rollback: func(tx *sqlx.Tx) error {
			query := `
            -- Drop new columns
			ALTER TABLE IF EXISTS reservation_orders DROP COLUMN IF EXISTS discount_amount;
			ALTER TABLE IF EXISTS reservation_orders DROP COLUMN IF EXISTS discount_reason;
			ALTER TABLE IF EXISTS reservation_orders DROP COLUMN IF EXISTS prevailing_tax;	
			ALTER TABLE IF EXISTS reservation_orders DROP COLUMN IF EXISTS tax;
			ALTER TABLE IF EXISTS reservation_orders DROP COLUMN IF EXISTS subtotal;
			ALTER TABLE IF EXISTS reservation_orders DROP COLUMN IF EXISTS final_total;

			-- Drop the foreign key constraint
			ALTER TABLE IF EXISTS reservation_orders DROP CONSTRAINT IF EXISTS fk_reservations_order_reservations;

			-- Drop the unique constraint
			ALTER TABLE IF EXISTS reservation_orders DROP CONSTRAINT IF EXISTS unique_reservation_order;
			`
			_, err := tx.Exec(query)
			return err
		},
	}
}
