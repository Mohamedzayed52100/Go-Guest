package repository

import (
	"context"
	"net/http"
	"sort"
	"strings"
	"sync"
	"time"

	domain2 "github.com/goplaceapp/goplace-guest/pkg/reservationservice/domain"

	guestProto "github.com/goplaceapp/goplace-guest/api/v1"
	common "github.com/goplaceapp/goplace-guest/internal/services/common/converters"
	"github.com/goplaceapp/goplace-guest/pkg/guestservice/converters"
	"github.com/goplaceapp/goplace-guest/utils"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"

	"github.com/goplaceapp/goplace-common/pkg/errorhelper"
)

var (
	joinedGuests            = false
	joinedReservationTables = false
)

func (r *ReservationRepository) GetAllReservations(ctx context.Context, req *guestProto.GetAllReservationsRequest) (*guestProto.GetAllReservationsResponse, error) {
	var (
		reservations []*domain2.Reservation
		totalRows    int64
		userRepo     = r.userClient.Client.UserService.Repository
	)

	if req.GetBranchId() != 0 && !userRepo.CheckForBranchAccess(ctx, req.GetBranchId()) {
		return nil, status.Error(http.StatusNotFound, "You don't have access to this branch")
	}

	baseQuery := r.GetTenantDBConnection(ctx).Model(&domain2.Reservation{})

	if req.GetBranchId() == 0 {
		baseQuery = baseQuery.Where("branch_id = ?", userRepo.GetCurrentBranchId(ctx))
	}

	offset := (req.GetPaginationParams().GetCurrentPage() - 1) * req.GetPaginationParams().GetPerPage()
	if err := reservationsQueryBuilder(baseQuery, req).
		Group("reservations.id").
		Count(&totalRows).Error; err != nil {
		joinedGuests = false
		joinedReservationTables = false
		return &guestProto.GetAllReservationsResponse{
			Pagination: common.BuildPaginationResponse(req.GetPaginationParams(), int32(totalRows), offset),
			Result:     []*guestProto.Reservation{},
		}, nil
	}

	if err := reservationsQueryBuilder(baseQuery, req).
		Group("reservations.id").
		Order("reservations.date desc, reservations.time asc").
		Select("reservations.*").
		Offset(int(offset)).
		Limit(int(req.GetPaginationParams().GetPerPage())).
		Find(&reservations).Error; err != nil {
		joinedGuests = false
		joinedReservationTables = false
		return &guestProto.GetAllReservationsResponse{
			Pagination: common.BuildPaginationResponse(req.GetPaginationParams(), int32(totalRows), offset),
			Result:     []*guestProto.Reservation{},
		}, nil
	}

	wg := sync.WaitGroup{}
	wg.Add(len(reservations))
	errChan := make(chan error, len(reservations))
	for i := range reservations {
		i := i
		go func() {
			defer wg.Done()
			var err error
			reservations[i], err = r.CommonRepo.GetAllReservationData(ctx, reservations[i])
			if err != nil {
				joinedGuests = false
				joinedReservationTables = false
				errChan <- status.Error(http.StatusInternalServerError, "Error getting reservation data: "+err.Error())
				return
			}
			errChan <- nil
		}()
		err := <-errChan
		if err != nil {
			return nil, err
		}
	}
	wg.Wait()
	close(errChan)

	for err := range errChan {
		if err != nil {
			return nil, err
		}
	}

	sort.SliceStable(reservations, func(i, j int) bool {
		time1, _ := time.Parse(time.TimeOnly, reservations[i].Time)
		time2, _ := time.Parse(time.TimeOnly, reservations[j].Time)

		if utils.IsAfterMidnight(time1) != utils.IsAfterMidnight(time2) {
			return true
		}

		return utils.CompareTimes(time1, time2)
	})

	joinedGuests = false
	joinedReservationTables = false

	return &guestProto.GetAllReservationsResponse{
		Pagination: common.BuildPaginationResponse(req.GetPaginationParams(), int32(totalRows), offset),
		Result:     converters.BuildAllReservationsResponse(reservations),
	}, nil
}

func reservationsQueryBuilder(query *gorm.DB, req *guestProto.GetAllReservationsRequest) *gorm.DB {
	if req.GetDate() != "" {
		query = query.Where("TO_CHAR(reservations.date, 'YYYY-MM-DD') LIKE ?", "%"+req.GetDate()+"%")
	}

	if req.GetBranchId() != 0 {
		query = query.Where("reservations.branch_id = ?", req.GetBranchId())
	}

	if len(req.GetStatusIds()) > 0 {
		query = query.Where("reservations.status_id IN (?)", req.GetStatusIds())
	}

	if req.GetShiftId() != 0 {
		query = query.Where("reservations.shift_id = ?", req.GetShiftId())
	}

	if len(req.GetTableIds()) > 0 && !joinedReservationTables {
		query = query.
			Joins("JOIN reservation_tables on reservations.id = reservation_tables.reservation_id").
			Where("reservation_tables.table_id IN (?)", req.GetTableIds())

		joinedReservationTables = true
	}

	if req.GetQuery() != "" {
		req.Query = strings.ReplaceAll(req.GetQuery(), " ", "")
		searchQuery := "%" + req.GetQuery() + "%"
		if !joinedGuests {
			query = query.Joins("JOIN guests on reservations.guest_id = guests.id")
			joinedGuests = true
		}
		query = query.Where("UPPER(CONCAT(guests.first_name,guests.last_name)) LIKE UPPER(?) OR "+
			"reservations.reservation_ref LIKE UPPER(?) OR "+
			"guests.phone_number LIKE ? OR "+
			"TO_CHAR(reservations.time, 'HH24:MI:SS') LIKE ?",
			searchQuery,
			searchQuery,
			searchQuery,
			searchQuery,
		)
	}

	return query
}

func (r *ReservationRepository) GetReservationByID(ctx context.Context, req *guestProto.GetReservationByIDRequest) (*guestProto.GetReservationByIDResponse, error) {
	var getReservation *domain2.Reservation

	if err := r.GetTenantDBConnection(ctx).
		First(&getReservation, "id = ?", req.GetId()).Error; err != nil {
		return nil, err
	}

	if getReservation.BranchID != r.userClient.Client.UserService.Repository.GetCurrentBranchId(ctx) {
		return nil, status.Error(http.StatusNotFound, errorhelper.ErrReservationNotFound)
	}

	var err error
	getReservation, err = r.CommonRepo.GetAllReservationData(ctx, getReservation)
	if err != nil {
		return nil, err
	}

	return &guestProto.GetReservationByIDResponse{
		Result: converters.BuildReservationProto(getReservation),
	}, nil
}
