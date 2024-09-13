package repository

import (
	"context"
	"net/http"
	"strings"
	"time"

	"github.com/goplaceapp/goplace-common/pkg/utils"

	"github.com/goplaceapp/goplace-common/pkg/errorhelper"
	"github.com/goplaceapp/goplace-common/pkg/rbac"
	guestProto "github.com/goplaceapp/goplace-guest/api/v1"
	logDomain "github.com/goplaceapp/goplace-guest/internal/services/guest-log/domain"
	"github.com/goplaceapp/goplace-guest/pkg/guestservice/converters"
	"github.com/goplaceapp/goplace-guest/pkg/guestservice/domain"
	tagDomain "github.com/goplaceapp/goplace-settings/pkg/guesttagservice/domain"
	"google.golang.org/grpc/status"
)

func (r *GuestRepository) UpdateGuest(ctx context.Context, req *guestProto.UpdateGuestRequest) (*guestProto.UpdateGuestResponse, error) {
	permissions := r.RoleClient.Client.RoleService.Repository.GetAllStringPermissions(ctx)

	currentGuest, err := r.CommonRepository.GetGuestByID(ctx, &guestProto.GetGuestByIDRequest{
		Id: req.GetParams().GetId(),
	})
	if err != nil {
		return nil, status.Error(http.StatusNotFound, errorhelper.ErrGuestNotFound)
	}

	updates := make(map[string]interface{})
	logs := []*logDomain.GuestLog{}

	birthdate := ""
	if currentGuest.GetResult().GetBirthDate() != nil {
		birthdate = currentGuest.GetResult().GetBirthDate().AsTime().Format(time.DateOnly)
	}

	if req.GetParams().GetAddress() != "" && req.GetParams().GetAddress() != currentGuest.Result.Address {
		updates["address"] = req.GetParams().GetAddress()

		logs = append(logs, &logDomain.GuestLog{
			GuestID:   req.GetParams().GetId(),
			FieldName: "address",
			Action:    "update",
			OldValue:  currentGuest.Result.Address,
			NewValue:  req.GetParams().GetAddress(),
		})
	}

	if req.GetParams().GetBirthDate() != "" && !req.GetParams().GetEmptyBirthdate() && birthdate != req.GetParams().GetBirthDate() {
		convertedBirthdate, err := time.Parse(time.DateOnly, req.GetParams().GetBirthDate())
		if err != nil {
			return nil, status.Error(http.StatusInternalServerError, err.Error())
		}
		updates["birthdate"] = convertedBirthdate

		logs = append(logs, &logDomain.GuestLog{
			GuestID:   req.GetParams().GetId(),
			FieldName: "birthdate",
			Action:    "update",
			OldValue:  birthdate,
			NewValue:  convertedBirthdate.Format(time.DateOnly),
		})
	} else if req.GetParams().GetBirthDate() == "" && !req.GetParams().GetEmptyBirthdate() && currentGuest.GetResult().GetBirthDate() != nil {
		updates["birthdate"] = nil
		logs = append(logs, &logDomain.GuestLog{
			GuestID:   req.GetParams().GetId(),
			FieldName: "birthdate",
			Action:    "delete",
			OldValue:  birthdate,
			NewValue:  "",
		})
	}

	if req.GetParams().GetFirstName() != "" &&
		currentGuest.GetResult().GetFirstName() != req.GetParams().GetFirstName() &&
		utils.HasPermission(permissions, rbac.ChangeGuestName.Name) {
		updates["first_name"] = req.GetParams().GetFirstName()
		logs = append(logs, &logDomain.GuestLog{
			GuestID:   req.GetParams().GetId(),
			FieldName: "first name",
			OldValue:  currentGuest.GetResult().GetFirstName(),
			NewValue:  updates["first_name"].(string),
		})
	} else if req.GetParams().GetFirstName() != "" &&
		!utils.HasPermission(permissions, rbac.ChangeGuestName.Name) {
		return nil, status.Error(http.StatusForbidden, "You don't have permission to change the guest's name")
	}

	if req.GetParams().GetLastName() != "" &&
		currentGuest.GetResult().GetLastName() != req.GetParams().GetLastName() &&
		utils.HasPermission(permissions, rbac.ChangeGuestName.Name) {
		updates["last_name"] = req.GetParams().GetLastName()
		logs = append(logs, &logDomain.GuestLog{
			GuestID:   req.GetParams().GetId(),
			FieldName: "last name",
			OldValue:  currentGuest.GetResult().GetLastName(),
			NewValue:  updates["last_name"].(string),
		})
	} else if req.GetParams().GetLastName() != "" &&
		!utils.HasPermission(permissions, rbac.ChangeGuestName.Name) {
		return nil, status.Error(http.StatusForbidden, "You don't have permission to change the guest's name")
	}

	if req.GetParams().GetPhoneNumber() != "" &&
		currentGuest.GetResult().GetPhoneNumber() != req.GetParams().GetPhoneNumber() &&
		utils.HasPermission(permissions, rbac.ChangeGuestPhone.Name) {
		updates["phone_number"] = req.GetParams().GetPhoneNumber()
		logs = append(logs, &logDomain.GuestLog{
			GuestID:   req.GetParams().GetId(),
			FieldName: "phone number",
			OldValue:  currentGuest.GetResult().GetPhoneNumber(),
			NewValue:  updates["phone_number"].(string),
		})
	} else if req.GetParams().GetPhoneNumber() != "" &&
		!utils.HasPermission(permissions, rbac.ChangeGuestPhone.Name) {
		return nil, status.Error(http.StatusForbidden, "You don't have permission to change the guest's phone number")
	}

	req.Params.Gender = strings.ToLower(req.Params.Gender)
	if req.GetParams().GetGender() != "" &&
		currentGuest.GetResult().GetGender() != req.GetParams().GetGender() {
		updates["gender"] = req.GetParams().GetGender()
		logs = append(logs, &logDomain.GuestLog{
			GuestID:   req.GetParams().GetId(),
			FieldName: "gender",
			Action:    "update",
			OldValue:  currentGuest.GetResult().GetGender(),
			NewValue:  updates["gender"].(string),
		})
	}

	if (req.GetParams().GetLanguage() != "" &&
		currentGuest.GetResult().GetLanguage() != req.GetParams().GetLanguage()) ||
		!req.GetParams().GetEmptyLanguage() {

		updates["language"] = req.GetParams().GetLanguage()

		if currentGuest.GetResult().GetLanguage() != req.GetParams().GetLanguage() {
			logs = append(logs, &logDomain.GuestLog{
				GuestID:   req.GetParams().GetId(),
				FieldName: "language",
				Action:    "update",
				OldValue:  currentGuest.GetResult().GetLanguage(),
				NewValue:  updates["language"].(string),
			})
		}
	}

	if (req.GetParams().GetEmail() != "" && currentGuest.GetResult().GetEmail() != req.GetParams().GetEmail()) ||
		!req.GetParams().GetEmptyEmail() {
		updates["email"] = req.GetParams().GetEmail()
		if currentGuest.GetResult().GetEmail() == "" || currentGuest.GetResult().GetEmail() != req.GetParams().GetEmail() {
			logs = append(logs, &logDomain.GuestLog{
				GuestID:   req.GetParams().GetId(),
				FieldName: "email",
				Action:    "update",
				OldValue:  currentGuest.GetResult().GetEmail(),
				NewValue:  updates["email"].(string),
			})
		}
	}

	if err := r.GetTenantDBConnection(ctx).
		Model(&domain.Guest{}).
		Where("id = ?", req.GetParams().GetId()).
		Updates(updates).Error; err != nil {
		return nil, status.Error(http.StatusInternalServerError, err.Error())
	}

	if !req.GetParams().GetEmptyTags() &&
		utils.HasPermission(permissions, rbac.EditGuestTags.Name) {
		var newTags []int32
		for _, tag := range req.GetParams().GetTags() {
			newTags = append(newTags, tag.Id)
		}
		var existingTags []int32
		if err := r.GetTenantDBConnection(ctx).
			Model(&tagDomain.GuestTagsAssignment{}).
			Where("guest_id = ?", req.GetParams().GetId()).
			Distinct("tag_id").
			Pluck("tag_id", &existingTags).
			Error; err != nil {
			return nil, status.Error(http.StatusInternalServerError, err.Error())
		}

		oldTagsMap := make(map[int32]bool)
		for _, tag := range existingTags {
			oldTagsMap[tag] = true
		}

		newTagsMap := make(map[int32]bool)
		for _, tag := range newTags {
			newTagsMap[tag] = true
		}

		var deletedTags []string
		for tag := range oldTagsMap {
			if _, ok := newTagsMap[tag]; !ok {
				var fullTag *tagDomain.GuestTag

				if err := r.GetTenantDBConnection(ctx).
					Model(&tagDomain.GuestTag{}).
					Where("id =?", tag).
					Find(&fullTag).
					Error; err != nil {
					return nil, status.Error(http.StatusInternalServerError, err.Error())
				}

				if err := r.GetTenantDBConnection(ctx).
					Model(&tagDomain.GuestTagsAssignment{}).
					Delete(tag, "tag_id = ?", tag).
					Error; err != nil {
					return nil, status.Error(http.StatusInternalServerError, err.Error())
				}

				deletedTags = append(deletedTags, fullTag.Name)
			}
		}

		var createdTags []string
		for tag := range newTagsMap {
			if _, ok := oldTagsMap[tag]; !ok {
				var fullTag *tagDomain.GuestTag

				if err := r.GetTenantDBConnection(ctx).
					Model(&tagDomain.GuestTagsAssignment{}).
					Create(&tagDomain.GuestTagsAssignment{
						GuestID: req.GetParams().GetId(),
						TagID:   tag,
					}).Error; err != nil {
					return nil, status.Error(http.StatusInternalServerError, err.Error())
				}

				if err := r.GetTenantDBConnection(ctx).
					Model(&tagDomain.GuestTag{}).
					Where("id =?", tag).
					Find(&fullTag).
					Error; err != nil {
					return nil, status.Error(http.StatusInternalServerError, err.Error())
				}

				createdTags = append(createdTags, fullTag.Name)
			}
		}

		if len(deletedTags) > 0 {
			var parsedDeletedTags string

			for i, tag := range deletedTags {
				if i == 0 {
					parsedDeletedTags += tag
				} else {
					parsedDeletedTags += ", " + tag
				}
			}

			logs = append(logs, &logDomain.GuestLog{
				GuestID:   req.GetParams().GetId(),
				FieldName: "tags",
				Action:    "delete",
				OldValue:  parsedDeletedTags,
			})
		}

		if len(createdTags) > 0 {
			var parsedCreatedTags string

			for i, tag := range createdTags {
				if i == 0 {
					parsedCreatedTags += tag
				} else {
					parsedCreatedTags += ", " + tag
				}
			}

			logs = append(logs, &logDomain.GuestLog{
				GuestID:   req.GetParams().GetId(),
				FieldName: "tags",
				Action:    "create",
				NewValue:  parsedCreatedTags,
			})
		}
	} else if !req.GetParams().GetEmptyTags() && !utils.HasPermission(permissions, rbac.EditGuestTags.Name) {
		return nil, status.Error(http.StatusForbidden, "You don't have permission to change the guest's tags")
	}

	for _, rlog := range logs {
		if rlog.Action == "" {
			rlog.Action = "update"
		}
	}

	if _, err := r.CreateGuestLogs(ctx, logs...); err != nil {
		return nil, status.Error(http.StatusInternalServerError, err.Error())
	}

	getUpdatedGuest, err := r.CommonRepository.GetAllGuestData(ctx, &domain.Guest{ID: currentGuest.Result.Id})
	if err != nil {
		return nil, status.Error(http.StatusNotFound, errorhelper.ErrGuestNotFound)
	}

	return &guestProto.UpdateGuestResponse{
		Result: converters.BuildGuestResponse(getUpdatedGuest),
	}, nil

}
