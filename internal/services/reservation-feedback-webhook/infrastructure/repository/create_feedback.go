package repository

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/goplaceapp/goplace-common/pkg/meta"
	"github.com/goplaceapp/goplace-common/pkg/utils"
	guestProto "github.com/goplaceapp/goplace-guest/api/v1"
	"github.com/goplaceapp/goplace-guest/internal/services/reservation-feedback/domain"
	"github.com/goplaceapp/goplace-guest/internal/services/reservation-feedback/infrastructure/repository"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"
)

func (r *ReservationFeedbackWebhookRepository) CreateReservationFeedbackFromWebhook(ctx context.Context, req *guestProto.CreateReservationFeedbackFromWebhookRequest) (*guestProto.CreateReservationFeedbackFromWebhookResponse, error) {
	feedback := &domain.ReservationFeedback{
		ReservationID: req.GetReservationId(),
		Rate:          req.GetRate(),
		Description:   req.GetFeedback(),
	}

	if feedback.Rate <= 3 {
		for i, v := range repository.FeedbackStatuses {
			if utils.CompareStr(v, meta.Pending) {
				feedback.StatusID = int32(i + 1)
				feedback.Status = v
				break
			}
		}
	}
	if err := r.GetTenantDBConnection(ctx).
		WithContext(ctx).
		Model(&domain.ReservationFeedback{}).
		Create(&feedback).Error; err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			if req.GetFeedback() != "" {
				if err := r.GetTenantDBConnection(ctx).
					WithContext(ctx).
					Model(&domain.ReservationFeedback{}).
					Where("reservation_id = ?", req.GetReservationId()).
					Updates(&domain.ReservationFeedback{
						Description: req.GetFeedback(),
					}).Error; err != nil {
					return nil, status.Error(http.StatusInternalServerError, fmt.Sprintf("Failed to update feedback: %v", err))
				}
			}

			return &guestProto.CreateReservationFeedbackFromWebhookResponse{
				Code:    200,
				Message: "Reservation feedback updated successfully",
			}, nil
		} else if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, status.Error(http.StatusNotFound, fmt.Sprintf("Reservation %d not found", req.GetReservationId()))
		}

		return nil, status.Error(http.StatusInternalServerError, fmt.Sprintf("Failed to create feedback: %v", err))
	}

	return &guestProto.CreateReservationFeedbackFromWebhookResponse{
		Code:    200,
		Message: "Reservation feedback created successfully",
	}, nil
}
