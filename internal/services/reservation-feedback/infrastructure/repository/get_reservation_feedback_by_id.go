package repository

import (
	"context"
	"net/http"

	reservationDomain "github.com/goplaceapp/goplace-guest/pkg/reservationservice/domain"

	"github.com/goplaceapp/goplace-common/pkg/meta"
	"github.com/goplaceapp/goplace-common/pkg/utils"
	guestProto "github.com/goplaceapp/goplace-guest/api/v1"
	"github.com/goplaceapp/goplace-guest/internal/services/reservation-feedback/domain"
	"github.com/goplaceapp/goplace-guest/pkg/guestservice/converters"
	guestDomain "github.com/goplaceapp/goplace-guest/pkg/guestservice/domain"
	"google.golang.org/grpc/status"
)

func (r *ReservationFeedbackRepository) GetReservationFeedbackByID(ctx context.Context, req *guestProto.GetReservationFeedbackByIDRequest) (*guestProto.GetReservationFeedbackByIDResponse, error) {
	var (
		feedback *domain.ReservationFeedback
		err      error
	)

	if err := r.GetTenantDBConnection(ctx).
		Model(feedback).
		Joins("JOIN reservations ON reservations.id = reservation_feedbacks.reservation_id").
		Select("reservation_feedbacks.*").
		First(&feedback, "reservation_feedbacks.id = ? AND reservation_id = ? AND branch_id = ?",
			req.GetFeedbackId(),
			req.GetReservationId(),
			r.userClient.Client.UserService.Repository.GetCurrentBranchId(ctx)).Error; err != nil {
		return nil, status.Error(http.StatusInternalServerError, "Reservation not found")
	}

	feedback, err = r.GetReservationFeedbackData(ctx, feedback)
	if err != nil {
		return nil, status.Error(http.StatusInternalServerError, err.Error())
	}

	return &guestProto.GetReservationFeedbackByIDResponse{
		Result: converters.BuildReservationFeedbackResponse(feedback),
	}, nil
}

func (r *ReservationFeedbackRepository) GetReservationFeedbackData(ctx context.Context, feedback *domain.ReservationFeedback) (*domain.ReservationFeedback, error) {
	var (
		err             error
		reservationData *reservationDomain.Reservation
	)

	feedback.Sections, err = r.GetAllReservationFeedbackSectionsByID(ctx, feedback.ID)
	if err != nil {
		return nil, err
	}

	reservationData, err = r.CommonRepo.GetReservationByID(ctx, feedback.ReservationID)
	if err != nil {
		return nil, err
	}

	if feedback.StatusID != 0 {
		feedback.Status = FeedbackStatuses[feedback.StatusID-1]
	}

	if feedback.Rate <= 3 && feedback.Status != meta.Solved {
		for j, v := range FeedbackStatuses {
			if utils.CompareStr(v, meta.Pending) {
				feedback.StatusID = int32(j)
				feedback.Status = v
				break
			}
		}
	}

	if feedback.SolutionID != 0 {
		var solution *domain.ReservationFeedbackSolution

		if err := r.GetTenantDBConnection(ctx).
			First(&solution, "id = ?", feedback.SolutionID).
			Error; err != nil {
			return nil, err
		}

		currentUser, err := r.userClient.Client.UserService.Repository.GetUserProfileByID(ctx, solution.CreatorID)
		if err != nil {
			return nil, err
		}

		solution.Creator = currentUser
		feedback.Solution = solution
	} else {
		feedback.Solution = nil
	}

	var primaryGuest *guestDomain.Guest
	r.GetTenantDBConnection(ctx).First(&primaryGuest, "id = ?", reservationData.GuestID)

	feedback.Guest = primaryGuest
	feedback.Reservation = reservationData

	return feedback, nil
}
