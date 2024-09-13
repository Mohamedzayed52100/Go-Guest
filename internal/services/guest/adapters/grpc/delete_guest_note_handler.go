package grpc

import (
	"context"
	guestProto "github.com/goplaceapp/goplace-guest/api/v1"
)

func (s *GuestServiceServer) DeleteGuestNote(ctx context.Context, req *guestProto.DeleteGuestNoteRequest) (*guestProto.DeleteGuestNoteResponse, error) {
	res, err := s.guestService.DeleteGuestNote(ctx, req)
	if err != nil {
		return nil, err
	}

	return res, nil
}
