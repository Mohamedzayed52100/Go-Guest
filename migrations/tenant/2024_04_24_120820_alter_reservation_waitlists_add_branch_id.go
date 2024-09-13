package migrations

import (
	"github.com/goplaceapp/goplace-common/pkg/dbhelper"
	"github.com/jmoiron/sqlx"
)

func AlterReservationWaitlistsAddBranchId() dbhelper.SqlxMigration {
	return dbhelper.SqlxMigration{
		ID: "2024_04_24_120820_alter_reservation_waitlists_add_branch_id",
		Migrate: func(tx *sqlx.Tx) error {
			query := `
			ALTER TABLE reservation_waitlists
			ADD COLUMN IF NOT EXISTS branch_id INTEGER,
			DROP CONSTRAINT IF EXISTS reservation_waitlists_branch_id_fkey,
			ADD CONSTRAINT reservation_waitlists_branch_id_fkey
			FOREIGN KEY (branch_id) REFERENCES branches(id)
			ON DELETE CASCADE;
			`
			_, err := tx.Exec(query)
			return err
		},
		Rollback: func(tx *sqlx.Tx) error {
			query := `
			ALTER TABLE reservation_waitlists
			DROP COLUMN IF EXISTS branch_id;
			`
			_, err := tx.Exec(query)
			if err != nil {
				return err
			}

			return nil
		},
	}
}
