package migrations

import (
	"github.com/goplaceapp/goplace-common/pkg/dbhelper"
	"github.com/jmoiron/sqlx"
)

func CreateReservationWaitlistsTable() dbhelper.SqlxMigration {
	return dbhelper.SqlxMigration{
		ID: "2024_04_24_120406_create_reservation_waitlists_table",
		Migrate: func(tx *sqlx.Tx) error {
			query := `
				CREATE TABLE IF NOT EXISTS reservation_waitlists (
    			id SERIAL PRIMARY KEY,
    			guest_id INTEGER NOT NULL,
    			note_id INTEGER NULL,
				shift_id INTEGER NOT NULL,
    			seating_area_id INTEGER NOT NULL,
    			guests_number INTEGER NOT NULL,
    			waiting_time INTEGER NOT NULL,
    			created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    			updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

				CONSTRAINT fk_waitlist_note FOREIGN KEY(note_id) REFERENCES reservation_notes(id),
				CONSTRAINT fk_waitlist_guest FOREIGN KEY(guest_id) REFERENCES guests(id),
				CONSTRAINT fk_waitlist_shift FOREIGN KEY(shift_id) REFERENCES shifts(id),
				CONSTRAINT fk_waitlist_seating_area FOREIGN KEY(seating_area_id) REFERENCES seating_areas(id),
				CONSTRAINT waitlist_guest_shift_id_ukey UNIQUE(guest_id, shift_id)
			)`

			if _, err := tx.Exec(query); err != nil {
				return err
			}

			return nil
		},
		Rollback: func(tx *sqlx.Tx) error {
			query := `DROP TABLE IF EXISTS reservation_waitlists`
			_, err := tx.Exec(query)

			return err
		},
	}
}
