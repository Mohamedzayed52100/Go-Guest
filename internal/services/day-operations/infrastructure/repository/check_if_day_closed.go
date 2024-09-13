package repository

import (
	"context"
	"github.com/goplaceapp/goplace-guest/pkg/reservationservice/domain"
	"net/http"
	"time"

	"github.com/goplaceapp/goplace-common/pkg/meta"
	guestProto "github.com/goplaceapp/goplace-guest/api/v1"
	waitlistDomain "github.com/goplaceapp/goplace-guest/internal/services/reservation-waitlist/domain"
	"google.golang.org/grpc/status"
)

func (r *DayOperationsRepository) CheckIfDayClosed(ctx context.Context, req *guestProto.CheckIfDayClosedRequest) (*guestProto.CheckIfDayClosedResponse, error) {
	var (
		countOfOpenReservations    int64
		countOfWaitingReservations int64
	)

	if req.GetDate() == "" {
		return nil, status.Error(http.StatusBadRequest, "Date is required")
	}

	currentUser, err := r.userClient.Client.UserService.Repository.GetLoggedInUser(ctx)
	if err != nil {
		return nil, status.Error(http.StatusInternalServerError, err.Error())
	}

	noShowStatus, err := r.reservationRepository.GetReservationStatusByName(ctx, meta.NoShow, currentUser.BranchID)
	if err != nil {
		return nil, status.Error(http.StatusInternalServerError, err.Error())
	}
	leftStatus, err := r.reservationRepository.GetReservationStatusByName(ctx, meta.Left, currentUser.BranchID)
	if err != nil {
		return nil, status.Error(http.StatusInternalServerError, err.Error())
	}
	cancelledStatus, err := r.reservationRepository.GetReservationStatusByName(ctx, meta.Cancelled, currentUser.BranchID)
	if err != nil {
		return nil, status.Error(http.StatusInternalServerError, err.Error())
	}

	excludedStatuses := []int{noShowStatus.ID, leftStatus.ID, cancelledStatus.ID}
	if err := r.GetTenantDBConnection(ctx).
		Model(&domain.Reservation{}).
		Where("date = ? AND branch_id = ? AND status_id NOT IN (?)",
			req.GetDate(),
			currentUser.BranchID,
			excludedStatuses,
		).Count(&countOfOpenReservations).Error; err != nil {
		return nil, status.Error(http.StatusInternalServerError, err.Error())
	}

	if err := r.GetTenantDBConnection(ctx).
		Model(&waitlistDomain.ReservationWaitlist{}).
		Where("branch_id = ? AND created_at <= ?",
			currentUser.BranchID,
			time.Now().Format("2006-01-02"),
		).
		Count(&countOfWaitingReservations).Error; err != nil {
		return nil, status.Error(http.StatusInternalServerError, err.Error())
	}

	if countOfOpenReservations != 0 || countOfWaitingReservations != 0 {
		return &guestProto.CheckIfDayClosedResponse{
			Closed: false,
		}, nil
	}

	return &guestProto.CheckIfDayClosedResponse{
		Closed: true,
	}, nil

}
