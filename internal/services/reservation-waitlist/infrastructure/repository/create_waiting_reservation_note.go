package repository

import (
	"context"
	"net/http"
	"time"

	"github.com/goplaceapp/goplace-common/pkg/errorhelper"
	guestProto "github.com/goplaceapp/goplace-guest/api/v1"
	"github.com/goplaceapp/goplace-guest/internal/services/reservation-waitlist/adapters/converters"
	reservationWaitlistDomain "github.com/goplaceapp/goplace-guest/internal/services/reservation-waitlist/domain"
	"google.golang.org/grpc/status"
)

func (r *ReservationWaitListRepository) CreateWaitingReservationNote(ctx context.Context, req *guestProto.CreateWaitingReservationNoteRequest) (*guestProto.CreateWaitingReservationNoteResponse, error) {
	loggedInUser, err := r.userClient.Client.UserService.Repository.GetLoggedInUser(ctx)
	if err != nil {
		return nil, status.Errorf(http.StatusNotFound, err.Error())
	}

	if err := r.GetTenantDBConnection(ctx).
		Model(&reservationWaitlistDomain.ReservationWaitlist{}).
		Where("id = ?", req.GetReservationWaitlistId()).
		First(&reservationWaitlistDomain.ReservationWaitlist{}).Error; err != nil {
		return nil, status.Errorf(http.StatusNotFound, errorhelper.ErrReservationNotFound)
	}

	note := &reservationWaitlistDomain.ReservationWaitlistNote{
		CreatorID:   loggedInUser.ID,
		Description: req.GetDescription(),
	}

	if err := r.GetTenantDBConnection(ctx).Create(note).Error; err != nil {
		return nil, status.Errorf(http.StatusInternalServerError, err.Error())
	}

	if err := r.GetTenantDBConnection(ctx).
		Model(&reservationWaitlistDomain.ReservationWaitlist{}).
		Where("id = ?", req.GetReservationWaitlistId()).
		Updates(&reservationWaitlistDomain.ReservationWaitlist{
			NoteID:    &note.ID,
			UpdatedAt: time.Now(),
		}).Error; err != nil {
		return nil, status.Errorf(http.StatusInternalServerError, err.Error())
	}

	note.Creator = loggedInUser

	if _, err := r.CreateReservationWaitListLogs(ctx, &reservationWaitlistDomain.ReservationWaitlistLog{
		ReservationWaitlistID: req.GetReservationWaitlistId(),
		Action:                "create",
		FieldName:             "note",
		NewValue:              note.Description,
	}); err != nil {
		return nil, status.Errorf(http.StatusNotFound, err.Error())
	}

	return &guestProto.CreateWaitingReservationNoteResponse{
		Result: converters.BuildWaitingReservationNoteResponse(note),
	}, nil
}
