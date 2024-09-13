package repository

import (
	"context"
	"net/http"

	guestProto "github.com/goplaceapp/goplace-guest/api/v1"
	"github.com/goplaceapp/goplace-guest/internal/services/reservation-waitlist/adapters/converters"
	reservationWailistDomain "github.com/goplaceapp/goplace-guest/internal/services/reservation-waitlist/domain"
	"google.golang.org/grpc/status"
)

func (r *ReservationWaitListRepository) UpdateWaitingReservationNote(ctx context.Context, req *guestProto.UpdateWaitingReservationNoteRequest) (*guestProto.UpdateWaitingReservationNoteResponse, error) {
	var oldNote *reservationWailistDomain.ReservationWaitlistNote

	if err := r.GetTenantDBConnection(ctx).First(&oldNote, "id = ?", req.GetId()).Error; err != nil {
		return nil, status.Error(http.StatusNotFound, err.Error())
	}

	currentUser, err := r.userClient.Client.UserService.Repository.GetLoggedInUser(ctx)
	if err != nil {
		return nil, status.Error(http.StatusInternalServerError, err.Error())
	}

	if req.GetDescription() != "" && oldNote.Description != req.GetDescription() {
		var (
			note       *reservationWailistDomain.ReservationWaitlistNote
			waitListId int32
		)

		if err := r.GetTenantDBConnection(ctx).
			Model(&reservationWailistDomain.ReservationWaitlistNote{}).
			Where("id = ?", req.GetId()).
			Update("description", req.GetDescription()).
			Error; err != nil {
			return nil, status.Error(http.StatusInternalServerError, err.Error())
		}

		if err := r.GetTenantDBConnection(ctx).First(&note, "id = ?", req.GetId()).Error; err != nil {
			return nil, status.Error(http.StatusInternalServerError, err.Error())
		}

		note.Creator = currentUser

		if err := r.GetTenantDBConnection(ctx).
			Model(&reservationWailistDomain.ReservationWaitlist{}).
			Where("note_id = ?", req.GetId()).
			Select("id").
			Scan(&waitListId).
			Error; err != nil {
			return nil, status.Error(http.StatusInternalServerError, err.Error())
		}

		if _, err := r.CreateReservationWaitListLogs(ctx, &reservationWailistDomain.ReservationWaitlistLog{
			ReservationWaitlistID: waitListId,
			Action:                "update",
			FieldName:             "note",
			OldValue:              oldNote.Description,
			NewValue:              note.Description,
		}); err != nil {
			return nil, status.Errorf(http.StatusNotFound, err.Error())
		}

		return &guestProto.UpdateWaitingReservationNoteResponse{
			Result: converters.BuildWaitingReservationNoteResponse(note),
		}, nil
	}

	oldNote.Creator = currentUser

	return &guestProto.UpdateWaitingReservationNoteResponse{
		Result: converters.BuildWaitingReservationNoteResponse(oldNote),
	}, nil
}
