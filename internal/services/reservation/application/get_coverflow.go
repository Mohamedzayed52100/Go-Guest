package application

import (
	"context"
	guestProto "github.com/goplaceapp/goplace-guest/api/v1"
)

func (s *ReservationService) GetReservationsCoverFlow(ctx context.Context, req *guestProto.GetReservationsCoverFlowRequest) (*guestProto.GetReservationsCoverFlowResponse, error) {
	res, err := s.Repository.GetReservationsCoverFlow(ctx, req)
	if err != nil {
		return nil, err
	}
	return res, nil
}
