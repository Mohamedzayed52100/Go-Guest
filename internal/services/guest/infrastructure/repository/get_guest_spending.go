package repository

import (
	"context"
	domain2 "github.com/goplaceapp/goplace-guest/pkg/reservationservice/domain"
	"time"

	guestProto "github.com/goplaceapp/goplace-guest/api/v1"
	"github.com/goplaceapp/goplace-guest/pkg/guestservice/converters"
	"github.com/goplaceapp/goplace-guest/internal/services/reservation/domain"
)

func (r *GuestRepository) GetGuestSpending(ctx context.Context, req *guestProto.GetGuestSpendingRequest) (*guestProto.GetGuestSpendingResponse, error) {
	var (
		years        map[int]map[string]float32
		reservations = []*domain2.Reservation{}
	)

	currentMonth := time.Date(
		time.Now().Year(), time.Now().Month(), 1, 0, 0, 0, 0, time.UTC,
	)
	nextMonth := time.Date(
		time.Now().Year(), time.Now().Month()+1, 1, 0, 0, 0, 0, time.UTC,
	)

	r.GetTenantDBConnection(ctx).Find(&reservations,
		"guest_id = ? AND date >= ? AND date < ?",
		req.GetGuestId(),
		currentMonth,
		nextMonth,
	)

	for _, res := range reservations {
		var totalSpent float64

		year := res.Date.Year()
		month := res.Date.Month().String()

		err := r.GetTenantDBConnection(ctx).
			Model(&domain.ReservationOrder{}).
			Select("final_total").
			Where("reservation_id = ?", res.ID).
			Scan(&totalSpent).
			Error
		if err != nil {
			totalSpent = 0
		}

		if years == nil {
			years = make(map[int]map[string]float32)
		}
		if years[year] == nil {
			years[year] = make(map[string]float32)
		}

		years[year][month] += float32(totalSpent)
	}

	for year, months := range years {
		for i := 1; i <= 12; i++ {
			if _, ok := months[time.Month(i).String()]; !ok {
				years[year][time.Month(i).String()] = 0
			}
		}
	}

	return &guestProto.GetGuestSpendingResponse{
		Result: converters.BuildGuestSpendingResponse(years),
	}, nil

}
