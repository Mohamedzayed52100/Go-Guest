package repository

import (
	"context"
	"net/http"

	"github.com/goplaceapp/goplace-common/pkg/errorhelper"
	guestProto "github.com/goplaceapp/goplace-guest/api/v1"
	reservationLogDomain "github.com/goplaceapp/goplace-guest/internal/services/reservation-log/domain"
	"github.com/goplaceapp/goplace-guest/internal/services/reservation/adapters/converters"
	"github.com/goplaceapp/goplace-guest/pkg/reservationservice/domain"
	"google.golang.org/grpc/status"
)

func (r *ReservationRepository) UpdateReservationNote(ctx context.Context, req *guestProto.UpdateReservationNoteRequest) (*guestProto.UpdateReservationNoteResponse, error) {
	var (
		oldNote = &domain.ReservationNote{}
		updates = make(map[string]interface{})
	)

	getCreator, err := r.userClient.Client.UserService.Repository.GetLoggedInUser(ctx)
	if err != nil {
		return nil, err
	}

	if err := r.GetTenantDBConnection(ctx).
		First(&oldNote, "id = ?", req.GetParams().GetId()).
		Error; err != nil {
		return nil, status.Error(http.StatusInternalServerError, errorhelper.ErrNoteNotFound)
	}

	if req.GetParams().GetDescription() != "" && oldNote.Description != req.GetParams().GetDescription() {
		updates["description"] = req.GetParams().GetDescription()

		if err := r.GetTenantDBConnection(ctx).
			Model(&domain.ReservationNote{}).
			Where("id = ?", req.GetParams().GetId()).
			Updates(updates).
			Error; err != nil {
			return nil, status.Error(http.StatusInternalServerError, err.Error())
		}

		var note *domain.ReservationNote
		if err := r.GetTenantDBConnection(ctx).
			First(&note, "id = ?", req.GetParams().GetId()).
			Error; err != nil {
			return nil, status.Error(http.StatusInternalServerError, err.Error())
		}

		note.Creator = getCreator

		if _, err := r.CreateReservationLogs(ctx, &reservationLogDomain.ReservationLog{
			ReservationID: note.ReservationID,
			Action:        "update",
			FieldName:     "note",
			OldValue:      oldNote.Description,
			NewValue:      note.Description,
		}); err != nil {
			return nil, status.Errorf(http.StatusNotFound, err.Error())
		}

		return &guestProto.UpdateReservationNoteResponse{
			Result: converters.BuildReservationNoteProto(note),
		}, nil
	}

	oldNote.Creator = getCreator

	res, err := r.CommonRepo.GetReservationByID(ctx, oldNote.ReservationID)
	if err != nil {
		return nil, status.Error(http.StatusInternalServerError, err.Error())
	}
	oldNote.Reservation = res

	return &guestProto.UpdateReservationNoteResponse{
		Result: converters.BuildReservationNoteProto(oldNote),
	}, nil
}
