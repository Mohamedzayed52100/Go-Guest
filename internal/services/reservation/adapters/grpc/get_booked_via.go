package grpc

import (
	"context"

	guestProto "github.com/goplaceapp/goplace-guest/api/v1"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (r *ReservationServiceServer) GetAllBookedVia(ctx context.Context, req *emptypb.Empty) (*guestProto.GetAllBookedViaResponse, error) {
	res, err := r.reservationService.GetAllBookedVia(ctx, req)
	if err != nil {
		return nil, err
	}

	return res, nil
}
