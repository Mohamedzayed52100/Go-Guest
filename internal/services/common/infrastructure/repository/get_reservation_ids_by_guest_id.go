package common

import (
	"context"
	"github.com/goplaceapp/goplace-guest/pkg/reservationservice/domain"
	"net/http"

	"google.golang.org/grpc/status"
)

/*
GetReservationIDsByGuestID retrieves all reservation IDs for a specific guest from the database.

Parameters:
- ctx: The context for timeout and cancellation signals.
- guestID: The ID of the guest.

Returns:
- A list of reservation IDs if successful.
- An error if there is an issue retrieving the reservations.
*/
func (r *CommonRepository) GetReservationIDsByGuestID(ctx context.Context, guestID int32) ([]int32, error) {
	var reservationIDs []int32

	if err := r.GetTenantDBConnection(ctx).
		Model(&domain.Reservation{}).
		Where("guest_id = ?", guestID).
		Pluck("id", &reservationIDs).Error; err != nil {
		return nil, status.Errorf(http.StatusNotFound, err.Error())
	}

	return reservationIDs, nil
}
