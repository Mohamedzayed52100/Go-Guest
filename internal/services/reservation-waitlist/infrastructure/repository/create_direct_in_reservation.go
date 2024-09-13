package repository

import (
	"context"
	"net/http"
	"time"

	"github.com/goplaceapp/goplace-common/pkg/errorhelper"
	"github.com/goplaceapp/goplace-common/pkg/meta"
	guestProto "github.com/goplaceapp/goplace-guest/api/v1"
	shiftDomain "github.com/goplaceapp/goplace-settings/pkg/shiftservice/domain"
	"google.golang.org/grpc/status"
)

func (r *ReservationWaitListRepository) CreateDirectInReservation(ctx context.Context, req *guestProto.CreateWaitingReservationRequest) (*guestProto.CreateReservationResponse, error) {
	var branchId int32
	if err := r.GetTenantDBConnection(ctx).
		Model(&shiftDomain.Shift{}).
		Where("id = ?", req.GetParams().GetShiftId()).
		Select("branch_id").
		Scan(&branchId).
		Error; err != nil {
		return nil, status.Error(http.StatusNotFound, errorhelper.ErrShiftNotFound)
	}

	convertedTime := r.CommonRepo.ConvertToLocalTime(ctx, time.Now())

	arrivedStatus, err := r.reservationRepository.GetReservationStatusByName(ctx, meta.Arrived, branchId)
	if err != nil {
		return nil, status.Error(http.StatusInternalServerError, err.Error())
	}

	res, err := r.reservationRepository.CreateReservation(ctx, &guestProto.CreateReservationRequest{
		Params: &guestProto.ReservationParams{
			GuestId:       req.GetParams().GetGuestId(),
			SeatingAreaId: req.GetParams().GetSeatingAreaId(),
			ShiftId:       req.GetParams().GetShiftId(),
			GuestsNumber:  req.GetParams().GetGuestsNumber(),
			BranchId:      branchId,
			ReservedVia:   "Direct in",
			Date:          req.GetParams().GetDate(),
			Time:          convertedTime.Format("15:04"),
			Tags:          req.GetParams().GetTags(),
			StatusId:      int32(arrivedStatus.ID),
		},
	})
	if err != nil {
		return nil, status.Error(http.StatusInternalServerError, err.Error())
	}

	return &guestProto.CreateReservationResponse{
		Result: res.GetResult(),
	}, nil
}
