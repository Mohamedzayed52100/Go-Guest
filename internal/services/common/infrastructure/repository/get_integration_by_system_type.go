package common

import (
	"context"

	"github.com/goplaceapp/goplace-guest/internal/services/reservation/domain"
)

func (r *CommonRepository) GetIntegrationBySystemType(ctx context.Context, systemType string, branchID int32) (*domain.Integration, error) {
	var record domain.Integration

	err := r.GetTenantDBConnection(ctx).
		Model(&record).
		Where("system_type = ? AND branch_id = ?", systemType, branchID).
		First(&record).Error
	if err != nil {
		return nil, err
	}

	return &record, nil
}
