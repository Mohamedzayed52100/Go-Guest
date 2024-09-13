package repository

import (
	"context"
	"encoding/json"

	"github.com/goplaceapp/goplace-common/pkg/logger"
	"github.com/goplaceapp/goplace-common/pkg/meta"
	guestProto "github.com/goplaceapp/goplace-guest/api/v1"
	"github.com/goplaceapp/goplace-guest/internal/services/reservation/infrastructure/database/listeners"
	"github.com/goplaceapp/goplace-guest/pkg/guestservice/converters"
	"github.com/lib/pq"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type reservationChange struct {
	Channel string `json:"channel"`
	Payload struct {
		Operation string `json:"operation"`
		Row       struct {
			ID        int32  `json:"id"`
			ShiftID   int32  `json:"shift_id"`
			Date      string `json:"date"`
			BranchID  int32  `json:"branch_id"`
			DeletedAt string `json:"deleted_at"`
		}
	}
}

func (r *ReservationRepository) GetRealTimeReservations(req *guestProto.GetRealtimeReservationsRequest, stream guestProto.Reservation_GetRealtimeReservationsServer) error {
	ctx := stream.Context()

	// Extract metadata in a more concise manner.
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return status.Error(codes.Internal, "error retrieving metadata")
	}

	clientDBName, authorization := md.Get(meta.TenantDBNameContextKey.String()), md.Get(meta.AuthorizationContextKey.String())
	if len(clientDBName) == 0 || len(authorization) == 0 {
		return status.Error(codes.Internal, "error retrieving client db name or authorization token")
	}

	ctx = context.WithValue(ctx, meta.TenantDBNameContextKey.String(), clientDBName[0])
	ctx = context.WithValue(ctx, meta.AuthorizationContextKey.String(), authorization[0])

	listener, cleanup := listeners.GetListener(clientDBName[0], "reservations_change")
	if listener == nil {
		return status.Error(codes.Internal, "error retrieving listener")
	}
	defer cleanup()

	done := make(chan struct{})

	go func() {
		defer close(done)

		for {
			select {
			case change, ok := <-listener.Notify:
				if !ok {
					logger.Default().Error("error listening to reservation changes")
					return
				}

				r.processChange(ctx, change, req, stream)

			case <-ctx.Done():
				return
			}
		}
	}()

	<-done
	return nil
}

// processChange processes each change notification from the database.
func (r *ReservationRepository) processChange(ctx context.Context, change *pq.Notification, req *guestProto.GetRealtimeReservationsRequest, stream guestProto.Reservation_GetRealtimeReservationsServer) {
	if change.Extra == "" {
		return
	}

	var resChange *reservationChange
	if err := json.Unmarshal([]byte(change.Extra), &resChange); err != nil {
		logger.Default().Error("error unmarshalling reservation change")
		return
	}

	// Filter out irrelevant changes.
	if resChange.Payload.Row.ID == 0 ||
		resChange.Payload.Row.ShiftID != req.ShiftId ||
		resChange.Payload.Row.Date != req.Date ||
		resChange.Payload.Row.BranchID != r.userClient.Client.UserService.Repository.GetCurrentBranchId(ctx) ||
		resChange.Payload.Row.DeletedAt != "" {
		return
	}

	res, err := r.CommonRepo.GetReservationByID(ctx, resChange.Payload.Row.ID)
	if err != nil {
		logger.Default().Errorf("Error retrieving reservation: %v", err)
		return
	}

	if err := stream.Send(&guestProto.GetReservationByIDResponse{
		Result: converters.BuildReservationProto(res),
	}); err != nil {
		logger.Default().Errorf("error sending reservation: %v", err)
	}
}
