package seeder

import (
	"database/sql"

	"github.com/goplaceapp/goplace-common/pkg/meta"
	userDomain "github.com/goplaceapp/goplace-user/pkg/userservice/domain"
)

type Executable interface {
	Exec(query string, args ...interface{}) (sql.Result, error)
	Query(query string, args ...interface{}) (*sql.Rows, error)
	QueryRow(query string, args ...interface{}) *sql.Row
}

func SpecialOccasionsSeeder(tx Executable) error {
	branches, err := GetAllBranches(tx)
	if err != nil {
		return err
	}

	query := `
    INSERT INTO special_occasions (name, color, icon, branch_id, created_at, updated_at)
    VALUES ($1, $2, $3, $4, NOW(), NOW())
	ON CONFLICT (branch_id, name)
    DO UPDATE SET
        color = $2,
        icon = $3,
        updated_at = NOW()
    `

	for _, branch := range branches {
		for _, sp := range meta.SpecialOccasions {
			_, err := tx.Exec(query, sp.Name, sp.Color, sp.Icon, branch.ID)
			if err != nil {
				return err
			}
		}
	}

	return nil

}

func GetAllBranches(tx Executable) ([]userDomain.Branch, error) {
	var branches []userDomain.Branch
	rows, err := tx.Query("SELECT id, name FROM branches")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var branch userDomain.Branch
		if err := rows.Scan(&branch.ID, &branch.Name); err != nil {
			return nil, err
		}
		branches = append(branches, branch)
	}

	return branches, nil
}
