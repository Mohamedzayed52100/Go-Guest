package repository

import (
	"context"
	"fmt"
	"github.com/goplaceapp/goplace-guest/pkg/reservationservice/domain"
	"net/http"
	"time"

	"github.com/goplaceapp/goplace-common/pkg/meta"
	guestProto "github.com/goplaceapp/goplace-guest/api/v1"
	logDomain "github.com/goplaceapp/goplace-guest/internal/services/reservation-log/domain"
	"google.golang.org/grpc/status"
)

func (r *ReservationRepository) UpdateReservationFromWebhook(ctx context.Context, req *guestProto.UpdateReservationFromWebhookRequest) (*guestProto.UpdateReservationFromWebhookResponse, error) {
	var (
		statusId int32
		log      *logDomain.ReservationLog
	)

	currentReservation, err := r.CommonRepo.GetReservationByID(ctx, req.GetReservationId())
	if err != nil {
		return nil, status.Error(http.StatusInternalServerError, fmt.Sprintf("Failed to get reservation: %v", err))
	}

	switch req.GetStatus() {
	case "Confirmed":
		getStatus, err := r.GetReservationStatusByName(ctx, meta.Confirmed, currentReservation.BranchID)
		if err != nil {
			return nil, status.Error(http.StatusInternalServerError, fmt.Sprintf("Failed to get reservation status: %v", err))
		}
		statusId = int32(getStatus.ID)

		if currentReservation.Status.Name != meta.Confirmed {
			log = &logDomain.ReservationLog{
				MadeBy:        "Guest",
				ReservationID: req.GetReservationId(),
				FieldName:     "status",
				OldValue:      currentReservation.Status.Name,
				NewValue:      getStatus.Name,
				Action:        "update",
			}

			if err := r.GetTenantDBConnection(ctx).
				Model(&logDomain.ReservationLog{}).
				Create(log).Error; err != nil {
				return nil, err
			}
		}
	case "Canceled":
		getStatus, err := r.GetReservationStatusByName(ctx, meta.Cancelled, currentReservation.BranchID)
		if err != nil {
			return nil, status.Error(http.StatusInternalServerError, fmt.Sprintf("Failed to get reservation status: %v", err))
		}
		statusId = int32(getStatus.ID)

		if currentReservation.Status.Name != meta.Cancelled {
			log = &logDomain.ReservationLog{
				MadeBy:        "Guest",
				ReservationID: req.GetReservationId(),
				FieldName:     "status",
				OldValue:      currentReservation.Status.Name,
				NewValue:      getStatus.Name,
				Action:        "update",
			}

			if err := r.GetTenantDBConnection(ctx).
				Model(&logDomain.ReservationLog{}).
				Create(log).Error; err != nil {
				return nil, err
			}
		}
	}

	if err := r.GetTenantDBConnection(ctx).
		WithContext(ctx).
		Model(&domain.Reservation{}).
		Where("id = ?", req.GetReservationId()).
		Updates(&domain.Reservation{
			StatusID:  statusId,
			UpdatedAt: time.Now(),
		}).Error; err != nil {
		return nil, status.Error(http.StatusInternalServerError, fmt.Sprintf("Failed to update reservation: %v", err))
	}

	return &guestProto.UpdateReservationFromWebhookResponse{
		Code:    200,
		Message: "Reservation updated successfully",
	}, nil
}
