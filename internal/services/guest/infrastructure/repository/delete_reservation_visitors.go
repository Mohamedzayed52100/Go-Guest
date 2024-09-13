package repository

import (
	"context"
	"net/http"

	guestProto "github.com/goplaceapp/goplace-guest/api/v1"
	"github.com/goplaceapp/goplace-guest/pkg/guestservice/domain"
)

func (r *GuestRepository) DeleteReservationVisitors(ctx context.Context, req *guestProto.DeleteReservationVisitorsRequest) (*guestProto.DeleteReservationVisitorsResponse, error) {
	for _, id := range req.GetVisitorIds() {
		r.GetTenantDBConnection(ctx).
			Delete(&domain.ReservationVisitor{}, "reservation_id = ? AND guest_id = ?", req.GetReservationId(), id)
	}

	return &guestProto.DeleteReservationVisitorsResponse{
		Code:    http.StatusOK,
		Message: "Deleted visitors successfully",
	}, nil
}
