package repository

import (
	"context"
	"errors"
	"net/http"

	"github.com/goplaceapp/goplace-common/pkg/errorhelper"
	guestProto "github.com/goplaceapp/goplace-guest/api/v1"
	"github.com/goplaceapp/goplace-guest/internal/services/reservation/adapters/converters"
	"github.com/goplaceapp/goplace-guest/pkg/reservationservice/domain"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"
)

func (r *ReservationRepository) DeleteReservationNote(ctx context.Context, req *guestProto.DeleteReservationNoteRequest) (*guestProto.DeleteReservationNoteResponse, error) {
	if err := r.GetTenantDBConnection(ctx).
		Delete(&domain.ReservationNote{}, "id = ? AND reservation_id = ?",
			req.GetId(),
			req.GetReservationId()).
		Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, status.Error(http.StatusNotFound, errorhelper.ErrNoteNotFound)
		}

		return nil, status.Error(http.StatusNotFound, err.Error())
	}

	reservation, err := r.CommonRepo.GetReservationByID(ctx, req.GetReservationId())
	if err != nil {
		return nil, status.Error(http.StatusInternalServerError, err.Error())
	}

	return &guestProto.DeleteReservationNoteResponse{
		Code:        http.StatusOK,
		Message:     "Reservation note deleted successfully",
		Reservation: converters.BuildReservationResponse(reservation),
	}, nil
}
