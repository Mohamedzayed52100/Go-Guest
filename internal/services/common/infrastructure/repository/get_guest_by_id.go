package common

import (
	"context"
	"net/http"
	"time"

	domain2 "github.com/goplaceapp/goplace-guest/pkg/reservationservice/domain"

	"github.com/goplaceapp/goplace-common/pkg/meta"
	"github.com/goplaceapp/goplace-common/pkg/rbac"
	guestProto "github.com/goplaceapp/goplace-guest/api/v1"
	reservationFeedbackDomain "github.com/goplaceapp/goplace-guest/internal/services/reservation-feedback/domain"
	reservationDomain "github.com/goplaceapp/goplace-guest/internal/services/reservation/domain"
	"github.com/goplaceapp/goplace-guest/pkg/guestservice/converters"
	"github.com/goplaceapp/goplace-guest/pkg/guestservice/domain"
	tagDomain "github.com/goplaceapp/goplace-settings/pkg/guesttagservice/domain"
	"github.com/goplaceapp/goplace-settings/utils"
	"google.golang.org/grpc/status"
)

func (r *CommonRepository) GetGuestByID(ctx context.Context, req *guestProto.GetGuestByIDRequest) (*guestProto.GetGuestByIDResponse, error) {
	var fetchedGuest *domain.Guest

	currentUser, err := r.userClient.Client.UserService.Repository.GetAuthenticatedUser(ctx)
	if err != nil {
		return nil, status.Error(http.StatusInternalServerError, err.Error())
	}

	if err := r.GetTenantDBConnection(ctx).
		Where("id = ?", req.GetId()).
		First(&fetchedGuest).
		Error; err != nil {
		return nil, err
	}

	fetchedGuest, err = r.GetAllGuestData(ctx, fetchedGuest)
	if err != nil {
		return nil, err
	}

	if !utils.ArrayContains(currentUser.GetResult().GetRole().GetPermissions(), rbac.ViewGuestFinancials.Name) {
		fetchedGuest.TotalSpent = 0
	}

	return &guestProto.GetGuestByIDResponse{
		Result: converters.BuildGuestResponse(fetchedGuest),
	}, nil
}

func (r *CommonRepository) GetAllGuestData(ctx context.Context, guest *domain.Guest) (*domain.Guest, error) {
	if err := r.GetTenantDBConnection(ctx).
		First(&guest, "id = ?", guest.ID).
		Error; err != nil {
		return nil, status.Error(http.StatusInternalServerError, err.Error())
	}

	lastVisit, err := r.GetGuestLastVisit(ctx, guest.ID)
	if err != nil {
		guest.LastVisit = &time.Time{}
	}

	totalVisits, err := r.GetGuestTotalVisits(ctx, guest.ID)
	if err != nil {
		guest.TotalVisits = 0
	}

	currentMood, err := r.GetGuestCurrentMood(ctx, guest.ID)
	if err != nil {
		guest.CurrentMood = ""
	}

	totalSpent, err := r.GetGuestTotalSpent(ctx, guest.ID)
	if err != nil {
		guest.TotalSpent = 0
	}

	totalNoShow, err := r.GetGuestTotalReservationByStatus(ctx, guest.ID, meta.NoShow)
	if err != nil {
		guest.TotalNoShow = 0
	}

	totalCancel, err := r.GetGuestTotalReservationByStatus(ctx, guest.ID, meta.Cancelled)
	if err != nil {
		guest.TotalCancel = 0
	}

	branches, err := r.GetGuestBranches(ctx, guest.ID)
	if err != nil {
		guest.Branches = nil
	}

	tags, err := r.GetGuestTags(ctx, guest.ID)
	if err != nil {
		guest.Tags = nil
	}

	notes, err := r.GetGuestNotes(ctx, guest.ID)
	if err != nil {
		guest.Notes = nil
	}

	upcomingReservationDate, err := r.GetGuestUpcomingReservationDate(ctx, guest.ID)
	if err != nil {
		guest.UpcomingReservation = ""
	}

	guest.LastVisit = lastVisit
	guest.TotalVisits = totalVisits
	guest.CurrentMood = currentMood
	guest.TotalSpent = float32(totalSpent)
	guest.TotalNoShow = totalNoShow
	guest.TotalCancel = totalCancel
	guest.Branches = branches
	guest.Tags = tags
	guest.Notes = notes
	guest.UpcomingReservation = upcomingReservationDate

	return guest, nil
}

func (r *CommonRepository) GetGuestLastVisit(ctx context.Context, guestId int32) (*time.Time, error) {
	var lastVisit *time.Time

	if err := r.GetTenantDBConnection(ctx).
		Model(domain2.Reservation{}).
		Where("guest_id = ?", guestId).
		Select("date").
		Order("date DESC").
		Take(&lastVisit).
		Error; err != nil {
		return &time.Time{}, err
	}

	return lastVisit, nil
}

func (r *CommonRepository) GetGuestTotalVisits(ctx context.Context, guestId int32) (int32, error) {
	var reservations []*domain2.Reservation

	if err := r.GetTenantDBConnection(ctx).
		Model(&domain2.Reservation{}).
		Where("guest_id = ?", guestId).
		Find(&reservations).
		Error; err != nil {
		return 0, err
	}

	return int32(len(reservations)), nil
}

func (r *CommonRepository) GetGuestCurrentMood(ctx context.Context, guestId int32) (string, error) {
	var feedbacks []*reservationFeedbackDomain.ReservationFeedback
	result := 0

	reservationIDs, err := r.GetReservationIDsByGuestID(ctx, guestId)
	if err != nil {
		return "", err
	}

	for _, id := range reservationIDs {
		var feedback *reservationFeedbackDomain.ReservationFeedback

		if err := r.GetTenantDBConnection(ctx).
			Model(&feedback).
			Where("reservation_id = ?", id).
			First(&feedback).
			Error; err != nil {
			continue
		}

		feedbacks = append(feedbacks, feedback)
	}

	for _, feedback := range feedbacks {
		result += int(feedback.Rate)
	}

	result = result / max(1, len(feedbacks))
	if len(feedbacks) == 0 {
		return "", nil
	}

	if result < 3 {
		return "Negative", nil
	} else if result > 3 {
		return "Positive", nil
	} else {
		return "Neutral", nil
	}
}

func (r *CommonRepository) GetGuestTotalSpent(ctx context.Context, guestId int32) (float64, error) {
	var (
		orders []*reservationDomain.ReservationOrder
		result = 0.0
	)

	if err := r.GetTenantDBConnection(ctx).
		Model(&reservationDomain.ReservationOrder{}).
		Joins("JOIN reservations ON reservation_orders.reservation_id = reservations.id").
		Where("reservations.guest_id = ?", guestId).
		Find(&orders).Error; err != nil {
		return 0, err
	}

	for _, order := range orders {
		result += order.FinalTotal
	}

	return result, nil
}

func (r *CommonRepository) GetGuestTotalReservationByStatus(ctx context.Context, guestId int32, status string) (int32, error) {
	var reservations []domain2.Reservation

	if err := r.GetTenantDBConnection(ctx).
		Model(&domain2.Reservation{}).
		Joins("JOIN reservation_statuses ON reservation_statuses.id = reservations.status_id").
		Where("reservations.guest_id = ? AND reservation_statuses.name = ?", guestId, status).
		Select("reservations.*").
		Find(&reservations).Error; err != nil {
		return 0, err
	}

	return int32(len(reservations)), nil
}

func (r *CommonRepository) GetGuestBranches(ctx context.Context, guestId int32) ([]*domain.GuestBranchVisit, error) {
	var (
		branches []*domain.GuestBranchVisit
	)

	if err := r.GetTenantDBConnection(ctx).
		Table("reservations").
		Joins("JOIN guests ON guests.id = reservations.guest_id").
		Joins("JOIN branches ON branches.id = reservations.branch_id").
		Where("reservations.guest_id = ?", guestId).
		Select("branches.name AS name, count(*) AS visits").
		Group("name").
		Find(&branches).
		Error; err != nil {
		return nil, err
	}

	return branches, nil
}

func (r *CommonRepository) GetGuestTags(ctx context.Context, guestId int32) ([]*tagDomain.GuestTag, error) {
	var tags []*tagDomain.GuestTag

	if err := r.GetTenantDBConnection(ctx).
		Model(&tagDomain.GuestTag{}).
		Joins("JOIN guest_tags_assignments ON guest_tags_assignments.tag_id = guest_tags.id").
		Where("guest_tags_assignments.guest_id = ?", guestId).
		Find(&tags).Error; err != nil {
		return nil, err
	}

	for _, tag := range tags {
		var category tagDomain.GuestTagCategory

		if err := r.GetTenantDBConnection(ctx).
			Model(&tagDomain.GuestTagCategory{}).
			Where("id = ?", tag.CategoryID).
			First(&category).Error; err != nil {
			return nil, err
		}

		tag.Category = &category
	}

	return tags, nil
}

func (r *CommonRepository) GetGuestNotes(ctx context.Context, guestID int32) ([]*domain.GuestNote, error) {
	var (
		guestNotes []*domain.GuestNote
	)

	if err := r.GetTenantDBConnection(ctx).
		Model(&domain.GuestNote{}).
		Where("guest_id = ?", guestID).
		Order("created_at desc").
		Find(&guestNotes).
		Error; err != nil {
		return nil, err
	}

	for _, guestNote := range guestNotes {
		getCreator, err := r.userClient.Client.UserService.Repository.GetUserProfileByID(ctx, guestNote.CreatorID)
		if err != nil {
			guestNote.Creator = nil
		}

		guestNote.Creator = getCreator
	}

	return guestNotes, nil
}

func (r *CommonRepository) GetGuestUpcomingReservationDate(ctx context.Context, guestID int32) (string, error) {
	var dateAndTime struct {
		Date time.Time `db:"date"`
		Time string    `db:"time"`
	}

	err := r.GetTenantDBConnection(ctx).
		Model(&domain2.Reservation{}).
		Where("guest_id = ? AND date >= ?", guestID, time.Now().Format("2006-01-02")).
		Order("date asc, time asc").
		Select("date, time").
		First(&dateAndTime).Error

	if err != nil {
		return "", err
	}

	formattedDateTime := dateAndTime.Date.Format("2006-01-02") + " " + dateAndTime.Time

	return formattedDateTime, nil
}
