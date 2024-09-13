package migrations

import (
	"github.com/goplaceapp/goplace-common/pkg/dbhelper"
	"github.com/jmoiron/sqlx"
)

func AlterReservationsAddReservationGuestUkey() dbhelper.SqlxMigration {
	return dbhelper.SqlxMigration{
		ID: "2024_04_24_121229_alter_reservations_add_reservation_guest_ukey",
		Migrate: func(tx *sqlx.Tx) error {
			query := `
			   CREATE UNIQUE INDEX IF NOT EXISTS reservations_guest_ukey
			   ON reservations (guest_id, date)
			   WHERE date >= '2024-03-04';
			`

			_, err := tx.Exec(query)
			return err
		},
	}
}
