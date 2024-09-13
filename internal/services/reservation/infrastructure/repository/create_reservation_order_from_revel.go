package repository

import (
	"context"
	"errors"
	"net/http"
	"strconv"

	"github.com/goplaceapp/goplace-guest/internal/services/reservation/domain"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"
)

func (r *ReservationRepository) CreateOrUpdateReservationOrderFromRevel(ctx context.Context, reservationId int, orderDetails map[string]interface{}) (*domain.ReservationOrder, error) {
	if orderDetails["discount_amount"] == nil {
		orderDetails["discount_amount"] = 0.0
	}

	record := domain.ReservationOrder{
		ReservationID:  reservationId,
		OrderID:        strconv.Itoa(int(orderDetails["id"].(float64))),
		DiscountAmount: orderDetails["discount_amount"].(float64),
		DiscountReason: orderDetails["discount_reason"].(string),
		PrevailingTax:  orderDetails["prevailing_tax"].(float64),
		Tax:            orderDetails["tax"].(float64),
		Subtotal:       orderDetails["subtotal"].(float64),
		FinalTotal:     orderDetails["final_total"].(float64),
	}

	if err := r.GetTenantDBConnection(ctx).Create(&record).Error; err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return nil, status.Error(http.StatusConflict, "Reservation order already exists")
		}

		return nil, status.Error(http.StatusInternalServerError, err.Error())
	}

	return &record, nil
}

func (r *ReservationRepository) CreateOrUpdateReservationOrderItemsFromRevel(ctx context.Context, reservationOrderId int, orderItems []map[string]interface{}) ([]*domain.ReservationOrderItem, error) {
	var items []*domain.ReservationOrderItem

	for _, orderItem := range orderItems {
		var record *domain.ReservationOrderItem

		if err := r.GetTenantDBConnection(ctx).
			Model(&domain.ReservationOrderItem{}).
			Where("id = ?", orderItem["id"]).
			Create(&domain.ReservationOrderItem{
				ReservationOrderID: reservationOrderId,
				ItemName:           orderItem["product_name_override"].(string),
				Cost:               orderItem["price"].(float64),
				Quantity:           int(orderItem["quantity"].(float64)),
			}).Scan(&record).Error; err != nil {
			if errors.Is(err, gorm.ErrDuplicatedKey) {
				if err := r.GetTenantDBConnection(ctx).
					Model(&domain.ReservationOrderItem{}).
					Where("id = ?", orderItem["id"]).
					Updates(&domain.ReservationOrderItem{
						ItemName: orderItem["product_name_override"].(string),
						Cost:     orderItem["price"].(float64),
						Quantity: orderItem["quantity"].(int),
					}).Scan(&record).Error; err != nil {
					return nil, err
				}
			}
		}

		items = append(items, record)
	}

	return items, nil
}
