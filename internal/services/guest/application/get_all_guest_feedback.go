package application

import (
	"context"
	guestProto "github.com/goplaceapp/goplace-guest/api/v1"
)

func (s *GuestService) GetAllGuestFeedback(ctx context.Context, req *guestProto.GetAllGuestFeedbackRequest) (*guestProto.GetAllGuestFeedbackResponse, error) {
	res, err := s.Repository.GetAllGuestFeedback(ctx, req)
	if err != nil {
		return nil, err
	}

	return res, nil
}
