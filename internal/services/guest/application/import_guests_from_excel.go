package application

import (
	"context"
	guestProto "github.com/goplaceapp/goplace-guest/api/v1"
)

func (s *GuestService) ImportGuestsFromExcel(ctx context.Context, req *guestProto.ImportGuestsFromExcelRequest) (*guestProto.ImportGuestsFromExcelResponse, error) {
	res, err := s.Repository.ImportGuestsFromExcel(ctx, req)
	if err != nil {
		return nil, err
	}

	return res, nil
}
