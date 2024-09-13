package repository

import (
	"context"

	guestProto "github.com/goplaceapp/goplace-guest/api/v1"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (r *ReservationRepository) GetAllBookedVia(ctx context.Context, req *emptypb.Empty) (*guestProto.GetAllBookedViaResponse, error) {
	var result []string
	r.GetTenantDBConnection(ctx).
		Table("reservations").
		Select("LOWER(reserved_via)").
		Group("reserved_via").
		Distinct("LOWER(reserved_via)").
		Pluck("LOWER(reserved_via)", &result)

	return &guestProto.GetAllBookedViaResponse{
		Result: result,
	}, nil
}
