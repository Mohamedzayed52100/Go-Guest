package migrations

import (
	"os"
	"time"

	"github.com/goplaceapp/goplace-common/pkg/dbhelper"
	"github.com/goplaceapp/goplace-common/pkg/meta"
	"github.com/goplaceapp/goplace-guest/pkg/guestservice/domain"
	"github.com/jmoiron/sqlx"
	"github.com/tinygg/gofaker"
)

func CreateGuestsTable() dbhelper.SqlxMigration {
	return dbhelper.SqlxMigration{
		ID: "2024_04_24_115452_create_guests_table",
		Migrate: func(tx *sqlx.Tx) error {
			query := `
            CREATE TABLE IF NOT EXISTS guests (
                id SERIAL PRIMARY KEY,
                first_name TEXT NOT NULL,
                last_name TEXT NOT NULL,
                email TEXT,
                phone_number TEXT NOT NULL UNIQUE,
                language TEXT NULL,
                birthdate DATE NULL,
                created_at timestamptz DEFAULT NOW(),
                updated_at timestamptz DEFAULT NOW()
            );

            CREATE UNIQUE INDEX IF NOT EXISTS idx_guests_email ON guests (email) WHERE email IS NOT NULL AND email != '';
            `
			_, err := tx.Exec(query)
			return err
		},
		Seed: func(tx *sqlx.Tx) error {
			environment := os.Getenv("ENVIRONMENT")
			if environment != meta.ProdEnvironment && environment != meta.StagingEnvironment {
				guests := []domain.Guest{}
				for i := 1; i <= 5; i++ {
					date := gofaker.Date()
					lastVisit := time.Date(2023, time.November, 25, 15, 0, 0, 0, time.UTC)
					guests = append(guests, domain.Guest{
						FirstName:   gofaker.Name(),
						LastName:    gofaker.Name(),
						PhoneNumber: gofaker.PhoneFormatted(),
						Language:    gofaker.Language(),
						Birthdate:   &date,
						LastVisit:   &lastVisit,
						TotalVisits: gofaker.Int32(),
						CurrentMood: gofaker.RandomString([]string{"Negative", "Positive", "Neutral"}),
						TotalSpent:  gofaker.Float32(),
						TotalNoShow: gofaker.Int32(),
						TotalCancel: gofaker.Int32(),
					})
				}

				for _, guest := range guests {
					query := `
                    INSERT INTO guests (first_name, last_name, email, phone_number, language, birthdate, created_at, updated_at) 
                    VALUES (:first_name, :last_name, :email, :phone_number, :language, :birthdate, NOW(), NOW())
`
					_, err := tx.NamedExec(query, guest)
					if err != nil {
						return err
					}
				}
			}

			return nil
		},
		Rollback: func(tx *sqlx.Tx) error {
			query := `DROP TABLE IF EXISTS guests`
			_, err := tx.Exec(query)
			return err
		},
	}
}
