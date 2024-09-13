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

func (r *GuestRepository) AddGuestNote(ctx context.Context, req *guestProto.AddGuestNoteRequest) (*guestProto.AddGuestNoteResponse, error) {
	var note *domain.GuestNote

	loggedInUser, err := r.UserClient.Client.UserService.Repository.GetLoggedInUser(ctx)
	if err != nil {
		return nil, status.Errorf(http.StatusNotFound, err.Error())
	}

	if err := r.GetTenantDBConnection(ctx).
		Model(&domain.Guest{}).
		Where("id = ?", req.GetParams().GetGuestId()).
		First(&domain.Guest{}).Error; err != nil {
		return nil, status.Errorf(http.StatusNotFound, errorhelper.ErrGuestNotFound)
	}

	if err := r.GetTenantDBConnection(ctx).
		Model(&domain.GuestNote{}).
		Create(&domain.GuestNote{
			GuestID:     req.GetParams().GetGuestId(),
			CreatorID:   loggedInUser.ID,
			Description: req.GetParams().GetDescription(),
		}).
		Scan(&note).
		Error; err != nil {
		return nil, status.Errorf(http.StatusNotFound, err.Error())
	}

	note.Creator = loggedInUser

	if _, err := r.CreateGuestLogs(ctx, &logDomain.GuestLog{
		GuestID:   note.GuestID,
		Action:    "create",
		FieldName: "note",
		NewValue:  note.Description,
	}); err != nil {
		return nil, status.Errorf(http.StatusNotFound, err.Error())
	}

	return &guestProto.AddGuestNoteResponse{
		Result: converters.BuildGuestNoteResponse(note),
	}, nil
}
