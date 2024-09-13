package domain

import "time"

type Integration struct {
	ID             int       `db:"id"`
	SystemName     string    `db:"system_name"`
	SystemType     string    `db:"system_type"`
	BaseURL        string    `db:"base_url"`
	CredentialType string    `db:"credential_type"`
	Credentials    string    `db:"credentials"`
	BranchID       int       `db:"branch_id"`
	CreatedAt      time.Time `db:"created_at"`
	UpdatedAt      time.Time `db:"updated_at"`
}
