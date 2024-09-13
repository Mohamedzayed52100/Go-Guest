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

func (r *ReservationFeedbackCommentRepository) UpdateReservationFeedbackComment(ctx context.Context, req *guestProto.UpdateReservationFeedbackCommentRequest) (*guestProto.UpdateReservationFeedbackCommentResponse, error) {
	updates := map[string]interface{}{}
	userRepo := r.userClient.Client.UserService.Repository

	if req.GetParams().GetComment() != "" {
		updates["comment"] = req.GetParams().GetComment()
	}
	if req.GetParams().GetReservationFeedbackId() != 0 {
		updates["reservation_feedback_id"] = req.GetParams().GetReservationFeedbackId()
	}

	if err := r.GetTenantDBConnection(ctx).Model(&domain.ReservationFeedbackComment{}).Where("id =?", req.GetParams().GetId()).Updates(updates).Error; err != nil {
		return nil, status.Error(http.StatusInternalServerError, err.Error())
	}

	var comment *domain.ReservationFeedbackComment
	if err := r.GetTenantDBConnection(ctx).First(&comment, "id =?", req.GetParams().GetId()).Error; err != nil {
		return nil, status.Error(http.StatusInternalServerError, err.Error())
	}

	var ReservationFeedback *feedbackDomain.ReservationFeedback
	if err := r.GetTenantDBConnection(ctx).First(&ReservationFeedback, "id = ?", comment.ReservationFeedbackID).Error; err != nil {
		return nil, status.Error(http.StatusInternalServerError, err.Error())
	}

	currentUser, err := userRepo.GetLoggedInUser(ctx)
	if err != nil {
		return nil, status.Error(http.StatusInternalServerError, err.Error())
	}

	var role *roleDomain.Role
	if err := r.GetTenantDBConnection(ctx).First(&role, "id = ?", currentUser.RoleID).Error; err != nil {
		return nil, status.Error(http.StatusInternalServerError, err.Error())
	}
	currentUser.Role = role
	comment.Creator = currentUser

	return &guestProto.UpdateReservationFeedbackCommentResponse{
		Result: converters.BuildReservationFeedbackCommentResponse(comment),
	}, nil
}
