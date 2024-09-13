package migrations

import (
	"github.com/goplaceapp/goplace-common/pkg/dbhelper"
	"github.com/jmoiron/sqlx"
)

func AlterReservationNotesAddReservationId() dbhelper.SqlxMigration {
	return dbhelper.SqlxMigration{
		ID: "2024_04_24_120720_alter_reservation_notes_add_reservation_id",
		Migrate: func(tx *sqlx.Tx) error {
			query := `
			ALTER TABLE reservation_notes
			ADD COLUMN IF NOT EXISTS reservation_id INTEGER,
			DROP CONSTRAINT IF EXISTS reservation_notes_reservation_id_fkey,
			ADD CONSTRAINT reservation_notes_reservation_id_fkey
			FOREIGN KEY (reservation_id) REFERENCES reservations(id)
			ON DELETE CASCADE;
			`
			_, err := tx.Exec(query)
			return err
		},
	}
}
