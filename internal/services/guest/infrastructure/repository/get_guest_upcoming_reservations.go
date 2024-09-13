package repository

import (
	"context"
	"net/http"
	"time"

	"github.com/goplaceapp/goplace-guest/pkg/reservationservice/domain"

	guestProto "github.com/goplaceapp/goplace-guest/api/v1"
	"github.com/goplaceapp/goplace-guest/pkg/guestservice/converters"
	"google.golang.org/grpc/status"
)

func (r *GuestRepository) GetGuestUpcomingReservations(ctx context.Context, req *guestProto.GetGuestUpcomingReservationsRequest) (*guestProto.GetGuestUpcomingReservationsResponse, error) {
	var upcomingReservations []*domain.Reservation

	convertedTime := r.CommonRepository.ConvertToLocalTime(ctx, time.Now())
	currentDate := convertedTime.Format(time.DateOnly)
	currentTime := convertedTime.Format(time.TimeOnly)
	branchId := r.UserClient.Client.UserService.Repository.GetCurrentBranchId(ctx)

	if err := r.GetTenantDBConnection(ctx).
		Model(&domain.Reservation{}).
		Where("guest_id = ? AND "+
			"branch_id = ? AND "+
			"date >= ? AND "+
			"(date > ? OR (date = ? AND time >= ?))",
			req.GetGuestId(),
			branchId,
			currentDate,
			currentDate,
			currentDate,
			currentTime,
		).
		Order("date asc, time asc").
		Find(&upcomingReservations).Error; err != nil {
		if err.Error() == "record not found" {
			return nil, status.Error(http.StatusNotFound, "No upcoming reservation found")
		}
		return nil, err
	}

	var reservations []*domain.Reservation
	for _, reservation := range upcomingReservations {
		upcomingReservation, err := r.CommonRepository.GetReservationByID(ctx, reservation.ID)
		if err != nil {
			return nil, err
		}

		reservations = append(reservations, upcomingReservation)
	}

	return &guestProto.GetGuestUpcomingReservationsResponse{
		Result: converters.BuildAllReservationsResponse(reservations),
	}, nil
}
