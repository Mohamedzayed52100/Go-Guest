package repository

import (
	"context"
	guestProto "github.com/goplaceapp/goplace-guest/api/v1"
	"github.com/goplaceapp/goplace-guest/pkg/guestservice/converters"
	"github.com/goplaceapp/goplace-guest/pkg/reservationservice/domain"
	"google.golang.org/grpc/status"
	"net/http"
)

func (r *GuestRepository) GetGuestReservationStatistics(ctx context.Context, req *guestProto.GetGuestReservationStatisticsRequest) (*guestProto.GetGuestReservationStatisticsResponse, error) {
	var (
		result       = make(map[string]int32)
		reservations = []*domain.Reservation{}
	)

	if req.GetFromDate() != "" && req.GetToDate() != "" {
		r.GetTenantDBConnection(ctx).
			Find(&reservations,
				"guest_id = ? AND date >= ? AND date <= ?",
				req.GetGuestId(),
				req.GetFromDate(),
				req.GetToDate(),
			)
	} else if req.GetFromDate() != "" {
		r.GetTenantDBConnection(ctx).
			Find(&reservations,
				"guest_id = ? AND date >= ?",
				req.GetGuestId(),
				req.GetFromDate(),
			)
	} else if req.GetToDate() != "" {
		r.GetTenantDBConnection(ctx).
			Find(&reservations, "guest_id = ? AND date <= ?", req.GetGuestId(), req.GetToDate())
	} else {
		r.GetTenantDBConnection(ctx).Find(&reservations, "guest_id = ?", req.GetGuestId())
	}

	if req.GetQueryType() == "status" {
		for _, res := range reservations {
			// Get status by id
			getStatus, err := r.CommonRepository.GetReservationStatusByID(ctx, res.StatusID, res.BranchID)
			if err != nil {
				return nil, status.Error(http.StatusInternalServerError, err.Error())
			}

			if _, ok := result[getStatus.Name]; !ok {
				result[getStatus.Name] = 0
			}
			result[getStatus.Name]++
		}
	} else if req.GetQueryType() == "source" {
		for _, res := range reservations {
			if _, ok := result[res.ReservedVia]; !ok {
				result[res.ReservedVia] = 0
			}
			result[res.ReservedVia]++
		}
	} else if req.GetQueryType() == "qualification" {
		for _, res := range reservations {
			neutral, positive, negative := int64(0), int64(0), int64(0)

			if err := r.GetTenantDBConnection(ctx).
				Table("reservations").
				Joins("JOIN reservation_feedbacks ON reservation_feedbacks.reservation_id = reservations.id").
				Where("reservations.id = ? AND reservation_feedbacks.rate < 3", res.ID).
				Count(&negative).Error; err != nil {
				return nil, status.Error(http.StatusInternalServerError, err.Error())
			}
			if err := r.GetTenantDBConnection(ctx).
				Table("reservations").
				Joins("JOIN reservation_feedbacks ON reservation_feedbacks.reservation_id = reservations.id").
				Where("reservations.id = ? AND reservation_feedbacks.rate = 3", res.ID).
				Count(&neutral).Error; err != nil {
				return nil, status.Error(http.StatusInternalServerError, err.Error())
			}
			if err := r.GetTenantDBConnection(ctx).
				Table("reservations").
				Joins("JOIN reservation_feedbacks ON reservation_feedbacks.reservation_id = reservations.id").
				Where("reservations.id = ? AND reservation_feedbacks.rate > 3", res.ID).
				Count(&positive).Error; err != nil {
				return nil, status.Error(http.StatusInternalServerError, err.Error())
			}

			result["Neutral"] = int32(neutral)
			result["Positive"] = int32(positive)
			result["Negative"] = int32(negative)
		}
	} else if req.GetQueryType() == "party-size" {
		for _, res := range reservations {
			if res.GuestsNumber == 1 {
				if _, ok := result["Single"]; !ok {
					result["Single"] = 0
					result["Multiple"] = 0
				}
				result["Single"]++
			} else {
				if _, ok := result["Multiple"]; !ok {
					result["Single"] = 0
					result["Multiple"] = 0
				}
				result["Multiple"]++
			}
		}
	}

	return &guestProto.GetGuestReservationStatisticsResponse{
		Result: converters.BuildGuestReservationStatisticsResponse(result),
	}, nil
}
