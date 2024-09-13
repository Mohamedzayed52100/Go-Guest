package domain

type InvoiceItem struct {
	ID       int32   `db:"id"`
	Name     string  `db:"name"`
	Price    float32 `db:"price"`
	Quantity int32   `db:"quantity"`
}
