package repository

import (
	"context"
	"net/http"
	"time"

	"github.com/goplaceapp/goplace-guest/internal/services/day-operations/domain"
	reservationDomain "github.com/goplaceapp/goplace-guest/pkg/reservationservice/domain"
	externalWaitlistDomain "github.com/goplaceapp/goplace-guest/pkg/waitlistservice/domain"
	userDomain "github.com/goplaceapp/goplace-user/pkg/userservice/domain"

	"github.com/goplaceapp/goplace-common/pkg/meta"
	guestProto "github.com/goplaceapp/goplace-guest/api/v1"
	waitlistDomain "github.com/goplaceapp/goplace-guest/internal/services/reservation-waitlist/domain"
	roleDomain "github.com/goplaceapp/goplace-user/pkg/roleservice/domain"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"
)

func (r *DayOperationsRepository) CloseDayOperations(ctx context.Context, req *guestProto.CloseDayOperationsRequest) (*guestProto.CloseDayOperationsResponse, error) {
	var (
		isOpen               int64
		leftReservations     = []*reservationDomain.Reservation{}
		upcomingReservations = []*reservationDomain.Reservation{}
		waitList             = []*waitlistDomain.ReservationWaitlist{}
	)

	isSuperAdmin := r.userClient.Client.UserService.Repository.CheckAdminPinCode(ctx, req.GetPinCode())

	currentUser, err := r.userClient.Client.UserService.Repository.GetLoggedInUser(ctx)
	if err != nil {
		return nil, status.Error(http.StatusInternalServerError, err.Error())
	}

	currentTime := r.CommonRepo.ConvertToLocalTime(ctx, time.Now())

	if req.GetDate() == "" {
		return nil, status.Error(http.StatusBadRequest, "Date is required")
	} else if req.GetDate() > currentTime.Format("2006-01-02") {
		return nil, status.Error(http.StatusBadRequest, "Cannot close future dates")
	}

	var endShiftPermission string
	if err := r.GetTenantDBConnection(ctx).
		Model(&roleDomain.RolePermissionAssignment{}).
		Joins("JOIN permissions ON permissions.id = role_permission_assignments.permission_id").
		Where(`role_id = ? AND permissions.name = 'end-day.shifts'`, currentUser.RoleID).
		Select("permissions.name").
		Scan(&endShiftPermission).Error; err != nil {
		return nil, status.Error(http.StatusInternalServerError, err.Error())
	}
	if endShiftPermission == "" && !isSuperAdmin {
		return nil, status.Error(http.StatusForbidden, "You do not have permission to close day operations")
	}

	noShowStatus, err := r.reservationRepository.GetReservationStatusByName(ctx, meta.NoShow, currentUser.BranchID)
	if err != nil {
		return nil, status.Error(http.StatusInternalServerError, err.Error())
	}
	leftStatus, err := r.reservationRepository.GetReservationStatusByName(ctx, meta.Left, currentUser.BranchID)
	if err != nil {
		return nil, status.Error(http.StatusInternalServerError, err.Error())
	}
	cancelledStatus, err := r.reservationRepository.GetReservationStatusByName(ctx, meta.Cancelled, currentUser.BranchID)
	if err != nil {
		return nil, status.Error(http.StatusInternalServerError, err.Error())
	}

	excludedStatuses := []int{noShowStatus.ID, leftStatus.ID, cancelledStatus.ID}
	if err := r.GetTenantDBConnection(ctx).
		Model(&reservationDomain.Reservation{}).
		Where("date = ? AND branch_id = ? AND status_id NOT IN (?)",
			req.GetDate(),
			currentUser.BranchID,
			excludedStatuses,
		).
		Count(&isOpen).Error; err != nil {
		return nil, status.Error(http.StatusInternalServerError, err.Error())
	}

	if isOpen == 0 {
		r.GetTenantDBConnection(ctx).
			Where("branch_id = ?", currentUser.BranchID).
			Find(&waitList)
		for _, w := range waitList {
			createdAt := r.CommonRepo.ConvertToLocalTime(ctx, w.CreatedAt.UTC())
			currentTime := r.CommonRepo.ConvertToLocalTime(ctx, time.Now().UTC())
			if !createdAt.After(currentTime) {
				isOpen++
			}
		}
		if isOpen == 0 {
			return nil, status.Error(http.StatusBadRequest, "The day has already ended")
		}
	}

	baseQuery := r.GetTenantDBConnection(ctx).
		Model(&reservationDomain.Reservation{}).
		Joins("JOIN reservation_statuses ON reservation_statuses.id = reservations.status_id").
		Where("reservations.branch_id = ? AND TO_CHAR(reservations.date, 'YYYY-MM-DD') <= ?",
			currentUser.BranchID,
			req.GetDate(),
		)

	if err := baseQuery.
		Where("reservation_statuses.type = 'Seated' OR reservation_statuses.type = 'Arrived'").
		Find(&leftReservations).Error; err != nil {
		return nil, status.Error(http.StatusInternalServerError, err.Error())
	}

	for _, leftReservation := range leftReservations {
		_, err := r.reservationRepository.UpdateReservation(ctx, &guestProto.UpdateReservationRequest{
			Params: &guestProto.ReservationParams{
				Id:       leftReservation.ID,
				StatusId: int32(leftStatus.ID),
			},
		})
		if err != nil {
			return nil, err
		}
	}

	baseQuery = r.GetTenantDBConnection(ctx).
		Model(&reservationDomain.Reservation{}).
		Joins("JOIN reservation_statuses ON reservation_statuses.id = reservations.status_id").
		Where("reservations.branch_id = ? AND TO_CHAR(reservations.date, 'YYYY-MM-DD') <= ?",
			currentUser.BranchID,
			req.GetDate(),
		)

	// Update upcoming reservations to No Show
	if err := baseQuery.
		Where("reservation_statuses.type = 'Upcoming'").
		Find(&upcomingReservations).Error; err != nil {
		return nil, status.Error(http.StatusInternalServerError, err.Error())
	}

	for _, upcomingReservation := range upcomingReservations {
		_, err := r.reservationRepository.UpdateReservation(ctx, &guestProto.UpdateReservationRequest{
			Params: &guestProto.ReservationParams{
				Id:       upcomingReservation.ID,
				StatusId: int32(noShowStatus.ID),
			},
		})
		if err != nil {
			return nil, err
		}
	}

	// Delete all waiting reservations
	if err := r.GetTenantDBConnection(ctx).Transaction(func(tx *gorm.DB) error {
		waitlist := []*waitlistDomain.ReservationWaitlist{}

		tx.Model(&waitlistDomain.ReservationWaitlist{}).
			Where("branch_id = ? AND created_at <= ?",
				currentUser.BranchID,
				time.Now().Format("2006-01-02")).
			Find(&waitlist)

		for _, w := range waitlist {
			if err := r.GetTenantDBConnection(ctx).
				Delete(&waitlistDomain.ReservationWaitlistLog{}, "reservation_waitlist_id = ?", w.ID).
				Error; err != nil {
				return status.Error(http.StatusInternalServerError, err.Error())
			}

			if err := tx.Delete(&externalWaitlistDomain.ReservationWaitlistTagsAssignment{},
				"reservation_id = ?", w.ID).
				Error; err != nil {
				return status.Error(http.StatusInternalServerError, err.Error())
			}
		}

		for _, w := range waitlist {
			if err := r.GetTenantDBConnection(ctx).Delete(&w).Error; err != nil {
				return status.Error(http.StatusInternalServerError, err.Error())
			}
		}

		return nil
	}); err != nil {
		return nil, err
	}

	//// Get report data
	//reportData, err := s.GetReportData(ctx, currentUser, req.GetDate())
	//if err != nil {
	//	return nil, status.Error(http.StatusInternalServerError, err.Error())
	//}
	//
	//// Send report WA message to managers about the day's operations
	//if err := common.SendCloseDayWhatsappReport(reportData); err != nil {
	//	return nil, status.Error(http.StatusInternalServerError, err.Error())
	//}

	return &guestProto.CloseDayOperationsResponse{
		Code:    http.StatusOK,
		Message: "Day operations closed successfully",
	}, nil
}

func (r *DayOperationsRepository) GetReportData(ctx context.Context, currentUser *userDomain.User, date string) (*domain.CloseDayOperationsReport, error) {
	var (
		branchName       string
		leftReservations struct {
			Count int64
			Sum   int64
		}
		noShowReservations struct {
			Count int64
			Sum   int64
		}
		walkInReservations struct {
			Count int64
			Sum   int64
		}
		cancelledReservations struct {
			Count int64
			Sum   int64
		}
		sales struct {
			TotalSales                 float64
			AverageCheckPerReservation float64
			AverageCheckPerGuest       float64
		}
	)

	// Get branch name
	branch, err := r.userClient.Client.UserService.Repository.GetBranchByID(ctx, currentUser.BranchID)
	if err != nil {
		return nil, status.Error(http.StatusInternalServerError, err.Error())
	}
	branchName = branch.Name

	// Get total left reservations
	if err := r.GetTenantDBConnection(ctx).
		Model(&reservationDomain.Reservation{}).
		Select("COUNT(*) as count, SUM(guests_number) as sum").
		Where("date = ? AND "+
			"branch_id = ? AND "+
			"status_id IN (SELECT id FROM reservation_statuses WHERE category = 'in-service' AND name = ?)",
			date,
			currentUser.BranchID,
			meta.Left).
		Scan(&leftReservations).Error; err != nil {
		return nil, status.Error(http.StatusInternalServerError, err.Error())
	}

	// Get total no show reservations
	if err := r.GetTenantDBConnection(ctx).
		Model(&reservationDomain.Reservation{}).
		Select("COUNT(*) as count, SUM(guests_number) as sum").
		Where("date = ? AND "+
			"branch_id = ? AND "+
			"status_id IN "+
			"(SELECT reservation_statuses.id FROM reservation_statuses WHERE reservation_statuses.category = 'pre-service' AND "+
			"reservation_statuses.name = ?)",
			date,
			currentUser.BranchID,
			meta.NoShow).
		Scan(&noShowReservations).Error; err != nil {
		return nil, status.Error(http.StatusInternalServerError, err.Error())
	}

	// Get total walk in and direct-in reservations
	if err := r.GetTenantDBConnection(ctx).
		Model(&reservationDomain.Reservation{}).
		Select("COUNT(*) as count, SUM(guests_number) as sum").
		Where("date = ? AND "+
			"branch_id = ? AND "+
			"(reserved_via = 'Walked in' OR reserved_via = 'Direct in')",
			date,
			currentUser.BranchID,
		).
		Scan(&walkInReservations).Error; err != nil {
		return nil, status.Error(http.StatusInternalServerError, err.Error())
	}

	// Get total cancellation reservations
	if err := r.GetTenantDBConnection(ctx).
		Model(&reservationDomain.Reservation{}).
		Select("COUNT(*) as count, SUM(guests_number) as sum").
		Where("date = ? AND "+
			"branch_id = ? AND "+
			"status_id IN "+
			"(SELECT reservation_statuses.id FROM "+
			"reservation_statuses WHERE "+
			"reservation_statuses.category = 'pre-service' AND "+
			"reservation_statuses.name = ?)",
			date,
			currentUser.BranchID,
			meta.Cancelled).
		Scan(&cancelledReservations).Error; err != nil {
		return nil, status.Error(http.StatusInternalServerError, err.Error())
	}

	// Get total sales
	if err := r.GetTenantDBConnection(ctx).
		Model(&reservationDomain.Reservation{}).
		Joins("JOIN reservation_orders ON reservation_orders.reservation_id = reservations.id").
		Where("reservations.branch_id = ? AND TO_CHAR(reservations.date, 'YYYY-MM-DD') = ?", currentUser.BranchID, date).
		Select("SUM(reservation_orders.final_total) as total_sales, " +
			"(SUM(reservation_orders.final_total) / COUNT(reservations)) as average_check_per_reservation, " +
			"(SUM(reservation_orders.final_total) / SUM(reservations.guests_number)) as average_check_per_guest").
		Scan(&sales).Error; err != nil {
		return nil, status.Error(http.StatusInternalServerError, err.Error())
	}

	return &domain.CloseDayOperationsReport{
		BranchName:                 branchName,
		Date:                       date,
		TotalLeftReservations:      leftReservations.Count,
		TotalLeftGuests:            leftReservations.Sum,
		TotalNoShowReservations:    noShowReservations.Count,
		TotalNoShowGuests:          noShowReservations.Sum,
		TotalWalkInReservations:    walkInReservations.Count,
		TotalWalkInGuests:          walkInReservations.Sum,
		TotalCancelledReservations: cancelledReservations.Count,
		TotalCancelledGuests:       cancelledReservations.Sum,
		TotalSales:                 int64(sales.TotalSales),
		AverageCheckPerReservation: int64(sales.AverageCheckPerReservation),
		AverageCheckPerGuest:       int64(sales.AverageCheckPerGuest),
	}, nil
}
