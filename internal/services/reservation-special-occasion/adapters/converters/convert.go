package converters

import (
	guestProto "github.com/goplaceapp/goplace-guest/api/v1"
	"github.com/goplaceapp/goplace-guest/internal/services/reservation-special-occasion/domain"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func BuildAllSpecialOccasionsResponse(specialOccasions []*domain.SpecialOccasion) []*guestProto.SpecialOccasion {
	var result []*guestProto.SpecialOccasion
	for _, s := range specialOccasions {
		result = append(result, &guestProto.SpecialOccasion{
			Id:        int32(s.ID),
			Name:      s.Name,
			Color:     s.Color,
			Icon:      s.Icon,
			CreatedAt: timestamppb.New(s.CreatedAt),
			UpdatedAt: timestamppb.New(s.UpdatedAt),
		})
	}

	return result
}
