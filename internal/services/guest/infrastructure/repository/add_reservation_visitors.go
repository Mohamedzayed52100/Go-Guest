package repository

import (
	"context"
	"net/http"

	guestProto "github.com/goplaceapp/goplace-guest/api/v1"
	"github.com/goplaceapp/goplace-guest/pkg/guestservice/domain"
)

func (r *GuestRepository) AddReservationVisitors(ctx context.Context, req *guestProto.AddReservationVisitorsRequest) (*guestProto.AddReservationVisitorsResponse, error) {
	for _, id := range req.GetVisitorIds() {
		r.GetTenantDBConnection(ctx).Create(&domain.ReservationVisitor{
			ReservationID: req.GetReservationId(),
			GuestID:       id,
		})
	}

	return &guestProto.AddReservationVisitorsResponse{
		Code:    http.StatusCreated,
		Message: "Added visitors successfully",
	}, nil
}
