package repository

import (
	"context"
	"net/http"

	guestProto "github.com/goplaceapp/goplace-guest/api/v1"
	"github.com/goplaceapp/goplace-guest/internal/services/guest-log/adapters/converters"
	"github.com/goplaceapp/goplace-guest/internal/services/guest-log/domain"
	"google.golang.org/grpc/status"
)

/*
GetAllGuestLogs retrieves all logs history for a guest.

Parameters:
- ctx: The context for timeout and cancellation signals.
- req: The request containing the ID of the guest.

The method:
1. Retrieves all logs for the guest, ordered by creation date in descending order.
2. For each log, if the creator ID is not zero, retrieves the creator's profile and assigns it to the log.

Returns:
- A response containing all the guest's logs if successful.
- An error if there is an issue retrieving the logs or the creator's profile.
*/
func (r *GuestLogRepository) GetAllGuestLogs(ctx context.Context, req *guestProto.GetAllGuestLogsRequest) (*guestProto.GetAllGuestLogsResponse, error) {
	var logs []*domain.GuestLog

	if err := r.GetTenantDBConnection(ctx).
		Model(&domain.GuestLog{}).
		Where("guest_id = ?", req.GetGuestId()).
		Order("created_at desc").
		Find(&logs).Error; err != nil {
		return nil, status.Error(http.StatusNotFound, err.Error())
	}

	for _, log := range logs {
		if log.CreatorID != 0 {
			getCreator, err := r.userClient.Client.UserService.Repository.GetUserProfileByID(ctx, log.CreatorID)
			if err != nil {
				return nil, status.Error(http.StatusNotFound, err.Error())
			}

			log.Creator = getCreator
		}
	}

	return &guestProto.GetAllGuestLogsResponse{
		Result: converters.BuildAllGuestLogsResponse(logs),
	}, nil
}
