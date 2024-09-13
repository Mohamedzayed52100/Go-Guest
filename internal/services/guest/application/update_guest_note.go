package application

import (
	"context"
	guestProto "github.com/goplaceapp/goplace-guest/api/v1"
)

func (s *GuestService) UpdateGuestNote(ctx context.Context, req *guestProto.UpdateGuestNoteRequest) (*guestProto.UpdateGuestNoteResponse, error) {
	res, err := s.Repository.UpdateGuestNote(ctx, req)
	if err != nil {
		return nil, err
	}

	return res, nil
}
