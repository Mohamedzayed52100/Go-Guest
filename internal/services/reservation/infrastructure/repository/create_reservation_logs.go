package repository

import (
	"context"
	"github.com/goplaceapp/goplace-guest/internal/services/reservation-log/domain"
	"google.golang.org/grpc/status"
	"net/http"
)

func (r *ReservationRepository) CreateReservationLogs(ctx context.Context, logs ...*domain.ReservationLog) ([]*domain.ReservationLog, error) {
	currentUser, err := r.userClient.Client.UserService.Repository.GetLoggedInUser(ctx)
	if err != nil {
		return nil, status.Error(http.StatusInternalServerError, err.Error())
	}

	for _, log := range logs {
		if log.CreatorID == 0 && currentUser != nil {
			log.CreatorID = currentUser.ID
			log.Creator = currentUser
		}

		if err := r.GetTenantDBConnection(ctx).Create(log).Error; err != nil {
			return nil, err
		}
	}

	return logs, nil
}
