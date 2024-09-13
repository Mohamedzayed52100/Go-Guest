package grpc

import (
	"context"
	guestProto "github.com/goplaceapp/goplace-guest/api/v1"
)

func (s *GuestServiceServer) UpdateGuestNote(ctx context.Context, req *guestProto.UpdateGuestNoteRequest) (*guestProto.UpdateGuestNoteResponse, error) {
	res, err := s.guestService.UpdateGuestNote(ctx, req)
	if err != nil {
		return nil, err
	}

	return res, nil
}
