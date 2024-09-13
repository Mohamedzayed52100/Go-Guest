package domain

type WhatsappTemplate struct {
	ID           int    `db:"id"`
	BranchID     int    `db:"branch_id"`
	TemplateName string `db:"template_name"`
	TemplateType string `db:"template_type"`
	CreatedAt    string `db:"created_at"`
	UpdatedAt    string `db:"updated_at"`
	DeleteAt     string `db:"deleted_at"`
}