package repository

import (
	"context"

	guestProto "github.com/goplaceapp/goplace-guest/api/v1"
	"github.com/goplaceapp/goplace-guest/pkg/guestservice/converters"
	"github.com/goplaceapp/goplace-guest/internal/services/reservation-feedback/domain"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (r *ReservationFeedbackRepository) GetAllReservationFeedbackSections(ctx context.Context, req *emptypb.Empty) (*guestProto.GetAllReservationsFeedbackSectionsResponse, error) {
	var sections []*domain.ReservationFeedbackSection

	loggedInUser, err := r.userClient.Client.UserService.Repository.GetLoggedInUser(ctx)
	if err != nil {
		return nil, err
	}

	if err := r.GetTenantDBConnection(ctx).
		Where("branch_id = ?", loggedInUser.BranchID).
		Find(&sections).Error; err != nil {
		return nil, err
	}

	return &guestProto.GetAllReservationsFeedbackSectionsResponse{
		Result: converters.BuildAllFeedbackSectionsResponse(sections),
	}, nil
}

func (r *ReservationFeedbackRepository) GetAllReservationFeedbackSectionsByID(ctx context.Context, feedbackId int32) ([]*domain.ReservationFeedbackSection, error) {
	sections := []*domain.ReservationFeedbackSection{}

	if err := r.GetTenantDBConnection(ctx).
		Model(sections).
		Joins("JOIN reservation_feedback_section_assignments ON "+
			"reservation_feedback_section_assignments.section_id = reservation_feedback_sections.id").
		Find(&sections, "reservation_feedback_section_assignments.feedback_id = ?", feedbackId).
		Error; err != nil {
		return nil, err
	}

	return sections, nil
}
