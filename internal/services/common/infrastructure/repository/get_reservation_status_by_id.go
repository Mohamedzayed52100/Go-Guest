package common

import (
	"context"

	"github.com/goplaceapp/goplace-guest/pkg/reservationservice/domain"
)

func (r *CommonRepository) GetReservationStatusByID(ctx context.Context, id int32, branchId int32) (*domain.ReservationStatus, error) {
	var result *domain.ReservationStatus

	if err := r.GetTenantDBConnection(ctx).
		Model(&domain.ReservationStatus{}).
		Where("id = ? AND branch_id = ?", id, branchId).
		First(&result).Error; err != nil {
		return nil, err
	}

	return result, nil
}
