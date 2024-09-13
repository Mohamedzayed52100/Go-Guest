package grpc

import (
	"context"
	guestProto "github.com/goplaceapp/goplace-guest/api/v1"
)

func (s *GuestServiceServer) AddGuestNote(ctx context.Context, req *guestProto.AddGuestNoteRequest) (*guestProto.AddGuestNoteResponse, error) {
	res, err := s.guestService.AddGuestNote(ctx, req)
	if err != nil {
		return nil, err
	}

	return res, nil
}
