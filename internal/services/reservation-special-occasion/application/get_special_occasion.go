package application

import (
	"context"

	guestProto "github.com/goplaceapp/goplace-guest/api/v1"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (s *ReservationSpecialOccasionService) GetAllSpecialOccasions(ctx context.Context, req *emptypb.Empty) (*guestProto.GetAllSpecialOccasionsResponse, error) {
	res, err := s.Repository.GetAllSpecialOccasions(ctx, req)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (s *ReservationSpecialOccasionService) GetWidgetAllSpecialOccasions(ctx context.Context, req *guestProto.GetWidgetAllSpecialOccasionsRequest) (*guestProto.GetAllSpecialOccasionsResponse, error) {
	res, err := s.Repository.GetWidgetAllSpecialOccasions(ctx, req)
	if err != nil {
		return nil, err
	}

	return res, nil
}
