package repository

import (
	"context"
	"net/http"

	"github.com/goplaceapp/goplace-common/pkg/errorhelper"
	guestProto "github.com/goplaceapp/goplace-guest/api/v1"
	"github.com/goplaceapp/goplace-guest/internal/services/reservation-log/adapters/converters"
	"github.com/goplaceapp/goplace-guest/internal/services/reservation-log/domain"
	reservationWaitlistDomain "github.com/goplaceapp/goplace-guest/internal/services/reservation-waitlist/domain"
	"google.golang.org/grpc/status"
)

func (r *ReservationLogRepository) GetAllReservationWaitlistLogs(ctx context.Context, req *guestProto.GetAllReservationWaitlistLogsRequest) (*guestProto.GetAllReservationWaitlistLogsResponse, error) {
	userRepo := r.userClient.Client.UserService.Repository

	if err := r.GetTenantDBConnection(ctx).
		First(
			&reservationWaitlistDomain.ReservationWaitlist{}, "id = ?",
			req.GetReservationWaitlistId(),
		).Error; err != nil {
		return nil, status.Error(http.StatusInternalServerError, errorhelper.ErrReservationNotFound)
	}

	if err := r.GetTenantDBConnection(ctx).
		First(
			&reservationWaitlistDomain.ReservationWaitlist{}, "id = ? AND branch_id = ?",
			req.GetReservationWaitlistId(),
			userRepo.GetCurrentBranchId(ctx),
		).Error; err != nil {
		return nil, status.Error(http.StatusInternalServerError, "You don't have access to this branch")
	}

	var logs []*domain.ReservationWaitlistLog

	if err := r.GetTenantDBConnection(ctx).
		Model(&domain.ReservationWaitlistLog{}).
		Where("reservation_waitlist_id = ?", req.GetReservationWaitlistId()).
		Order("created_at DESC").
		Find(&logs).Error; err != nil {
		return nil, status.Error(http.StatusNotFound, err.Error())
	}

	for _, log := range logs {
		if log.CreatorID == 0 {
			continue
		}

		getCreator, err := userRepo.GetUserProfileByID(ctx, log.CreatorID)
		if err != nil {
			return nil, status.Error(http.StatusNotFound, err.Error())
		}

		log.Creator = getCreator
	}

	return &guestProto.GetAllReservationWaitlistLogsResponse{
		Result: converters.BuildAllReservationWaitlistLogsResponse(r.CommonRepository, ctx, logs),
	}, nil
}
