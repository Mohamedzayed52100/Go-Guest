package repository

import (
	"context"
	"net/http"

	"github.com/goplaceapp/goplace-common/pkg/errorhelper"
	guestProto "github.com/goplaceapp/goplace-guest/api/v1"
	logDomain "github.com/goplaceapp/goplace-guest/internal/services/guest-log/domain"
	"github.com/goplaceapp/goplace-guest/pkg/guestservice/domain"
	"google.golang.org/grpc/status"
)

/*
DeleteGuestNote deletes a specific note of a guest from the database.

Parameters:
- ctx: The context for timeout and cancellation signals.
- req: The request containing the ID of the note and the ID of the guest.

The method:
1. Checks if the note exists in the database.
2. Deletes the note from the database.

Returns:
- A response with a success message if the note is successfully deleted.
- An error if the note does not exist or if there is an error during deletion.
*/

func (r *GuestRepository) DeleteGuestNote(ctx context.Context, req *guestProto.DeleteGuestNoteRequest) (*guestProto.DeleteGuestNoteResponse, error) {
	var oldNote *domain.GuestNote
	if err := r.GetTenantDBConnection(ctx).
		First(&oldNote, "id = ? AND guest_id = ?", req.GetId(), req.GetGuestId()).
		Error; err != nil {
		return nil, status.Error(http.StatusNotFound, errorhelper.ErrNoteNotFound)
	}

	if err := r.GetTenantDBConnection(ctx).
		Delete(&domain.GuestNote{}, "id = ? AND guest_id = ?", req.GetId(), req.GetGuestId()).
		Error; err != nil {
		return nil, status.Error(http.StatusNotFound, err.Error())
	}

	if _, err := r.CreateGuestLogs(ctx, &logDomain.GuestLog{
		GuestID:   req.GetGuestId(),
		Action:    "delete",
		FieldName: "note",
		OldValue:  oldNote.Description,
	}); err != nil {
		return nil, status.Errorf(http.StatusNotFound, err.Error())
	}

	return &guestProto.DeleteGuestNoteResponse{
		Code:    http.StatusOK,
		Message: "Guest note deleted successfully",
	}, nil
}
