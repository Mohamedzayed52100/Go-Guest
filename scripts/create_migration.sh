#!/bin/bash

# Check if the correct number of arguments was provided
if [ "$#" -ne 2 ]; then
    echo "Usage: $0 migration_name migration_database"
    exit 1
fi

MIGRATION_NAME=$1
MIGRATION_DATABASE=$2

# Use regex to extract the action from the migration name
if [[ $MIGRATION_NAME =~ ^(create|alter|drop|add)_ ]]; then
    ACTION="${BASH_REMATCH[1]}"
else
    echo "Migration name must start with 'create_', 'alter_', 'drop_', or 'add_'."
    exit 2
fi

# Generate a unique ID based on the current date and time
DATE=$(date +%Y_%m_%d_%H%M%S)
MIGRATION_ID="${DATE}_${MIGRATION_NAME}"

# Define the directory based on migration database
if [ "$MIGRATION_DATABASE" = "shared" ]; then
    DIRECTORY="migrations/shared"
elif [ "$MIGRATION_DATABASE" = "tenant" ]; then
    DIRECTORY="migrations/tenant"
else
    echo "Invalid migration database: $MIGRATION_DATABASE. Use 'shared' or 'tenant'."
    exit 3
fi

# Ensure the directory exists
mkdir -p $DIRECTORY

# Path to the new migration file
FILE_PATH="$DIRECTORY/$MIGRATION_ID.go"

# Method name should be PascalCase and without spaces or special characters
METHOD_NAME=$(echo $MIGRATION_NAME | sed -r 's/[^a-zA-Z0-9]+/_/g' | sed -r 's/(^|_)([a-z])/\U\2/g')

# Determine SQL query based on the action
case $ACTION in
    create)
        TABLE_NAME=$(echo "$MIGRATION_NAME" | sed -r 's/^create_//g' | sed -r 's/_table$//g' | sed -r 's/[^a-zA-Z0-9]+/_/g' | tr '[:upper:]' '[:lower:]')
        SQL_QUERY="CREATE TABLE IF NOT EXISTS ${TABLE_NAME} (
        );"
        ROLLBACK_QUERY="DROP TABLE IF EXISTS ${TABLE_NAME};"
        ;;
    alter)
        TABLE_NAME=$(echo "$MIGRATION_NAME" | sed -r 's/^alter_//g' | sed -r 's/_table$//g' | sed -r 's/[^a-zA-Z0-9]+/_/g' | tr '[:upper:]' '[:lower:]')
        SQL_QUERY="ALTER TABLE ${TABLE_NAME}"
        ROLLBACK_QUERY=""
        ;;
    drop)
        TABLE_NAME=$(echo "$MIGRATION_NAME" | sed -r 's/^drop_//g' | sed -r 's/_table$//g' | sed -r 's/[^a-zA-Z0-9]+/_/g' | tr '[:upper:]' '[:lower:]')
        SQL_QUERY="DROP TABLE IF EXISTS ${TABLE_NAME};"
        ROLLBACK_QUERY="CREATE TABLE ${TABLE_NAME} (
        );"
        ;;
    add)
        COLUMN_NAME=$(echo "$MIGRATION_NAME" | sed -r 's/^add_//g' | sed -r 's/_to_.*$//g' | sed -r 's/[^a-zA-Z0-9]+/_/g' | tr '[:upper:]' '[:lower:]')
        TABLE_NAME=$(echo "$MIGRATION_NAME" | sed -r 's/^add_.*_to_//g' | sed -r 's/_table$//g' | sed -r 's/[^a-zA-Z0-9]+/_/g' | tr '[:upper:]' '[:lower:]')
        SQL_QUERY="ALTER TABLE ${TABLE_NAME} ADD COLUMN ${COLUMN_NAME} BOOLEAN DEFAULT FALSE"
        ROLLBACK_QUERY="ALTER TABLE ${TABLE_NAME} DROP COLUMN ${COLUMN_NAME}"
        ;;
    *)
        echo "Invalid action derived from name: $ACTION."
        exit 4
        ;;
esac

# Write the migration template to the new file
cat <<EOF > $FILE_PATH
package migrations

import (
    "github.com/goplaceapp/goplace-common/pkg/dbhelper"
    "github.com/jmoiron/sqlx"
)

func ${METHOD_NAME}() dbhelper.SqlxMigration {
    return dbhelper.SqlxMigration{
        ID: "$MIGRATION_ID",
        Migrate: func(tx *sqlx.Tx) error {
            query := \`
            $SQL_QUERY
            \`
            _, err := tx.Exec(query)
            return err
        },
        Seed: func(tx *sqlx.Tx) error {
            return nil
        },
        Rollback: func(tx *sqlx.Tx) error {
            query := "$ROLLBACK_QUERY"
            _, err := tx.Exec(query)
            return err
        },
    }
}
EOF

echo "Migration created at $FILE_PATH"

exit 0