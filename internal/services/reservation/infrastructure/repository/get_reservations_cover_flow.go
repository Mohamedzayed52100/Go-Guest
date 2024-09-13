package repository

import (
	"context"
	"github.com/goplaceapp/goplace-common/pkg/meta"
	guestProto "github.com/goplaceapp/goplace-guest/api/v1"
	"github.com/goplaceapp/goplace-guest/pkg/guestservice/converters"
	"github.com/goplaceapp/goplace-guest/internal/services/reservation/domain"
	domain2 "github.com/goplaceapp/goplace-guest/pkg/reservationservice/domain"
	settingsProto "github.com/goplaceapp/goplace-settings/api/v1"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
	"net/http"
	"time"
)

func (r *ReservationRepository) GetReservationsCoverFlow(ctx context.Context, req *guestProto.GetReservationsCoverFlowRequest) (*guestProto.GetReservationsCoverFlowResponse, error) {
	var (
		reservationsStatuses []*domain2.ReservationStatus
		cancelledStatusID    int
		noShowStatusID       int
		coverFlow            = []*domain.CoverFlow{}
	)

	if err := r.GetTenantDBConnection(ctx).
		Where("branch_id = ?", r.userClient.Client.UserService.Repository.GetCurrentBranchId(ctx)).
		Find(&reservationsStatuses).
		Error; err != nil {
		return nil, status.Error(http.StatusInternalServerError, err.Error())
	}

	shift, err := r.shiftClient.Client.ShiftService.Repository.GetShiftByID(ctx, &settingsProto.GetShiftByIDRequest{
		Id: req.GetShiftId(),
	})
	if err != nil {
		return nil, status.Error(http.StatusInternalServerError, err.Error())
	}

	for _, s := range reservationsStatuses {
		if s.Name == meta.Cancelled {
			cancelledStatusID = s.ID
		}
		if s.Name == meta.NoShow {
			noShowStatusID = s.ID
		}
	}

	for i := shift.Result.From; ; {
		var (
			reservations          = []*domain2.Reservation{}
			coverFlowReservations = []*domain.CoverFlowReservation{}
		)

		iTime := i.AsTime()
		newTime := iTime.Add(time.Duration(shift.Result.TimeInterval) * time.Minute)
		i = timestamppb.New(newTime)

		r.GetTenantDBConnection(ctx).
			Where("shift_id = ? AND "+
				"TO_CHAR(date, 'YYYY-MM-DD') = ? AND "+
				"TO_CHAR(time, 'HH24:MI') = ? AND "+
				"status_id NOT IN (?) AND "+
				"seating_area_id IN ? AND "+
				"deleted_at IS NULL",
				req.GetShiftId(),
				req.GetDate(),
				iTime.Format("15:04"),
				[]int{cancelledStatusID, noShowStatusID},
				req.GetSeatingArea(),
			).
			Order("guests_number DESC").
			Find(&reservations)

		for _, r := range reservations {
			var reservationStatus *domain2.ReservationStatus

			for _, s := range reservationsStatuses {
				if int32(s.ID) == r.StatusID {
					reservationStatus = s
					break
				}
			}

			if reservationStatus == nil {
				continue
			}

			coverFlowReservations = append(coverFlowReservations, &domain.CoverFlowReservation{
				ID:           r.ID,
				GuestsNumber: r.GuestsNumber,
				Status: domain.CoverFlowReservationStatus{
					ID:    int32(reservationStatus.ID),
					Name:  reservationStatus.Name,
					Color: reservationStatus.Color,
					Icon:  reservationStatus.Icon,
				},
			})
		}

		coverFlow = append(coverFlow, &domain.CoverFlow{
			Time:         iTime.Format("15:04"),
			Reservations: coverFlowReservations,
		})

		if iTime.Hour()*60+iTime.Minute() == shift.Result.To.AsTime().Hour()*60+shift.Result.To.AsTime().Minute() {
			break
		}
	}

	return &guestProto.GetReservationsCoverFlowResponse{
		Result: converters.BuildCoverFlowResponse(coverFlow),
	}, nil
}
