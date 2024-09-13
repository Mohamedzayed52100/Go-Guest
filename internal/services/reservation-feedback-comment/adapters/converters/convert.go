package converters

import (
	guestProto "github.com/goplaceapp/goplace-guest/api/v1"
	"github.com/goplaceapp/goplace-guest/internal/services/reservation-feedback-comment/domain"
	feedbackDomain "github.com/goplaceapp/goplace-guest/internal/services/reservation-feedback/domain"
	userDomain "github.com/goplaceapp/goplace-user/pkg/userservice/domain"
	"google.golang.org/protobuf/types/known/timestamppb"
)

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

func BuildReservationFeedbackCommentResponse(comment *domain.ReservationFeedbackComment) *guestProto.ReservationFeedbackComment {
	res := &guestProto.ReservationFeedbackComment{
		Id:        int32(comment.ID),
		Creator:   BuildCreatorResponse(comment.Creator),
		Comment:   comment.Comment,
		CreatedAt: timestamppb.New(comment.CreatedAt),
		UpdatedAt: timestamppb.New(comment.UpdatedAt),
	}

	return res
}

func BuildAllReservationFeedbackCommentsResponse(comments []*domain.ReservationFeedbackComment) []*guestProto.ReservationFeedbackComment {
	response := []*guestProto.ReservationFeedbackComment{}
	for _, comment := range comments {
		response = append(response, BuildReservationFeedbackCommentResponse(comment))
	}

	return response
}

func BuildReservationFeedbackSolutionResponse(solution *feedbackDomain.ReservationFeedbackSolution) *guestProto.ReservationFeedbackSolution {
	return &guestProto.ReservationFeedbackSolution{
		Id:        int32(solution.ID),
		Creator:   BuildCreatorResponse(solution.Creator),
		Solution:  solution.Solution,
		CreatedAt: timestamppb.New(solution.CreatedAt),
		UpdatedAt: timestamppb.New(solution.UpdatedAt),
	}
}
