package migrations

import (
	"github.com/goplaceapp/goplace-common/pkg/dbhelper"
	"github.com/jmoiron/sqlx"
)

func CreateReservationsTable() dbhelper.SqlxMigration {
	return dbhelper.SqlxMigration{
		ID: "2024_04_24_115528_create_reservations_table",
		Migrate: func(tx *sqlx.Tx) error {
			query := `
			CREATE TABLE IF NOT EXISTS reservations(
				id SERIAL PRIMARY KEY,
				guest_id INTEGER NOT NULL,
				branch_id INTEGER NOT NULL,
				shift_id INTEGER NOT NULL,
				seating_area_id INTEGER NOT NULL,
				status_id INTEGER NOT NULL,
				special_occasion_id INTEGER NULL,
				note_id INTEGER NULL,
				guests_number INTEGER NOT NULL,
				seated_guests INTEGER NULL,
				date DATE NOT NULL,
				time TIMESTAMP NOT NULL,
				creation_duration NUMERIC(10,2) NOT NULL,
				reserved_via TEXT NOT NULL,
				check_in TIMESTAMP NULL,
				check_out TIMESTAMP NULL,
				creator_id INTEGER NULL,

				created_at TIMESTAMPTZ DEFAULT NOW(),
				updated_at TIMESTAMPTZ DEFAULT NOW(),

				CONSTRAINT fk_reservations_guests FOREIGN KEY (guest_id) REFERENCES guests(id),
				CONSTRAINT fk_reservations_branches FOREIGN KEY (branch_id) REFERENCES branches(id),
				CONSTRAINT fk_reservations_shifts FOREIGN KEY (shift_id) REFERENCES shifts(id),
				CONSTRAINT fk_reservations_seating_areas FOREIGN KEY (seating_area_id) REFERENCES seating_areas(id),
				CONSTRAINT fk_reservations_statuses FOREIGN KEY (status_id) REFERENCES reservation_statuses(id),
				CONSTRAINT fk_reservations_special_occasions FOREIGN KEY (special_occasion_id) REFERENCES special_occasions(id),
				CONSTRAINT fk_reservations_note FOREIGN KEY (note_id) REFERENCES reservation_notes(id)
			)
			`

			if _, err := tx.Exec(query); err != nil {
				return err
			}

			return nil
		},
		Seed: func(tx *sqlx.Tx) error {
			return nil
		},
		Rollback: func(tx *sqlx.Tx) error {
			query := `DROP TABLE IF EXISTS reservations`
			_, err := tx.Exec(query)
			return err
		},
	}
}
