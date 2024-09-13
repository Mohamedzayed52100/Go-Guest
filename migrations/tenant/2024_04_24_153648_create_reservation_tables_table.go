package migrations

import (
	"github.com/goplaceapp/goplace-common/pkg/dbhelper"
	"github.com/jmoiron/sqlx"
)

func CreateReservationTablesTable() dbhelper.SqlxMigration {
	return dbhelper.SqlxMigration{
		ID: "2024_04_24_153648_create_reservation_tables_table",
		Migrate: func(tx *sqlx.Tx) error {
			query := `
            CREATE TABLE IF NOT EXISTS reservation_tables (
                id SERIAL PRIMARY KEY,
                reservation_id INTEGER NOT NULL,
                table_id INTEGER NOT NULL,
				created_at TIMESTAMPTZ DEFAULT NOW(),
				updated_at TIMESTAMPTZ DEFAULT NOW(),

                CONSTRAINT fk_reservation_tables_reservation FOREIGN KEY (reservation_id) REFERENCES reservations(id),
                CONSTRAINT fk_reservation_tables_table FOREIGN KEY (table_id) REFERENCES tables(id)
            )
            `
			_, err := tx.Exec(query)
			return err
		},
		Rollback: func(tx *sqlx.Tx) error {
			query := `DROP TABLE IF EXISTS reservations_table`
			_, err := tx.Exec(query)
			return err
		},
	}
}
