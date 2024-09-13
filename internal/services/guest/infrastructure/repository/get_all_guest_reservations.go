package repository

import (
	"context"
	"github.com/goplaceapp/goplace-guest/pkg/reservationservice/domain"
	"net/http"

	"github.com/goplaceapp/goplace-common/pkg/errorhelper"
	guestProto "github.com/goplaceapp/goplace-guest/api/v1"
	common "github.com/goplaceapp/goplace-guest/internal/services/common/converters"
	"github.com/goplaceapp/goplace-guest/pkg/guestservice/converters"
	"google.golang.org/grpc/status"
)

func (r *GuestRepository) GetAllGuestReservations(ctx context.Context, req *guestProto.GetAllGuestReservationsRequest) (*guestProto.GetAllGuestReservationsResponse, error) {
	var (
		reservations []*domain.Reservation
		totalRows    int64
		offset       = (req.GetParams().GetCurrentPage() - 1) * req.GetParams().GetPerPage()
	)

	if err := r.GetTenantDBConnection(ctx).
		Model(&domain.Reservation{}).
		Where("guest_id = ? AND branch_id IN (?)",
			req.GetGuestId(),
			r.UserClient.Client.UserService.Repository.GetAllUserBranchesIDs(ctx),
		).
		Count(&totalRows).Error; err != nil {
		return nil, status.Error(http.StatusNotFound, errorhelper.ErrNoReservationsFound)
	}

	if err := r.GetTenantDBConnection(ctx).
		Model(&domain.Reservation{}).
		Offset(int(offset)).
		Limit(int(req.GetParams().GetPerPage())).
		Where("guest_id = ? AND branch_id IN (?)",
			req.GetGuestId(),
			r.UserClient.Client.UserService.Repository.GetAllUserBranchesIDs(ctx),
		).
		Order("created_at DESC").
		Find(&reservations).Error; err != nil {
		return nil, status.Error(http.StatusNotFound, errorhelper.ErrNoReservationsFound)
	}

	for i := range reservations {
		getReservationByID, err := r.CommonRepository.GetReservationByID(ctx, reservations[i].ID)
		if err != nil {
			return nil, err
		}

		reservations[i] = getReservationByID
	}

	return &guestProto.GetAllGuestReservationsResponse{
		Pagination: common.BuildPaginationResponse(req.GetParams(), int32(totalRows), offset),
		Result:     converters.BuildAllReservationsResponse(reservations),
	}, nil
}
