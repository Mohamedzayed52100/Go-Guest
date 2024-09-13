package repository

import (
	"context"
	"net/http"
	"strings"
	"time"

	"github.com/goplaceapp/goplace-common/pkg/errorhelper"
	guestProto "github.com/goplaceapp/goplace-guest/api/v1"
	logDomain "github.com/goplaceapp/goplace-guest/internal/services/guest-log/domain"
	"github.com/goplaceapp/goplace-guest/pkg/guestservice/converters"
	"github.com/goplaceapp/goplace-guest/pkg/guestservice/domain"
	tagDomain "github.com/goplaceapp/goplace-settings/pkg/guesttagservice/domain"
	"google.golang.org/grpc/status"
)

/*
CreateGuest creates a new guest in the database.

Parameters:
- ctx: The context for timeout and cancellation signals.
- req: The request containing the guest details.

The method:
1. Parses the birthdate from the request.
2. Checks if the guest tags exist in the database.
3. Removes the '+' sign from the phone number.
4. Checks if the phone number and email are unique.
5. Creates the guest in the database.
6. Assigns the tags to the guest.
7. Retrieves the guest data after creation.
8. Creates a guest log for the creation of the guest record.

Returns:
- A response containing the created guest details if successful.
- An error if any operation fails.
*/

func (r *GuestRepository) CreateGuest(ctx context.Context, req *guestProto.CreateGuestRequest) (*guestProto.CreateGuestResponse, error) {
	birthdate, err := time.Parse(time.DateOnly, req.GetParams().GetBirthDate())
	if err != nil {
		birthdate = time.Now().Truncate(time.Hour)
	}

	// Check if the guest tags exist
	for _, tag := range req.GetParams().GetTags() {
		var queryTag tagDomain.GuestTag

		if err := r.GetTenantDBConnection(ctx).
			Model(&tagDomain.GuestTag{}).
			Where("id = ? AND category_id = ?", tag.Id, tag.CategoryId).
			First(&queryTag).
			Error; err != nil {
			return nil, status.Error(http.StatusBadRequest, errorhelper.ErrGuestTagNotFound)
		}
	}

	// Remove the '+' sign from the phone number
	if req.GetParams().GetPhoneNumber()[0] == '+' {
		req.GetParams().PhoneNumber = req.GetParams().PhoneNumber[1:]
	}

	req.Params.Gender = strings.ToLower(req.Params.Gender)

	// Create the guest object with the provided details
	createdGuest := &domain.Guest{
		FirstName:   req.GetParams().GetFirstName(),
		LastName:    req.GetParams().GetLastName(),
		Email:       nil,
		PhoneNumber: req.GetParams().GetPhoneNumber(),
		Language:    req.GetParams().GetLanguage(),
		Branches:    make([]*domain.GuestBranchVisit, 0),
		Tags:        make([]*tagDomain.GuestTag, 0),
		Address:     req.GetParams().GetAddress(),
		Gender:      req.GetParams().GetGender(),
	}

	if birthdate != time.Now().Truncate(time.Hour) {
		createdGuest.Birthdate = &birthdate
	} else {
		createdGuest.Birthdate = nil
	}

	// Check if the phone number is unique
	if err := r.GetTenantDBConnection(ctx).
		Model(&domain.Guest{}).
		Where("phone_number = ?", req.GetParams().GetPhoneNumber()).
		First(&domain.Guest{}).
		Error; err == nil {
		return nil, status.Error(http.StatusConflict, errorhelper.ErrDuplicateGuestPhoneNumber)
	}

	// Check if the email is unique
	if req.GetParams().GetEmail() != "" {
		if err := r.GetTenantDBConnection(ctx).
			Where("email = ?", req.GetParams().GetEmail()).
			First(&domain.Guest{}).Error; err == nil {
			return nil, status.Error(http.StatusConflict, errorhelper.ErrDuplicateGuestEmail)
		}

		email := req.GetParams().GetEmail()
		createdGuest.Email = &email
	}

	// Create the guest after all validations have passed
	if err := r.GetTenantDBConnection(ctx).
		Model(&domain.Guest{}).
		Create(&createdGuest).
		Error; err != nil {
		return nil, status.Error(http.StatusInternalServerError, err.Error())
	}

	// Assign the tags to the guest
	for _, tag := range req.GetParams().GetTags() {
		tagAssignment := &tagDomain.GuestTagsAssignment{
			TagID:   tag.GetId(),
			GuestID: createdGuest.ID,
		}

		if err := r.GetTenantDBConnection(ctx).
			Model(&tagDomain.GuestTagsAssignment{}).
			Create(&tagAssignment).
			Error; err != nil {
			return nil, status.Error(http.StatusInternalServerError, err.Error())
		}
	}

	// Get the guest data after creation to return to the client
	createdGuest, err = r.CommonRepository.GetAllGuestData(ctx, createdGuest)
	if err != nil {
		return nil, err
	}

	// Create a guest log for the creation of the guest record
	if _, err := r.CreateGuestLogs(ctx, &logDomain.GuestLog{
		GuestID:   createdGuest.ID,
		Action:    "create",
		FieldName: "guest",
	}); err != nil {
		return nil, status.Error(http.StatusInternalServerError, err.Error())
	}

	return &guestProto.CreateGuestResponse{
		Result: converters.BuildGuestResponse(createdGuest),
	}, nil

}
