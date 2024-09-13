package repository

import (
	"context"
	"net/http"

	"github.com/goplaceapp/goplace-common/pkg/meta"
	guestProto "github.com/goplaceapp/goplace-guest/api/v1"
	"github.com/goplaceapp/goplace-guest/internal/services/reservation-feedback-comment/adapters/converters"
	feedbackDomain "github.com/goplaceapp/goplace-guest/internal/services/reservation-feedback/domain"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"
)

func (r *ReservationFeedbackCommentRepository) AddReservationFeedbackSolution(ctx context.Context, req *guestProto.AddReservationFeedbackSolutionRequest) (*guestProto.AddReservationFeedbackSolutionResponse, error) {
	userRepo := r.userClient.Client.UserService.Repository

	currentUser, err := userRepo.GetLoggedInUser(ctx)
	if err != nil {
		return nil, status.Error(http.StatusInternalServerError, err.Error())
	}

	solution := &feedbackDomain.ReservationFeedbackSolution{
		Solution:  req.GetParams().GetSolution(),
		CreatorID: currentUser.ID,
	}

	var reservationFeedback *feedbackDomain.ReservationFeedback
	if err := r.GetTenantDBConnection(ctx).First(&reservationFeedback, "id = ?", req.GetParams().GetFeedbackId()).Error; err != nil {
		return nil, status.Error(http.StatusInternalServerError, err.Error())
	}

	if reservationFeedback.SolutionID != 0 {
		return nil, status.Error(http.StatusBadRequest, "Feedback already has a solution")
	}

	if err := r.GetTenantDBConnection(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(solution).Error; err != nil {
			return status.Error(http.StatusInternalServerError, err.Error())
		}

		if reservationFeedback.StatusID != 0 && feedbackStatuses[reservationFeedback.StatusID-1] != meta.Solved {
			var solvedStatusId int
			for i, v := range feedbackStatuses {
				if v == meta.Solved {
					solvedStatusId = i + 1
					break
				}
			}

			if err := tx.Model(&feedbackDomain.ReservationFeedback{}).
				Where("id = ?", req.GetParams().GetFeedbackId()).
				Updates(map[string]interface{}{
					"solution_id": solution.ID,
					"status_id":   solvedStatusId,
				}).Error; err != nil {
				return status.Error(http.StatusInternalServerError, err.Error())
			}
		} else {
			if err := tx.Model(&feedbackDomain.ReservationFeedback{}).
				Where("id = ?", req.GetParams().GetFeedbackId()).
				Update("solution_id", solution.ID).Error; err != nil {
				return status.Error(http.StatusInternalServerError, err.Error())
			}
		}

		return nil
	}); err != nil {
		return nil, err
	}

	solution.Creator = currentUser

	return &guestProto.AddReservationFeedbackSolutionResponse{
		Result: converters.BuildReservationFeedbackSolutionResponse(solution),
	}, nil
}
