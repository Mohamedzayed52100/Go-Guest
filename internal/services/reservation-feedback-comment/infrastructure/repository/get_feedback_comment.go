package repository

import (
	"context"
	reservationDomain "github.com/goplaceapp/goplace-guest/pkg/reservationservice/domain"
	"net/http"

	guestProto "github.com/goplaceapp/goplace-guest/api/v1"
	"github.com/goplaceapp/goplace-guest/internal/services/reservation-feedback-comment/adapters/converters"
	"github.com/goplaceapp/goplace-guest/internal/services/reservation-feedback-comment/domain"
	"google.golang.org/grpc/status"
)

func (r *ReservationFeedbackCommentRepository) GetAllReservationFeedbackComments(ctx context.Context, req *guestProto.GetAllReservationFeedbackCommentsRequest) (*guestProto.GetAllReservationFeedbackCommentsResponse, error) {
	userRepo := r.userClient.Client.UserService.Repository

	if err := r.GetTenantDBConnection(ctx).
		Model(reservationDomain.Reservation{}).
		Joins("JOIN reservation_feedbacks ON reservations.id = reservation_feedbacks.reservation_id").
		Joins("JOIN reservation_feedback_comments ON reservation_feedbacks.id = reservation_feedback_comments.reservation_feedback_id").
		Where("branch_id = ?", userRepo.GetCurrentBranchId(ctx)).
		Select("reservations.*").Scan(&reservationDomain.Reservation{}).Error; err != nil {
		return nil, status.Error(http.StatusInternalServerError, "You don't have access to this branch")
	}

	var comments []*domain.ReservationFeedbackComment
	r.GetTenantDBConnection(ctx).Order("id DESC").Find(&comments, "reservation_feedback_id = ?", req.GetReservationFeedbackId())

	currentUser, err := userRepo.GetLoggedInUser(ctx)
	if err != nil {
		return nil, status.Error(http.StatusInternalServerError, err.Error())
	}
	for i := range comments {
		comments[i].Creator = currentUser
	}

	return &guestProto.GetAllReservationFeedbackCommentsResponse{
		Result: converters.BuildAllReservationFeedbackCommentsResponse(comments),
	}, nil
}
