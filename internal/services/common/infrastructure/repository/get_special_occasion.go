package common

import (
	"context"

	"github.com/goplaceapp/goplace-guest/internal/services/reservation-special-occasion/domain"
)

func (r *CommonRepository) GetReservationSpecialOccasionByID(ctx context.Context, specialOccasionID int32) (*domain.SpecialOccasion, error) {
	var specialOccasion domain.SpecialOccasion

	if err := r.GetTenantDBConnection(ctx).
		Model(&domain.SpecialOccasion{}).
		Where("id = ?", specialOccasionID).
		First(&specialOccasion).Error; err != nil {
		return nil, err
	}

	return &specialOccasion, nil
}
