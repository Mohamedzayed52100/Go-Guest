package repository

import (
	"context"
	"net/http"

	guestProto "github.com/goplaceapp/goplace-guest/api/v1"
	"github.com/goplaceapp/goplace-guest/internal/services/reservation-special-occasion/adapters/converters"
	"github.com/goplaceapp/goplace-guest/internal/services/reservation-special-occasion/domain"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (r *ReservationSpecialOccasionRepository) GetAllSpecialOccasions(ctx context.Context, req *emptypb.Empty) (*guestProto.GetAllSpecialOccasionsResponse, error) {
	result := []*domain.SpecialOccasion{}
	userRepo := r.userClient.Client.UserService.Repository
	currentBranch := userRepo.GetCurrentBranchId(ctx)

	r.GetTenantDBConnection(ctx).
		Find(&result, "branch_id = ?", currentBranch)
	return &guestProto.GetAllSpecialOccasionsResponse{
		Result: converters.BuildAllSpecialOccasionsResponse(result),
	}, nil
}

func (r *ReservationSpecialOccasionRepository) GetWidgetAllSpecialOccasions(ctx context.Context, req *guestProto.GetWidgetAllSpecialOccasionsRequest) (*guestProto.GetAllSpecialOccasionsResponse, error) {
	result := []*domain.SpecialOccasion{}

	if req.GetBranchId() == 0 {
		return nil, status.Error(http.StatusBadRequest, "Branch id query is required")
	}

	r.GetTenantDBConnection(ctx).
		Find(&result, "branch_id = ?", req.GetBranchId())

	return &guestProto.GetAllSpecialOccasionsResponse{
		Result: converters.BuildAllSpecialOccasionsResponse(result),
	}, nil
}
