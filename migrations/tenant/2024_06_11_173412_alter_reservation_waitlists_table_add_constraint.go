package migrations

import (
	"github.com/goplaceapp/goplace-common/pkg/dbhelper"
	"github.com/jmoiron/sqlx"
)

func AlterReservationWaitlistsTableAddConstraint() dbhelper.SqlxMigration {
	return dbhelper.SqlxMigration{
		ID: "2024_06_11_173412_alter_reservation_waitlists_table_add_constraint",
		Migrate: func(tx *sqlx.Tx) error {
			query := `
            ALTER TABLE reservation_waitlists
            DROP CONSTRAINT waitlist_guest_shift_id_ukey;

            ALTER TABLE reservation_waitlists
            ADD CONSTRAINT waitlist_guest_shift_id_date_ukey UNIQUE(guest_id, shift_id, date);
            `
			_, err := tx.Exec(query)
			return err
		},
		Seed: func(tx *sqlx.Tx) error {
			return nil
		},
		Rollback: func(tx *sqlx.Tx) error {
			query := "ALTER TABLE reservation_waitlists DROP CONSTRAINT waitlist_guest_shift_id_date_ukey;"
			_, err := tx.Exec(query)
			return err
		},
	}
}
