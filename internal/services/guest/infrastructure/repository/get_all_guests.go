package repository

import (
	"context"
	"net/http"
	"strings"

	"github.com/goplaceapp/goplace-common/pkg/errorhelper"
	guestProto "github.com/goplaceapp/goplace-guest/api/v1"
	common "github.com/goplaceapp/goplace-guest/internal/services/common/converters"
	"github.com/goplaceapp/goplace-guest/pkg/guestservice/converters"
	"github.com/goplaceapp/goplace-guest/pkg/guestservice/domain"
	"google.golang.org/grpc/status"
)

func (r *GuestRepository) GetAllGuests(ctx context.Context, req *guestProto.GetAllGuestsRequest) (*guestProto.GetAllGuestsResponse, error) {
	var (
		guests    []*domain.Guest
		totalRows int64
		offset    = (req.GetPaginationParams().GetCurrentPage() - 1) * req.GetPaginationParams().GetPerPage()
	)

	req.Query = strings.ReplaceAll(req.GetQuery(), " ", "")

	if err := r.GetTenantDBConnection(ctx).Model(&domain.Guest{}).
		Where("UPPER(CONCAT(first_name,last_name)) LIKE UPPER(?)", "%"+req.GetQuery()+"%").
		Or("UPPER(email) LIKE UPPER(?)", "%"+req.GetQuery()+"%").
		Or("phone_number LIKE ?", "%"+req.GetQuery()+"%").
		Count(&totalRows).Error; err != nil {
		return nil, status.Error(http.StatusNotFound, errorhelper.ErrNoGuestsFound)
	}

	if err := r.GetTenantDBConnection(ctx).
		Model(&domain.Guest{}).
		Offset(int(offset)).
		Limit(int(req.GetPaginationParams().GetPerPage())).
		Where("UPPER(CONCAT(first_name,last_name)) LIKE UPPER(?)", "%"+req.GetQuery()+"%").
		Or("UPPER(email) LIKE UPPER(?)", "%"+req.GetQuery()+"%").
		Or("phone_number LIKE ?", "%"+req.GetQuery()+"%").
		Order("id desc").
		Find(&guests).Error; err != nil {
		return &guestProto.GetAllGuestsResponse{
			Pagination: common.BuildPaginationResponse(req.GetPaginationParams(), int32(totalRows), offset),
			Result:     []*guestProto.Guest{},
		}, nil
	}

	for i := range guests {
		getGuestByID, err := r.CommonRepository.GetAllGuestData(ctx, guests[i])
		if err != nil {
			return nil, status.Error(http.StatusInternalServerError, errorhelper.ErrGuestNotFound)
		}

		guests[i] = getGuestByID
	}

	return &guestProto.GetAllGuestsResponse{
		Pagination: common.BuildPaginationResponse(req.GetPaginationParams(), int32(totalRows), offset),
		Result:     converters.BuildAllGuestsResponse(guests),
	}, nil
}
