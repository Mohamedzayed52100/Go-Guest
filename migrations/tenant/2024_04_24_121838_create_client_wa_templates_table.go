package migrations

import (
	"github.com/goplaceapp/goplace-common/pkg/dbhelper"
	"github.com/jmoiron/sqlx"
)

func CreateClientWaTemplatesTable() dbhelper.SqlxMigration {
	return dbhelper.SqlxMigration{
		ID: "2024_04_24_121838_create_client_wa_templates_table",
		Migrate: func(tx *sqlx.Tx) error {
			_, err := tx.Exec(`
                CREATE TABLE IF NOT EXISTS whatsapp_templates (
                    id SERIAL PRIMARY KEY,
                    branch_id INT NOT NULL,
                    template_name TEXT NOT NULL,
                    template_type TEXT NOT NULL,
                    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
                    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
                    deleted_at TIMESTAMP WITH TIME ZONE,
                    
                    CONSTRAINT fk_whatsapp_templates_branch_id FOREIGN KEY (branch_id) REFERENCES branches (id),
                    CONSTRAINT uq_whatsapp_templates_branch_id_name UNIQUE (branch_id, template_name, template_type)
                );
            `)
			return err
		},
		Rollback: func(tx *sqlx.Tx) error {
			_, err := tx.Exec(`DROP TABLE whatsapp_templates;`)
			return err
		},
	}
}
