package database

import (
	"github.com/goplaceapp/goplace-common/pkg/dbhelper"
	"github.com/goplaceapp/goplace-guest/migrations"
	"gorm.io/gorm"
)

type Migrator interface {
	SharedMigrationsUp(db *gorm.DB) error
	TenantMigrationsUp(db *gorm.DB) error
}

func SharedMigrationsUp(db *gorm.DB) error {
	sqlMigration := dbhelper.Sqlx{
		Migrations: migrations.SharedMigrations,
	}

	conn, _ := db.DB()
	err := sqlMigration.Migrate(conn, "postgres")
	if err != nil {
		return err
	}

	return nil
}

func TenantMigrationsUp(db *gorm.DB) error {
	sqlMigration := dbhelper.Sqlx{
		Migrations: migrations.TenantMigrations,
	}

	conn, _ := db.DB()
	err := sqlMigration.Migrate(conn, "postgres")
	if err != nil {
		return err
	}

	if err := SyncSeeder(conn); err != nil {
		return err
	}

	if err := SetupReservationListener(conn); err != nil {
		return err
	}

	if err := SetupReservationsWaitlistListener(conn); err != nil {
		return err
	}

	return nil
}
