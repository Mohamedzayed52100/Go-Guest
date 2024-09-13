package repository

import (
	"context"
	"github.com/goplaceapp/goplace-guest/pkg/reservationservice/domain"
)

func (r *ReservationRepository) GetReservationStatusByID(ctx context.Context, statusId int32, branchId int32) (*domain.ReservationStatus, error) {
	var result *domain.ReservationStatus

	if err := r.GetTenantDBConnection(ctx).
		Model(&domain.ReservationStatus{}).
		Where("id = ? AND branch_id = ?", statusId, branchId).
		First(&result).Error; err != nil {
		return nil, err
	}

	return result, nil
}
