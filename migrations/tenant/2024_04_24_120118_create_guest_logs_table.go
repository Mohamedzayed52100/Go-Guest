package migrations

import (
	"os"

	"github.com/goplaceapp/goplace-common/pkg/dbhelper"
	"github.com/goplaceapp/goplace-common/pkg/meta"
	logDomain "github.com/goplaceapp/goplace-guest/internal/services/guest-log/domain"
	"github.com/goplaceapp/goplace-guest/pkg/guestservice/domain"
	"github.com/jmoiron/sqlx"
)

func CreateGuestLogsTable() dbhelper.SqlxMigration {
	return dbhelper.SqlxMigration{
		ID: "2024_04_24_120118_create_guest_logs_table",
		Migrate: func(tx *sqlx.Tx) error {
			query := `
            CREATE TABLE IF NOT EXISTS guest_logs (
                id SERIAL PRIMARY KEY,
				guest_id INTEGER NOT NULL,
                creator_id INTEGER,
				made_by TEXT NOT NULL DEFAULT 'System',
				field_name TEXT NOT NULL,
				old_value TEXT NULL,
				new_value TEXT NULL,
                action TEXT NULL,
				created_at TIMESTAMPTZ DEFAULT NOW(),
				updated_at TIMESTAMPTZ DEFAULT NOW(),

				CONSTRAINT fk_logs_guest FOREIGN KEY (guest_id) REFERENCES guests(id)
            )
            `
			_, err := tx.Exec(query)
			return err
		},
		Seed: func(tx *sqlx.Tx) error {
			environment := os.Getenv("ENVIRONMENT")
			if environment != meta.ProdEnvironment && environment != meta.StagingEnvironment {
				var guests []domain.Guest
				if err := tx.Select(&guests, `SELECT id FROM guests LIMIT 2`); err != nil {
					return err
				}
				for _, g := range guests {
					log := logDomain.GuestLog{
						CreatorID: 1,
						GuestID:   g.ID,
						Action:    "create",
						FieldName: "guest",
					}

					if _, err := tx.NamedExec(`INSERT INTO guest_logs (guest_id, creator_id, action, field_name, created_at, updated_at) VALUES (:guest_id, :creator_id, :action, :field_name, NOW(), NOW())`, log); err != nil {
						return err
					}
				}
			}
			return nil
		},
		Rollback: func(tx *sqlx.Tx) error {
			query := `DROP TABLE IF EXISTS activity_logs`
			_, err := tx.Exec(query)
			return err
		},
	}
}
