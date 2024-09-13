package repository

import (
	"context"
	"github.com/goplaceapp/goplace-guest/internal/services/reservation/domain"
)

func (r *ReservationWaitListRepository) GetWaitingReservationTags(ctx context.Context, reservationID int32) ([]*domain.ReservationTag, error) {
	var tags []*domain.ReservationTag

	if err := r.GetTenantDBConnection(ctx).
		Model(&tags).
		Joins("JOIN reservation_waitlist_tags_assignments ON reservation_waitlist_tags_assignments.tag_id = reservation_tags.id").
		Where("reservation_waitlist_tags_assignments.reservation_id = ?", reservationID).
		Distinct().
		Find(&tags).Error; err != nil {
		return []*domain.ReservationTag{}, nil
	}

	for _, tag := range tags {
		var category domain.ReservationTagCategory

		if err := r.GetTenantDBConnection(ctx).
			Model(&domain.ReservationTagCategory{}).
			Where("id = ?", tag.CategoryID).
			First(&category).Error; err != nil {
			return nil, err
		}

		tag.Category = &category
	}

	return tags, nil
}
