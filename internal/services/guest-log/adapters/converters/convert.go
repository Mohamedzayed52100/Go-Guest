package converters

import (
	guestProto "github.com/goplaceapp/goplace-guest/api/v1"
	"github.com/goplaceapp/goplace-guest/internal/services/guest-log/domain"
	userDomain "github.com/goplaceapp/goplace-user/pkg/userservice/domain"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func BuildGuestLogResponse(log *domain.GuestLog) *guestProto.GuestLog {
	res := &guestProto.GuestLog{
		Id:        log.ID,
		GuestId:   log.GuestID,
		MadeBy:    log.MadeBy,
		FieldName: log.FieldName,
		OldValue:  log.OldValue,
		NewValue:  log.NewValue,
		Action:    log.Action,
		CreatedAt: timestamppb.New(log.CreatedAt),
		UpdatedAt: timestamppb.New(log.UpdatedAt),
	}

	if log.Creator != nil {
		res.Creator = BuildCreatorResponse(log.Creator)
	} else {
		res.Creator = nil
	}

	return res
}

func BuildAllGuestLogsResponse(logs []*domain.GuestLog) []*guestProto.GuestLog {
	result := make([]*guestProto.GuestLog, 0)
	for _, log := range logs {
		result = append(result, BuildGuestLogResponse(log))
	}

	return result
}
func BuildCreatorResponse(creator *userDomain.User) *guestProto.CreatorProfile {
	return &guestProto.CreatorProfile{
		Id:          creator.ID,
		FirstName:   creator.FirstName,
		LastName:    creator.LastName,
		Email:       creator.Email,
		PhoneNumber: creator.PhoneNumber,
		Avatar:      creator.Avatar,
		Role:        creator.Role.DisplayName,
	}
}
