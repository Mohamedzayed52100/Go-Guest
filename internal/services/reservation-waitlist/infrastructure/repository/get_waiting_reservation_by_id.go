package repository

import (
	"context"

	"github.com/goplaceapp/goplace-guest/internal/services/reservation-waitlist/domain"
)

func (r *ReservationWaitListRepository) GetWaitingReservationByID(ctx context.Context, id int32) (*domain.ReservationWaitlist, error) {
	var res domain.ReservationWaitlist

	if err := r.GetTenantDBConnection(ctx).
		Where("id = ?", id).
		First(&res).Error; err != nil {
		return nil, err
	}

	return &res, nil
}
