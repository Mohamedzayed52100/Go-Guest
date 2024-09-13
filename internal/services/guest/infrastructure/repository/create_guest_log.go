package repository

import (
	"context"
	"net/http"

	logDomain "github.com/goplaceapp/goplace-guest/internal/services/guest-log/domain"
	"google.golang.org/grpc/status"
)

func (r *GuestRepository) CreateGuestLogs(ctx context.Context, logs ...*logDomain.GuestLog) ([]*logDomain.GuestLog, error) {
	currentUser, err := r.UserClient.Client.UserService.Repository.GetLoggedInUser(ctx)
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
