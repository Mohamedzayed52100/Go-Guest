package repository

import (
	"context"
	"net/http"

	guestProto "github.com/goplaceapp/goplace-guest/api/v1"
	reservationLogDomain "github.com/goplaceapp/goplace-guest/internal/services/reservation-log/domain"
	"github.com/goplaceapp/goplace-guest/internal/services/reservation/adapters/converters"
	extDomain "github.com/goplaceapp/goplace-guest/pkg/reservationservice/domain"
	"google.golang.org/grpc/status"
)

func (r *ReservationRepository) AddReservationNote(ctx context.Context, req *guestProto.AddReservationNoteRequest) (*guestProto.AddReservationNoteResponse, error) {
	var (
		note      *extDomain.ReservationNote
		creatorID int32
	)

	loggedInUser, err := r.userClient.Client.UserService.Repository.GetLoggedInUser(ctx)
	if err != nil {
		return nil, status.Error(http.StatusInternalServerError, err.Error())
	}

	if loggedInUser == nil {
		creatorID = 0
	} else {
		creatorID = loggedInUser.ID
	}

	if err := r.GetTenantDBConnection(ctx).
		Model(&extDomain.ReservationNote{}).
		Create(&extDomain.ReservationNote{
			CreatorID:     creatorID,
			Description:   req.GetParams().GetDescription(),
			ReservationID: req.GetParams().GetReservationId(),
		}).Scan(&note).Error; err != nil {
		return nil, status.Error(http.StatusInternalServerError, err.Error())
	}

	if creatorID != 0 {
		getCreator, err := r.userClient.Client.UserService.Repository.GetUserProfileByID(ctx, note.CreatorID)
		if err != nil {
			return nil, status.Error(http.StatusInternalServerError, err.Error())
		}

		note.Creator = getCreator
	}

	if _, err := r.CreateReservationLogs(ctx, &reservationLogDomain.ReservationLog{
		ReservationID: note.ReservationID,
		Action:        "create",
		FieldName:     "note",
		NewValue:      note.Description,
	}); err != nil {
		return nil, status.Errorf(http.StatusNotFound, err.Error())
	}

	res, err := r.CommonRepo.GetReservationByID(ctx, note.ReservationID)
	if err != nil {
		return nil, status.Error(http.StatusInternalServerError, err.Error())
	}
	note.Reservation = res

	return &guestProto.AddReservationNoteResponse{
		Result: converters.BuildReservationNoteProto(note),
	}, nil
}
