package repository

import (
	"context"
	reservationDomain "github.com/goplaceapp/goplace-guest/pkg/reservationservice/domain"
	"net/http"

	guestProto "github.com/goplaceapp/goplace-guest/api/v1"
	"github.com/goplaceapp/goplace-guest/internal/services/reservation-log/adapters/converters"
	"github.com/goplaceapp/goplace-guest/internal/services/reservation-log/domain"
	"google.golang.org/grpc/status"
)

func (r *ReservationLogRepository) GetAllReservationLogs(ctx context.Context, req *guestProto.GetAllReservationLogsRequest) (*guestProto.GetAllReservationLogsResponse, error) {
	userRepo := r.userClient.Client.UserService.Repository

	if err := r.GetTenantDBConnection(ctx).
		First(
			&reservationDomain.Reservation{}, "id = ? AND branch_id IN (?)",
			req.GetReservationId(),
			userRepo.GetAllUserBranchesIDs(ctx),
		).Error; err != nil {
		return nil, status.Error(http.StatusInternalServerError, "Reservation not found")
	}

	var logs []*domain.ReservationLog

	if err := r.GetTenantDBConnection(ctx).
		Model(&domain.ReservationLog{}).
		Where("reservation_id = ?", req.GetReservationId()).
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

	return &guestProto.GetAllReservationLogsResponse{
		Result: converters.BuildAllReservationLogsResponse(r.CommonRepository, ctx, logs),
	}, nil
}
