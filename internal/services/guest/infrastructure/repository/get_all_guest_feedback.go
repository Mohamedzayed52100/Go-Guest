package repository

import (
	"context"
	"net/http"

	guestProto "github.com/goplaceapp/goplace-guest/api/v1"
	"github.com/goplaceapp/goplace-guest/internal/services/reservation-feedback/domain"
	reservationFeedbackDomain "github.com/goplaceapp/goplace-guest/internal/services/reservation-feedback/domain"
	"github.com/goplaceapp/goplace-guest/pkg/guestservice/converters"
	"google.golang.org/grpc/status"
)

/*
GetAllGuestFeedback retrieves all feedback for a specific guest from the database.

Parameters:
- ctx: The context for timeout and cancellation signals.
- req: The request containing the ID of the guest.

The method:
1. Retrieves all reservation IDs for the guest.
2. For each reservation, retrieves the reservation details and feedback.
3. Appends the feedback to a list if it exists.

Returns:
- A response containing all the guest's feedback if successful.
- An error if there is an issue retrieving the reservations or feedback.
*/

func (r *GuestRepository) GetAllGuestFeedback(ctx context.Context, req *guestProto.GetAllGuestFeedbackRequest) (*guestProto.GetAllGuestFeedbackResponse, error) {
	var reservationFeedbacks []*reservationFeedbackDomain.ReservationFeedback

	guestReservations, err := r.CommonRepository.GetReservationIDsByGuestID(ctx, req.GuestId)
	if err != nil {
		return nil, err
	}

	for _, guestReservation := range guestReservations {
		var feedback *domain.ReservationFeedback

		if err := r.GetTenantDBConnection(ctx).
			Find(&feedback, "reservation_id = ?", guestReservation).
			Error; err != nil {
			return nil, status.Error(http.StatusInternalServerError, err.Error())
		}

		feedback, _ = r.ReservationFeedbackRepository.GetReservationFeedbackData(ctx, feedback)

		if feedback != nil {
			reservationFeedbacks = append(reservationFeedbacks, feedback)
		}
	}

	return &guestProto.GetAllGuestFeedbackResponse{
		Result: converters.BuildAllReservationFeedbacksResponse(reservationFeedbacks),
	}, nil

}
