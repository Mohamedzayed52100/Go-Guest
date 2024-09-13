package database

import (
	"database/sql"
	"fmt"

	"github.com/goplaceapp/goplace-guest/internal/services/reservation-special-occasion/infrastructure/database/seeder"
)

type Executable interface {
	Exec(query string, args ...interface{}) (sql.Result, error)
	Query(query string, args ...interface{}) (*sql.Rows, error)
	QueryRow(query string, args ...interface{}) *sql.Row
}

func SyncSeeder(tx Executable) error {
	var err error
	fmt.Println("Syncing seeders...")

	err = seeder.SpecialOccasionsSeeder(tx)
	if err != nil {
		fmt.Println("Error syncing seed: Special occasions")
		return err
	}
	fmt.Println("Synced Seed: Special occasions")

	return nil
}
