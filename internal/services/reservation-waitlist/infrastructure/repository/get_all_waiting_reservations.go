package repository

import (
	"context"
	"net/http"

	settingsDomain "github.com/goplaceapp/goplace-settings/pkg/shiftservice/domain"

	"github.com/goplaceapp/goplace-common/pkg/errorhelper"
	guestProto "github.com/goplaceapp/goplace-guest/api/v1"
	"github.com/goplaceapp/goplace-guest/internal/services/reservation-waitlist/adapters/converters"
	"github.com/goplaceapp/goplace-guest/internal/services/reservation-waitlist/domain"
	"google.golang.org/grpc/status"
)

func (r *ReservationWaitListRepository) GetAllWaitingReservations(ctx context.Context, req *guestProto.GetWaitingReservationRequest) (*guestProto.GetWaitingReservationsResponse, error) {
	var reservations []*domain.ReservationWaitlist

	if err := r.GetTenantDBConnection(ctx).
		First(
			&settingsDomain.Shift{}, "id = ? AND branch_id = ?",
			req.GetShiftId(),
			r.userClient.Client.UserService.Repository.GetCurrentBranchId(ctx),
		).Error; err != nil {
		return nil, status.Error(http.StatusInternalServerError, errorhelper.ErrShiftNotFound)
	}

	if err := r.GetTenantDBConnection(ctx).
		Model(&reservations).
		Where("shift_id = ? AND (? = '' OR TO_CHAR(date, 'YYYY-MM-DD') = ?)", req.GetShiftId(), req.GetDate(), req.GetDate()).
		Find(&reservations).Error; err != nil {
		return &guestProto.GetWaitingReservationsResponse{
			Result: []*guestProto.ReservationWaitlist{},
		}, nil
	}

	for i, ele := range reservations {
		resData, err := r.GetWaitingReservationData(ctx, ele.ID)
		if err != nil {
			return nil, status.Error(http.StatusInternalServerError, err.Error())
		}

		reservations[i] = resData
	}

	return &guestProto.GetWaitingReservationsResponse{
		Result: converters.BuildAllReservationWaitListsResponse(reservations),
	}, nil
}
