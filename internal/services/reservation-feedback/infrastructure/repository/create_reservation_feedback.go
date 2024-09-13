package repository

import (
	"context"
	"errors"
	"net/http"

	reservationDomain "github.com/goplaceapp/goplace-guest/pkg/reservationservice/domain"

	"github.com/goplaceapp/goplace-common/pkg/errorhelper"
	"github.com/goplaceapp/goplace-common/pkg/meta"
	"github.com/goplaceapp/goplace-common/pkg/utils"
	guestProto "github.com/goplaceapp/goplace-guest/api/v1"
	"github.com/goplaceapp/goplace-guest/internal/services/reservation-feedback/domain"
	"github.com/goplaceapp/goplace-guest/pkg/guestservice/converters"
	guestDomain "github.com/goplaceapp/goplace-guest/pkg/guestservice/domain"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"
)

var FeedbackStatuses = []string{
	meta.Pending,
	meta.Solved,
}

func (r *ReservationFeedbackRepository) CreateReservationFeedback(ctx context.Context, req *guestProto.CreateReservationFeedbackRequest) (*guestProto.CreateReservationFeedbackResponse, error) {
	var (
		err         error
		reservation *reservationDomain.Reservation
	)

	feedback := &domain.ReservationFeedback{
		ReservationID: req.GetParams().GetReservationId(),
		Rate:          req.GetParams().GetRate(),
		Description:   req.GetParams().GetDescription(),
	}

	if feedback.Rate <= 3 {
		for i, v := range FeedbackStatuses {
			if utils.CompareStr(v, meta.Pending) {
				feedback.StatusID = int32(i + 1)
				feedback.Status = v
				break
			}
		}
	}

	if err := r.GetTenantDBConnection(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&feedback).Error; err != nil {
			if errors.Is(err, gorm.ErrDuplicatedKey) {
				return status.Error(http.StatusConflict, "Duplicated feedback for this reservation")
			} else if errors.Is(err, gorm.ErrForeignKeyViolated) {
				return status.Error(http.StatusInternalServerError, errorhelper.ErrReservationNotFound)
			}
			return status.Error(http.StatusInternalServerError, err.Error())
		}

		sectionIds := req.GetParams().GetSectionIds()
		for _, sectionId := range sectionIds {
			assignment := domain.ReservationFeedbackSectionAssignment{
				FeedbackID: feedback.ID,
				SectionID:  sectionId,
			}

			if err := tx.Create(&assignment).Error; err != nil {
				return status.Error(http.StatusInternalServerError, err.Error())
			}
		}
		return nil
	}); err != nil {
		return nil, err
	}

	feedback.Sections, err = r.GetAllReservationFeedbackSectionsByID(ctx, feedback.ID)
	if err != nil {
		return nil, status.Error(http.StatusInternalServerError, err.Error())
	}

	reservation, err = r.CommonRepo.GetReservationByID(ctx, feedback.ReservationID)
	if err != nil {
		return nil, status.Error(http.StatusInternalServerError, errorhelper.ErrReservationNotFound)
	}
	feedback.Reservation = reservation

	getGuest, err := r.CommonRepo.GetAllGuestData(ctx, &guestDomain.Guest{ID: reservation.GuestID})
	if err != nil {
		return nil, status.Error(http.StatusInternalServerError, errorhelper.ErrGuestNotFound)
	}
	feedback.Guest = getGuest

	return &guestProto.CreateReservationFeedbackResponse{
		Result: converters.BuildReservationFeedbackResponse(feedback),
	}, nil
}
