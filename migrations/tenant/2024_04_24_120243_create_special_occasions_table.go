package migrations

import (
	"github.com/goplaceapp/goplace-common/pkg/dbhelper"
	"github.com/jmoiron/sqlx"
)

func CreateSpecialOccasionsTable() dbhelper.SqlxMigration {
	return dbhelper.SqlxMigration{
		ID: "2024_04_24_120243_create_special_occasions_table",
		Migrate: func(tx *sqlx.Tx) error {
			query := `
			CREATE TABLE IF NOT EXISTS special_occasions(
				id SERIAL PRIMARY KEY,
				name TEXT NOT NULL,
				color TEXT NOT NULL,
				icon TEXT NOT NULL,
				created_at TIMESTAMPTZ DEFAULT NOW(),
				updated_at TIMESTAMPTZ DEFAULT NOW(),
			    
			    CONSTRAINT special_occasions_name_unique UNIQUE (name)
			)
			`

			_, err := tx.Exec(query)
			return err
		},
		Rollback: func(tx *sqlx.Tx) error {
			query := `DROP TABLE IF EXISTS special_occasions`
			_, err := tx.Exec(query)
			return err
		},
	}
}
