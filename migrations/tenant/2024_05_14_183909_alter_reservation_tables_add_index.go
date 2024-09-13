package migrations

import (
    "github.com/goplaceapp/goplace-common/pkg/dbhelper"
    "github.com/jmoiron/sqlx"
)

func AlterReservationTablesAddIndex() dbhelper.SqlxMigration {
    return dbhelper.SqlxMigration{
        ID: "2024_05_14_183909_alter_reservation_tables_add_index",
        Migrate: func(tx *sqlx.Tx) error {
            query := `
            CREATE INDEX IF NOT EXISTS reservation_tables_idx ON reservation_tables (table_id);
            `
            _, err := tx.Exec(query)
            return err
        },
    }
}
