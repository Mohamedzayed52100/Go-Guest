package migrations

import (
	"github.com/goplaceapp/goplace-common/pkg/dbhelper"
	"github.com/jmoiron/sqlx"
)

func AlterReservationsDropNoteId() dbhelper.SqlxMigration {
	return dbhelper.SqlxMigration{
		ID: "2024_04_24_120655_alter_reservations_drop_note_id",
		Migrate: func(tx *sqlx.Tx) error {
			query := `
			ALTER TABLE reservations
			DROP COLUMN IF EXISTS note_id;
			`
			_, err := tx.Exec(query)
			return err
		},
	}
}
