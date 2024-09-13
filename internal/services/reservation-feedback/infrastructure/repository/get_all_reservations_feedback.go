package repository

import (
	"context"
	"net/http"
	"strings"

	"github.com/goplaceapp/goplace-common/pkg/meta"
	guestProto "github.com/goplaceapp/goplace-guest/api/v1"
	common "github.com/goplaceapp/goplace-guest/internal/services/common/converters"
	"github.com/goplaceapp/goplace-guest/internal/services/reservation-feedback/domain"
	"github.com/goplaceapp/goplace-guest/pkg/guestservice/converters"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"
)

func (r *ReservationFeedbackRepository) GetAllReservationsFeedbacks(ctx context.Context, req *guestProto.GetAllReservationsFeedbacksRequest) (*guestProto.GetAllReservationsFeedbacksResponse, error) {
	var (
		feedbacks       []*domain.ReservationFeedback
		totalRows       int64
		totalPositive   int64
		totalNegative   int64
		totalPending    int64
		totalSolved     int64
		pendingStatusId int
		solvedStatusId  int
		userBranches    = []int32{}
	)

	if err := r.QueryBuilder(ctx, req).Count(&totalRows).Error; err != nil {
		return nil, status.Error(http.StatusNotFound, err.Error())
	}

	offset := (req.GetPaginationParams().GetCurrentPage() - 1) * req.GetPaginationParams().GetPerPage()
	if err := r.QueryBuilder(ctx, req).
		Offset(int(offset)).
		Limit(int(req.GetPaginationParams().GetPerPage())).
		Order("reservation_feedbacks.id DESC").
		Find(&feedbacks).Error; err != nil {
		return nil, status.Error(http.StatusNotFound, err.Error())
	}

	for i := range feedbacks {
		feedbackData, err := r.GetReservationFeedbackData(ctx, feedbacks[i])
		if err != nil {
			return nil, status.Error(http.StatusInternalServerError, err.Error())
		}

		feedbacks[i] = feedbackData
	}

	if len(req.GetBranchIds()) > 0 {
		for _, branchId := range req.BranchIds {
			if r.userClient.Client.UserService.Repository.CheckForBranchAccess(ctx, branchId) {
				userBranches = append(userBranches, branchId)
			}
		}
	} else {
		userBranches = r.userClient.Client.UserService.Repository.GetAllUserBranchesIDs(ctx)
	}

	if err := r.GetTenantDBConnection(ctx).
		Model(&domain.ReservationFeedback{}).
		Joins("JOIN reservations on reservation_feedbacks.reservation_id = reservations.id").
		Where("reservation_feedbacks.rate < 3 AND "+
			"reservations.deleted_at IS NULL").
		Where("reservations.branch_id IN ?", userBranches).
		Count(&totalNegative).Error; err != nil {
		return nil, status.Error(http.StatusInternalServerError, err.Error())
	}

	if err := r.GetTenantDBConnection(ctx).
		Model(&domain.ReservationFeedback{}).
		Joins("JOIN reservations on reservation_feedbacks.reservation_id = reservations.id").
		Where("reservation_feedbacks.rate > 3 AND "+
			"reservations.deleted_at IS NULL").
		Where("reservations.branch_id IN ?", userBranches).
		Count(&totalPositive).Error; err != nil {
		return nil, status.Error(http.StatusInternalServerError, err.Error())
	}

	for i, v := range FeedbackStatuses {
		if v == meta.Pending {
			pendingStatusId = i + 1
		}
		if v == meta.Solved {
			solvedStatusId = i + 1
		}
	}

	if err := r.GetTenantDBConnection(ctx).
		Model(&domain.ReservationFeedback{}).
		Joins("JOIN reservations on reservation_feedbacks.reservation_id = reservations.id").
		Where("reservation_feedbacks.status_id = ? AND reservations.branch_id IN (?)",
			pendingStatusId,
			userBranches).
		Count(&totalPending).Error; err != nil {
		return nil, status.Error(http.StatusInternalServerError, err.Error())
	}

	if err := r.GetTenantDBConnection(ctx).
		Model(&domain.ReservationFeedback{}).
		Joins("JOIN reservations on reservation_feedbacks.reservation_id = reservations.id").
		Where("reservation_feedbacks.status_id = ? AND reservations.branch_id IN (?)",
			solvedStatusId,
			userBranches).
		Count(&totalSolved).Error; err != nil {
		return nil, status.Error(http.StatusInternalServerError, err.Error())
	}

	return &guestProto.GetAllReservationsFeedbacksResponse{
		Pagination:    common.BuildPaginationResponse(req.GetPaginationParams(), int32(totalRows), offset),
		Result:        converters.BuildAllReservationFeedbacksResponse(feedbacks),
		TotalPositive: int32(totalPositive),
		TotalNegative: int32(totalNegative),
		TotalPending:  int32(totalPending),
		TotalSolved:   int32(totalSolved),
	}, nil
}

func (r *ReservationFeedbackRepository) QueryBuilder(ctx context.Context, req *guestProto.GetAllReservationsFeedbacksRequest) *gorm.DB {
	queryBuilder := r.GetTenantDBConnection(ctx).
		Model(&domain.ReservationFeedback{}).
		Joins("JOIN reservations on reservation_feedbacks.reservation_id = reservations.id").
		Joins("JOIN guests on reservations.guest_id = guests.id").
		Joins("LEFT JOIN reservation_tables on reservations.id = reservation_tables.reservation_id").
		Joins("LEFT JOIN tables on reservation_tables.table_id = tables.id").
		Where("reservations.branch_id IN (?) AND reservations.deleted_at IS NULL",
			r.userClient.Client.UserService.Repository.GetAllUserBranchesIDs(ctx))

	// Adding Branch ID condition
	if len(req.GetBranchIds()) > 0 {
		queryBuilder = queryBuilder.Where("reservations.branch_id IN ?", req.GetBranchIds())
	}

	// Adding Status ID condition
	if len(req.GetStatusIds()) > 0 {
		queryBuilder = queryBuilder.Where("reservation_feedbacks.status_id IN ?", req.GetStatusIds())
	}

	// Adding Date Range condition
	fromDate := req.GetFromDate()
	if fromDate == "" {
		fromDate = "1980-01-01"
	}
	toDate := req.GetToDate()
	if toDate == "" {
		toDate = meta.InfiniteDate
	}
	queryBuilder = queryBuilder.Where("TO_CHAR(reservations.date, 'YYYY-MM-DD') BETWEEN ? AND ?", fromDate, toDate)

	// Adding Rate condition
	if len(req.GetRate()) > 0 {
		rates := req.GetRate()
		for i := range rates {
			rates[i] = strings.ToUpper(rates[i])
		}
		rateCase := "CASE " +
			"WHEN reservation_feedbacks.rate < 3 THEN 'negative' " +
			"WHEN reservation_feedbacks.rate > 3 THEN 'positive' " +
			"ELSE 'neutral' " +
			"END"
		queryBuilder = queryBuilder.Where("UPPER("+rateCase+") IN ?", rates)
	}

	// Adding Search Query conditions
	searchQuery := "%" + req.GetQuery() + "%"
	if req.GetQuery() != "" {
		queryBuilder = queryBuilder.Where("(UPPER(guests.first_name) LIKE UPPER(?) OR "+
			"UPPER(guests.last_name) LIKE UPPER(?) OR "+
			"guests.phone_number LIKE ? OR "+
			"tables.table_number LIKE ? OR "+
			"TO_CHAR(reservations.time, 'HH24:MI:SS') LIKE ?)",
			searchQuery, searchQuery, searchQuery, searchQuery, searchQuery)
	}

	return queryBuilder.Group("reservation_feedbacks.id").Distinct()
}
