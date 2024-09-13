package repository

import (
	"context"
	"github.com/goplaceapp/goplace-guest/pkg/reservationservice/domain"
)

func (r *ReservationRepository) GetReservationStatusByName(ctx context.Context, status string, branchId int32) (*domain.ReservationStatus, error) {
	var result *domain.ReservationStatus

	if err := r.GetTenantDBConnection(ctx).
		Model(&domain.ReservationStatus{}).
		Where("name = ? AND branch_id = ?", status, branchId).
		First(&result).Error; err != nil {
		return nil, err
	}

	return result, nil
}
