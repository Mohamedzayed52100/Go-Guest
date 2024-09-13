package repository

import (
	"context"
	"github.com/goplaceapp/goplace-guest/pkg/reservationservice/domain"
	"net/http"

	guestProto "github.com/goplaceapp/goplace-guest/api/v1"
	"github.com/goplaceapp/goplace-guest/pkg/guestservice/converters"
	"google.golang.org/grpc/status"
)

func (r *GuestRepository) GetGuestStatistics(ctx context.Context, req *guestProto.GetGuestStatisticsRequest) (*guestProto.GetGuestStatisticsResponse, error) {
	var (
		err                error
		totalReservations  int64
		totalSpent         float32
		publicSatisfaction string
	)

	r.GetTenantDBConnection(ctx).
		Find(&domain.Reservation{}, "guest_id = ?", req.GuestId).
		Count(&totalReservations)

	r.GetTenantDBConnection(ctx).
		Table("reservations").
		Joins("JOIN reservation_orders ON reservation_orders.reservation_id = reservations.id").
		Joins("JOIN reservation_order_items ON reservation_order_items.reservation_order_id = reservation_orders.id").
		Where("reservations.guest_id = ?", req.GuestId).
		Select("COALESCE(reservation_orders.final_total)").
		Scan(&totalSpent)

	publicSatisfaction, err = r.CommonRepository.GetGuestCurrentMood(ctx, req.GetGuestId())
	if err != nil {
		return nil, status.Error(http.StatusInternalServerError, err.Error())
	}

	return &guestProto.GetGuestStatisticsResponse{
		Result: converters.BuildGuestStatisticsResponse(totalReservations, totalSpent, publicSatisfaction),
	}, nil
}
