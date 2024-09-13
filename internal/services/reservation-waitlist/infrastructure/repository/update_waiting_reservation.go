package repository

import (
	"context"
	"net/http"
	"strconv"

	guestProto "github.com/goplaceapp/goplace-guest/api/v1"
	"github.com/goplaceapp/goplace-guest/internal/services/reservation-waitlist/adapters/converters"
	reservationWaitlistDomain "github.com/goplaceapp/goplace-guest/internal/services/reservation-waitlist/domain"
	externalWaitlistDomain "github.com/goplaceapp/goplace-guest/pkg/waitlistservice/domain"

	"github.com/goplaceapp/goplace-guest/internal/services/reservation/domain"
	seatingAreaDomain "github.com/goplaceapp/goplace-settings/pkg/seatingareaservice/domain"
	shiftDomain "github.com/goplaceapp/goplace-settings/pkg/shiftservice/domain"
	"google.golang.org/grpc/status"
)

func (r *ReservationWaitListRepository) UpdateWaitingReservation(ctx context.Context, req *guestProto.UpdateWaitingReservationDetailsRequest) (*guestProto.UpdateWaitingReservationDetailsResponse, error) {
	var (
		logs        = []*reservationWaitlistDomain.ReservationWaitlistLog{}
		oldWaitList *reservationWaitlistDomain.ReservationWaitlist
		updates     = make(map[string]interface{})
		err         error
	)

	if err := r.GetTenantDBConnection(ctx).
		First(&oldWaitList, "id = ? AND branch_id = ?",
			req.GetParams().GetId(),
			r.userClient.Client.UserService.Repository.GetCurrentBranchId(ctx),
		).Error; err != nil {
		return nil, status.Error(http.StatusInternalServerError, err.Error())
	}

	oldWaitList, err = r.GetWaitingReservationData(ctx, oldWaitList.ID)
	if err != nil {
		return nil, status.Error(http.StatusInternalServerError, err.Error())
	}

	if req.GetParams().GetGuestId() != 0 && oldWaitList.GuestID != req.GetParams().GetGuestId() {
		updates["guest_id"] = req.GetParams().GetGuestId()

		logs = append(logs, &reservationWaitlistDomain.ReservationWaitlistLog{
			ReservationWaitlistID: req.GetParams().GetId(),
			FieldName:             "guest",
			Action:                "update",
			OldValue:              strconv.FormatInt(int64(oldWaitList.GuestID), 10),
			NewValue:              strconv.FormatInt(int64(req.GetParams().GetGuestId()), 10),
		})
	}

	if req.GetParams().GetGuestsNumber() != 0 && oldWaitList.GuestsNumber != req.GetParams().GetGuestsNumber() {
		updates["guests_number"] = req.GetParams().GetGuestsNumber()

		logs = append(logs, &reservationWaitlistDomain.ReservationWaitlistLog{
			ReservationWaitlistID: req.GetParams().GetId(),
			FieldName:             "guests number",
			Action:                "update",
			OldValue:              strconv.FormatInt(int64(oldWaitList.GuestsNumber), 10),
			NewValue:              strconv.FormatInt(int64(req.GetParams().GetGuestsNumber()), 10),
		})
	}

	if req.GetParams().GetSeatingAreaId() != 0 && oldWaitList.SeatingAreaID != req.GetParams().GetSeatingAreaId() {
		var newSeatingArea *seatingAreaDomain.SeatingArea
		updates["seating_area_id"] = req.GetParams().GetSeatingAreaId()

		if err := r.GetTenantDBConnection(ctx).
			First(&newSeatingArea, "id = ? AND branch_id IN ?",
				req.GetParams().GetSeatingAreaId(),
				r.userClient.Client.UserService.Repository.GetAllUserBranchesIDs(ctx),
			).Error; err != nil {
			return nil, status.Error(http.StatusInternalServerError, err.Error())
		}

		logs = append(logs, &reservationWaitlistDomain.ReservationWaitlistLog{
			ReservationWaitlistID: req.GetParams().GetId(),
			FieldName:             "seating area",
			Action:                "update",
			OldValue:              oldWaitList.SeatingArea.Name,
			NewValue:              newSeatingArea.Name,
		})
	}

	if req.GetParams().GetShiftId() != 0 && oldWaitList.ShiftID != req.GetParams().GetShiftId() {
		var newShift *shiftDomain.Shift
		updates["shift_id"] = req.GetParams().GetShiftId()

		if err := r.GetTenantDBConnection(ctx).
			First(&newShift, "id = ? AND branch_id IN ?",
				req.GetParams().GetShiftId(),
				r.userClient.Client.UserService.Repository.GetAllUserBranchesIDs(ctx),
			).Error; err != nil {
			return nil, status.Error(http.StatusInternalServerError, err.Error())
		}

		logs = append(logs, &reservationWaitlistDomain.ReservationWaitlistLog{
			ReservationWaitlistID: req.GetParams().GetId(),
			FieldName:             "shift",
			Action:                "update",
			OldValue:              oldWaitList.Shift.Name,
			NewValue:              newShift.Name,
		})
	}

	if req.GetParams().GetWaitingTime() != 0 && oldWaitList.WaitingTime != req.GetParams().GetWaitingTime() {
		updates["waiting_time"] = req.GetParams().GetWaitingTime()
		logs = append(logs, &reservationWaitlistDomain.ReservationWaitlistLog{
			ReservationWaitlistID: req.GetParams().GetId(),
			FieldName:             "waiting time",
			Action:                "update",
			OldValue:              strconv.FormatInt(int64(oldWaitList.WaitingTime), 10),
			NewValue:              strconv.FormatInt(int64(req.GetParams().GetWaitingTime()), 10),
		})
	}

	if err := r.GetTenantDBConnection(ctx).Model(&oldWaitList).Updates(updates).Error; err != nil {
		return nil, status.Error(http.StatusInternalServerError, err.Error())
	}

	if req.GetParams().GetDeleteTags() {
		var (
			newTags     = req.GetParams().GetTags()
			oldTags     = []*domain.ReservationTag{}
			deletedTags string
			createdTags string
		)

		r.GetTenantDBConnection(ctx).
			Model(&domain.ReservationTag{}).
			Joins("JOIN reservation_waitlist_tags_assignments ON reservation_waitlist_tags_assignments.tag_id = reservation_tags.id").
			Where("reservation_waitlist_tags_assignments.reservation_id = ?", oldWaitList.ID).
			Select("reservation_tags.*").
			Scan(&oldTags)

		for _, oldTag := range oldTags {
			var isFound bool
			for _, newTag := range newTags {
				if newTag.GetId() == oldTag.ID {
					isFound = true
					break
				}
			}
			if !isFound {
				r.GetTenantDBConnection(ctx).
					Delete(&externalWaitlistDomain.ReservationWaitlistTagsAssignment{}, "tag_id = ?", oldTag.ID)

				deletedTags += oldTag.Name + ","
			}
		}
		for _, newTag := range newTags {
			var isFound bool
			for _, oldTag := range oldTags {
				if newTag.GetId() == oldTag.ID {
					isFound = true
					break
				}
			}
			if !isFound {
				tagAssignment := &externalWaitlistDomain.ReservationWaitlistTagsAssignment{
					TagID:         newTag.GetId(),
					ReservationID: oldWaitList.ID,
				}

				if err := r.GetTenantDBConnection(ctx).Create(&tagAssignment).Error; err != nil {
					return nil, status.Error(http.StatusInternalServerError, err.Error())
				}

				var newTagName string
				if err := r.GetTenantDBConnection(ctx).
					Model(&domain.ReservationTag{}).
					Where("id = ?", newTag.Id).
					Select("name").
					Scan(&newTagName).
					Error; err != nil {
					return nil, status.Error(http.StatusInternalServerError, err.Error())
				}

				createdTags += newTagName + ","
			}
		}
		if deletedTags != "" {
			deletedTags = deletedTags[:len(deletedTags)-1]
			logs = append(logs, &reservationWaitlistDomain.ReservationWaitlistLog{
				ReservationWaitlistID: req.GetParams().GetId(),
				FieldName:             "tags",
				Action:                "delete",
				OldValue:              deletedTags,
			})
		}
		if createdTags != "" {
			createdTags = createdTags[:len(createdTags)-1]
			logs = append(logs, &reservationWaitlistDomain.ReservationWaitlistLog{
				ReservationWaitlistID: req.GetParams().GetId(),
				FieldName:             "tags",
				Action:                "create",
				NewValue:              createdTags,
			})
		}
	}

	waitingReservation, err := r.GetWaitingReservationData(ctx, req.GetParams().GetId())
	if err != nil {
		return nil, status.Error(http.StatusInternalServerError, err.Error())
	}

	if _, err := r.CreateReservationWaitListLogs(ctx, logs...); err != nil {
		return nil, status.Error(http.StatusInternalServerError, err.Error())
	}

	result := &guestProto.UpdateWaitingReservationDetailsResponse{
		Result: converters.BuildReservationWaitListResponse(waitingReservation),
	}

	return result, nil
}
