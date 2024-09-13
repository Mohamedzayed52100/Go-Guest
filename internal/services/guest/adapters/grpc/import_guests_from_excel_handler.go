package grpc

import (
	"context"
	guestProto "github.com/goplaceapp/goplace-guest/api/v1"
)

func (s *GuestServiceServer) ImportGuestsFromExcel(ctx context.Context, req *guestProto.ImportGuestsFromExcelRequest) (*guestProto.ImportGuestsFromExcelResponse, error) {
	res, err := s.guestService.ImportGuestsFromExcel(ctx, req)
	if err != nil {
		return nil, err
	}

	return res, nil
}
