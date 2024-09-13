package repository

import (
	"context"
	"net/http"

	guestProto "github.com/goplaceapp/goplace-guest/api/v1"
	"github.com/goplaceapp/goplace-guest/internal/services/reservation-feedback-comment/adapters/converters"
	"github.com/goplaceapp/goplace-guest/internal/services/reservation-feedback-comment/domain"
	feedbackDomain "github.com/goplaceapp/goplace-guest/internal/services/reservation-feedback/domain"
	roleDomain "github.com/goplaceapp/goplace-user/pkg/roleservice/domain"
	"google.golang.org/grpc/status"
)

func (r *ReservationFeedbackCommentRepository) CreateReservationFeedbackComment(ctx context.Context, req *guestProto.CreateReservationFeedbackCommentRequest) (*guestProto.CreateReservationFeedbackCommentResponse, error) {
	userRepo := r.userClient.Client.UserService.Repository

	currentUser, err := userRepo.GetLoggedInUser(ctx)
	if err != nil {
		return nil, status.Error(http.StatusInternalServerError, err.Error())
	}

	comment := &domain.ReservationFeedbackComment{
		ReservationFeedbackID: req.GetParams().GetReservationFeedbackId(),
		Comment:               req.GetParams().GetComment(),
		CreatorID:             currentUser.ID,
	}

	if err := r.GetTenantDBConnection(ctx).Create(&comment).Error; err != nil {
		return nil, status.Error(http.StatusInternalServerError, err.Error())
	}

	var ReservationFeedback *feedbackDomain.ReservationFeedback
	if err := r.GetTenantDBConnection(ctx).First(&ReservationFeedback, "id = ?", req.GetParams().GetReservationFeedbackId()).Error; err != nil {
		return nil, status.Error(http.StatusInternalServerError, err.Error())
	}

	var role *roleDomain.Role
	if err := r.GetTenantDBConnection(ctx).First(&role, "id = ?", currentUser.RoleID).Error; err != nil {
		return nil, status.Error(http.StatusInternalServerError, err.Error())
	}
	currentUser.Role = role
	comment.Creator = currentUser

	return &guestProto.CreateReservationFeedbackCommentResponse{
		Result: converters.BuildReservationFeedbackCommentResponse(comment),
	}, nil
}
