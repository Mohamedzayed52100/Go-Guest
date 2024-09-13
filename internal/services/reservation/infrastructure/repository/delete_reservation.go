package repository

import (
	"context"
	"net/http"

	guestProto "github.com/goplaceapp/goplace-guest/api/v1"
	reservationDomain "github.com/goplaceapp/goplace-guest/pkg/reservationservice/domain"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (r *ReservationRepository) DeleteReservation(ctx context.Context, req *guestProto.DeleteReservationRequest) (*guestProto.DeleteReservationResponse, error) {
	var reservation *reservationDomain.Reservation

	if err := r.GetTenantDBConnection(ctx).
		Where("id = ?", req.GetId()).
		First(&reservation).
		Error; err != nil {
		return nil, status.Error(http.StatusInternalServerError, err.Error())
	}

	if err := r.GetTenantDBConnection(ctx).
		Delete(&reservationDomain.Reservation{}, "id = ?", req.GetId()).
		Error; err != nil {
		return nil, status.Error(http.StatusInternalServerError, err.Error())
	}

	return &guestProto.DeleteReservationResponse{
		Result: &guestProto.Reservation{
			Id:   req.GetId(),
			Date: timestamppb.New(reservation.Date),
			Shift: &guestProto.ReservationShift{
				Id: reservation.ShiftID,
			},
			Branch: &guestProto.ReservationBranch{
				Id: reservation.BranchID,
			},
		},
	}, nil
}
