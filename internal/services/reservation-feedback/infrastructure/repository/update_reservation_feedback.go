package repository

import (
	"context"
	"errors"
	"net/http"

	"github.com/goplaceapp/goplace-common/pkg/errorhelper"
	"github.com/goplaceapp/goplace-common/pkg/meta"
	guestProto "github.com/goplaceapp/goplace-guest/api/v1"
	"github.com/goplaceapp/goplace-guest/internal/services/reservation-feedback/domain"
	"github.com/goplaceapp/goplace-guest/pkg/guestservice/converters"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"
)

func (r *ReservationFeedbackRepository) UpdateReservationFeedback(ctx context.Context, req *guestProto.UpdateReservationFeedbackRequest) (*guestProto.UpdateReservationFeedbackResponse, error) {
	var (
		feedback = &domain.ReservationFeedback{}
		updates  = make(map[string]interface{})
		params   = req.GetParams()
		err      error
	)

	if params.GetReservationId() != 0 {
		updates["reservation_id"] = params.GetReservationId()
	}

	if params.GetDescription() != "" {
		updates["description"] = params.GetDescription()
	}

	if params.GetRate() != 0 {
		updates["rate"] = params.GetRate()
	}

	if params.GetStatus() > 0 {
		updates["status_id"] = params.GetStatus()
		if params.GetStatus() > int32(len(FeedbackStatuses)) {
			return nil, status.Error(http.StatusBadRequest, "Invalid status id")
		}
		feedback.StatusID = params.GetStatus()
		feedback.Status = FeedbackStatuses[params.GetStatus()-1]
		if feedback.Status != meta.Solved {
			updates["solution_id"] = 0
		}
		updates["status_id"] = feedback.StatusID
	}

	if err := r.GetTenantDBConnection(ctx).
		Model(feedback).
		Joins("JOIN reservations ON reservations.id = reservation_feedbacks.reservation_id").
		Select("reservation_feedbacks.*").
		First(&feedback, "reservation_feedbacks.id = ? AND reservation_id = ? AND branch_id = ?",
			req.GetParams().GetId(),
			req.GetParams().GetReservationId(),
			r.userClient.Client.UserService.Repository.GetCurrentBranchId(ctx)).Error; err != nil {
		return nil, status.Error(http.StatusInternalServerError, "Reservation not found")
	}

	if err := r.GetTenantDBConnection(ctx).
		Model(&domain.ReservationFeedback{}).
		Where("id = ? AND reservation_id = ?",
			params.GetId(),
			params.GetReservationId()).
		Updates(updates).
		First(&feedback).
		Error; err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return nil, status.Error(http.StatusConflict, "Duplicated feedback for this reservation")
		} else if errors.Is(err, gorm.ErrForeignKeyViolated) {
			return nil, status.Error(http.StatusInternalServerError, errorhelper.ErrReservationNotFound)
		}

		return nil, status.Error(http.StatusInternalServerError, err.Error())
	}

	if !req.GetParams().GetEmptySections() {
		if err := r.GetTenantDBConnection(ctx).Transaction(func(tx *gorm.DB) error {
			oldSectionIds := []int32{}

			if err := tx.
				Model(&domain.ReservationFeedbackSectionAssignment{}).
				Select("section_id").
				Where("feedback_id = ?", params.GetId()).
				Find(&oldSectionIds).
				Error; err != nil {
				return err
			}

			for _, sectionId := range oldSectionIds {
				if err := tx.Model(&domain.ReservationFeedbackSectionAssignment{}).
					Where("feedback_id =? AND section_id =?", params.GetId(), sectionId).
					Delete(&domain.ReservationFeedbackSectionAssignment{}).
					Error; err != nil {
					return err
				}
			}

			for _, sectionId := range req.GetParams().GetSectionIds() {
				if err := tx.Create(&domain.ReservationFeedbackSectionAssignment{
					FeedbackID: params.GetId(),
					SectionID:  sectionId,
				}).Error; err != nil {
					return err
				}
			}

			return nil
		}); err != nil {
			return nil, status.Error(http.StatusInternalServerError, err.Error())
		}
	}

	feedback, err = r.GetReservationFeedbackData(ctx, feedback)
	if err != nil {
		return nil, status.Error(http.StatusInternalServerError, err.Error())
	}

	return &guestProto.UpdateReservationFeedbackResponse{
		Result: converters.BuildReservationFeedbackResponse(feedback),
	}, nil
}
