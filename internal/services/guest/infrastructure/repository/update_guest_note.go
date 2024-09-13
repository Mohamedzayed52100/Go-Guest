package repository

import (
	"context"
	"net/http"

	"github.com/goplaceapp/goplace-common/pkg/errorhelper"
	guestProto "github.com/goplaceapp/goplace-guest/api/v1"
	logDomain "github.com/goplaceapp/goplace-guest/internal/services/guest-log/domain"
	"github.com/goplaceapp/goplace-guest/pkg/guestservice/converters"
	"github.com/goplaceapp/goplace-guest/pkg/guestservice/domain"
	"google.golang.org/grpc/status"
)

func (r *GuestRepository) UpdateGuestNote(ctx context.Context, req *guestProto.UpdateGuestNoteRequest) (*guestProto.UpdateGuestNoteResponse, error) {
	oldNote := &domain.GuestNote{}
	updates := make(map[string]interface{})

	getCreator, err := r.UserClient.Client.UserService.Repository.GetLoggedInUser(ctx)
	if err != nil {
		return nil, err
	}

	if err := r.GetTenantDBConnection(ctx).
		First(&oldNote, "id = ?", req.GetParams().GetId()).
		Error; err != nil {
		return nil, status.Error(http.StatusInternalServerError, errorhelper.ErrNoteNotFound)
	}

	if req.GetParams().GetDescription() != "" &&
		oldNote.Description != req.GetParams().GetDescription() {
		updates["description"] = req.GetParams().GetDescription()

		if err := r.GetTenantDBConnection(ctx).
			Model(&domain.GuestNote{}).
			Where("id = ?", req.GetParams().GetId()).
			Updates(updates).
			Error; err != nil {
			return nil, status.Error(http.StatusInternalServerError, err.Error())
		}

		var note *domain.GuestNote
		if err := r.GetTenantDBConnection(ctx).
			First(&note, "id = ?", req.GetParams().GetId()).
			Error; err != nil {
			return nil, status.Error(http.StatusInternalServerError, err.Error())
		}

		note.Creator = getCreator

		if _, err := r.CreateGuestLogs(ctx, &logDomain.GuestLog{
			GuestID:   note.GuestID,
			Action:    "update",
			FieldName: "note",
			OldValue:  oldNote.Description,
			NewValue:  note.Description,
		}); err != nil {
			return nil, status.Errorf(http.StatusNotFound, err.Error())
		}

		return &guestProto.UpdateGuestNoteResponse{
			Result: converters.BuildGuestNoteResponse(note),
		}, nil
	}

	oldNote.Creator = getCreator

	return &guestProto.UpdateGuestNoteResponse{
		Result: converters.BuildGuestNoteResponse(oldNote),
	}, nil
}
