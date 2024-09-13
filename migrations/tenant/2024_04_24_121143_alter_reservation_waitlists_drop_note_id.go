package migrations

import (
	"github.com/goplaceapp/goplace-common/pkg/dbhelper"
	"github.com/jmoiron/sqlx"
)

func AlterReservationWaitlistsDropNoteId() dbhelper.SqlxMigration {
	return dbhelper.SqlxMigration{
		ID: "2024_04_24_121143_alter_reservation_waitlists_drop_note_id",
		Migrate: func(tx *sqlx.Tx) error {
			query := `
			ALTER TABLE reservation_waitlists
			DROP CONSTRAINT fk_waitlist_note,
			ADD CONSTRAINT fk_waitlist_note FOREIGN KEY (note_id) REFERENCES reservation_waitlist_notes(id) ON DELETE CASCADE;
			`

			_, err := tx.Exec(query)
			return err

		},
	}
}
