package migrations

import (
    "github.com/goplaceapp/goplace-common/pkg/dbhelper"
    "github.com/jmoiron/sqlx"
)

func AlterPaymentRequests() dbhelper.SqlxMigration {
    return dbhelper.SqlxMigration{
        ID: "2024_06_24_154046_alter_payment_requests",
        Migrate: func(tx *sqlx.Tx) error {
            query := `
            ALTER TABLE payment_requests
            DROP COLUMN IF EXISTS note,
            DROP COLUMN IF EXISTS special_request_id;
            `
            _, err := tx.Exec(query)
            return err
        },
    }
}
