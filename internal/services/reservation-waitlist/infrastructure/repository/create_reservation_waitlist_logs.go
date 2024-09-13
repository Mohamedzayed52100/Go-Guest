package repository

import (
	"context"
	"net/http"

	"github.com/goplaceapp/goplace-guest/internal/services/reservation-waitlist/domain"
	"google.golang.org/grpc/status"
)

func (r *ReservationWaitListRepository) CreateReservationWaitListLogs(ctx context.Context, logs ...*domain.ReservationWaitlistLog) ([]*domain.ReservationWaitlistLog, error) {
	currentUser, err := r.userClient.Client.UserService.Repository.GetLoggedInUser(ctx)
	if err != nil {
		return nil, status.Error(http.StatusInternalServerError, err.Error())
	}

	for _, log := range logs {
		log.CreatorID = currentUser.ID
		log.Creator = currentUser

		if err := r.GetTenantDBConnection(ctx).Create(log).Error; err != nil {
			return nil, err
		}
	}

	return logs, nil
}
