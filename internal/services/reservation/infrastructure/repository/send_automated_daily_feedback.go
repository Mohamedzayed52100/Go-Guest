package repository

import (
	"context"
	"time"

	"github.com/goplaceapp/goplace-common/pkg/logger"
	"github.com/goplaceapp/goplace-common/pkg/meta"
)

func (r *ReservationRepository) SendAutomatedDailyFeedback(ctx context.Context, branchId int32) (res bool, err error) {
	var reservationIds []int32

	if err := r.GetTenantDBConnection(ctx).
		Table("reservations").
		Joins("JOIN reservation_statuses ON reservation_statuses.id = reservations.status_id").
		Where("reservation_statuses.type = ? AND "+
			"reservations.branch_id = ? AND "+
			"reservations.date = ?",
			meta.Left,
			branchId,
			time.Now().AddDate(0, 0, -1).Format(time.DateOnly),
		).
		Pluck("reservations.id", &reservationIds).
		Error; err != nil {
		logger.Default().Errorf("Failed to get reservation ids: %v", err)
		return false, err
	}

	if _, err := r.SendBulkReservationWhatsappFeedback(ctx, reservationIds, branchId); err != nil {
		logger.Default().Errorf("Failed to send bulk reservation whatsapp feedback: %v", err)
		return false, err
	}

	return true, nil
}
