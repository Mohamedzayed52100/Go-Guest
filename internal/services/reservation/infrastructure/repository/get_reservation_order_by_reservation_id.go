package repository

import (
	"context"
	guestProto "github.com/goplaceapp/goplace-guest/api/v1"
	"github.com/goplaceapp/goplace-guest/internal/services/reservation/domain"
	"github.com/goplaceapp/goplace-guest/pkg/guestservice/converters"
	domain2 "github.com/goplaceapp/goplace-guest/pkg/reservationservice/domain"
	"google.golang.org/grpc/status"
	"net/http"
)

func (r *ReservationRepository) GetReservationOrderByReservationID(ctx context.Context, req *guestProto.GetReservationOrderByReservationIDRequest) (*guestProto.GetReservationOrderByReservationIDResponse, error) {
	var (
		order *domain.ReservationOrder
		items = []*domain.ReservationOrderItem{}
	)

	if err := r.GetTenantDBConnection(ctx).
		First(&domain2.Reservation{}, "id = ? AND branch_id IN (?)",
			req.GetReservationId(), r.userClient.Client.UserService.Repository.GetAllUserBranchesIDs(ctx)).
		Error; err != nil {
		return nil, status.Error(http.StatusNotFound, "Reservation not found")
	}

	if err := r.GetTenantDBConnection(ctx).
		First(&order, "reservation_id = ?", req.GetReservationId()).
		Error; err != nil {
		return &guestProto.GetReservationOrderByReservationIDResponse{}, nil
	}

	r.GetTenantDBConnection(ctx).Find(&items, "reservation_order_id = ?", order.ID)
	order.Items = items

	return &guestProto.GetReservationOrderByReservationIDResponse{
		Result: converters.BuildReservationOrderResponse(order),
	}, nil
}
