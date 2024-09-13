package migrations

import (
    "github.com/goplaceapp/goplace-common/pkg/dbhelper"
    "github.com/jmoiron/sqlx"
)

func AlterReservationsDropReservationGuestUkey() dbhelper.SqlxMigration {
    return dbhelper.SqlxMigration{
        ID: "2024_04_24_121434_alter_reservations_drop_reservation_guest_ukey",
        Migrate: func(tx *sqlx.Tx) error {
			query := `
                ALTER TABLE reservations DROP CONSTRAINT IF EXISTS reservations_guest_ukey;
                DROP INDEX IF EXISTS reservations_guest_ukey;
            `

			_, err := tx.Exec(query)
			return err
		},
    }
}
