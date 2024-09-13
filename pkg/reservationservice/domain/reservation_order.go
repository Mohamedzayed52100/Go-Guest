package domain

import (
	"encoding/json"
	"time"
)

type ReservationOrder struct {
	ID             int                     `db:"id"`
	ReservationID  int                     `db:"reservation_id"`
	OrderID        string                  `db:"order_id"`
	Items          []*ReservationOrderItem `gorm:"-"`
	DiscountAmount float64                 `db:"discount_amount"`
	DiscountReason string                  `db:"discount_reason"`
	PrevailingTax  float64                 `db:"prevailing_tax"`
	Tax            float64                 `db:"tax"`
	Subtotal       float64                 `db:"subtotal"`
	FinalTotal     float64                 `db:"final_total"`
	CreatedAt      time.Time               `db:"created_at"`
	UpdatedAt      time.Time               `db:"updated_at"`
}

type ReservationOrderItem struct {
	ID                 int       `db:"id"`
	ReservationOrderID int       `db:"reservation_order_id"`
	ItemName           string    `db:"item_name"`
	Cost               float64   `db:"cost"`
	Quantity           int       `db:"quantity"`
	CreatedAt          time.Time `db:"created_at"`
	UpdatedAt          time.Time `db:"updated_at"`
}

func (ReservationOrderItem) TableName() string {
	return "reservation_order_items"
}

func (order *ReservationOrder) ToString() string {
	result, err := json.Marshal(order)
	if err != nil {
		return ""
	}

	return string(result)
}

func (orderItem *ReservationOrderItem) ToString() string {
	result, err := json.Marshal(orderItem)
	if err != nil {
		return ""
	}

	return string(result)
}
