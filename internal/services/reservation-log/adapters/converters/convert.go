package converters

import (
	"context"
	"time"

	guestProto "github.com/goplaceapp/goplace-guest/api/v1"
	common "github.com/goplaceapp/goplace-guest/internal/services/common/infrastructure/repository"
	"github.com/goplaceapp/goplace-guest/internal/services/reservation-log/domain"
	userDomain "github.com/goplaceapp/goplace-user/pkg/userservice/domain"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func BuildAllReservationLogsResponse(r *common.CommonRepository, ctx context.Context, logs []*domain.ReservationLog) []*guestProto.ReservationLog {
	result := make([]*guestProto.ReservationLog, 0)
	for _, log := range logs {
		result = append(result, BuildReservationLogResponse(r, ctx, log))
	}

	return result
}
func BuildReservationLogResponse(r *common.CommonRepository, ctx context.Context, log *domain.ReservationLog) *guestProto.ReservationLog {
	log.CreatedAt = r.ConvertToLocalTime(ctx, log.CreatedAt)
	log.CreatedAt = time.Date(log.CreatedAt.Year(), log.CreatedAt.Month(), log.CreatedAt.Day(), log.CreatedAt.Hour(), log.CreatedAt.Minute(), 0, 0, time.UTC)

	log.UpdatedAt = r.ConvertToLocalTime(ctx, log.UpdatedAt)
	log.UpdatedAt = time.Date(log.UpdatedAt.Year(), log.UpdatedAt.Month(), log.UpdatedAt.Day(), log.UpdatedAt.Hour(), log.UpdatedAt.Minute(), 0, 0, time.UTC)

	res := &guestProto.ReservationLog{
		Id:            log.ID,
		ReservationId: log.ReservationID,
		MadeBy:        log.MadeBy,
		FieldName:     log.FieldName,
		OldValue:      log.OldValue,
		NewValue:      log.NewValue,
		Action:        log.Action,
		CreatedAt:     timestamppb.New(log.CreatedAt),
		UpdatedAt:     timestamppb.New(log.UpdatedAt),
	}

	if log.Creator != nil {
		res.Creator = BuildCreatorResponse(log.Creator)
	} else {
		res.Creator = nil
	}

	return res
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

func BuildReservationWaitlistLogResponse(r *common.CommonRepository, ctx context.Context, log *domain.ReservationWaitlistLog) *guestProto.ReservationWaitlistLog {
	log.CreatedAt = r.ConvertToLocalTime(ctx, log.CreatedAt)
	log.CreatedAt = time.Date(log.CreatedAt.Year(), log.CreatedAt.Month(), log.CreatedAt.Day(), log.CreatedAt.Hour(), log.CreatedAt.Minute(), 0, 0, time.UTC)

	log.UpdatedAt = r.ConvertToLocalTime(ctx, log.UpdatedAt)
	log.UpdatedAt = time.Date(log.UpdatedAt.Year(), log.UpdatedAt.Month(), log.UpdatedAt.Day(), log.UpdatedAt.Hour(), log.UpdatedAt.Minute(), 0, 0, time.UTC)

	res := &guestProto.ReservationWaitlistLog{
		Id:                    log.ID,
		ReservationWaitlistId: log.ReservationWaitlistID,
		MadeBy:                log.MadeBy,
		FieldName:             log.FieldName,
		OldValue:              log.OldValue,
		NewValue:              log.NewValue,
		Action:                log.Action,
		CreatedAt:             timestamppb.New(log.CreatedAt),
		UpdatedAt:             timestamppb.New(log.UpdatedAt),
	}

	if log.Creator != nil {
		res.Creator = BuildCreatorResponse(log.Creator)
	} else {
		res.Creator = nil
	}

	return res
}

func BuildAllReservationWaitlistLogsResponse(r *common.CommonRepository, ctx context.Context, logs []*domain.ReservationWaitlistLog) []*guestProto.ReservationWaitlistLog {
	result := make([]*guestProto.ReservationWaitlistLog, 0)
	for _, log := range logs {
		result = append(result, BuildReservationWaitlistLogResponse(r, ctx, log))
	}

	return result
}
