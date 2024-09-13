package repository

import (
	"context"
	"encoding/json"
	"time"

	"github.com/goplaceapp/goplace-common/pkg/logger"
	"github.com/goplaceapp/goplace-common/pkg/meta"
	guestProto "github.com/goplaceapp/goplace-guest/api/v1"
	"github.com/goplaceapp/goplace-guest/internal/services/reservation-waitlist/adapters/converters"
	"github.com/goplaceapp/goplace-guest/internal/services/reservation/infrastructure/database/listeners"
	"github.com/lib/pq"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type WaitingReservationChange struct {
	Channel string `json:"channel"`
	Payload struct {
		Operation string `json:"operation"`
		Row       struct {
			ID      int32 `json:"id"`
			ShiftID int32 `json:"shift_id"`
		}
	}
}

func (r *ReservationWaitListRepository) GetRealtimeWaitingReservations(req *guestProto.GetWaitingReservationRequest, stream guestProto.ReservationWaitlist_GetRealtimeWaitingReservationsServer) error {
	ctx := stream.Context()

	md, ok := metadata.FromIncomingContext(stream.Context())
	if !ok {
		return status.Error(codes.Internal, "error retrieving metadata")
	}

	clientDBName, authorization := md.Get(meta.TenantDBNameContextKey.String()), md.Get(meta.AuthorizationContextKey.String())
	if len(clientDBName) == 0 || len(authorization) == 0 {
		return status.Error(codes.Internal, "error retrieving client db name or authorization token")
	}

	ctx = context.WithValue(ctx, meta.TenantDBNameContextKey.String(), clientDBName[0])
	ctx = context.WithValue(ctx, meta.AuthorizationContextKey.String(), authorization[0])

	listener, cleanup := listeners.GetListener(clientDBName[0], "reservation_waitlists_change")
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

				r.processWaitingReservationChange(ctx, change, req, stream)

			case <-ctx.Done():
				return
			}
		}
	}()

	<-done
	return nil
}

func (r *ReservationWaitListRepository) processWaitingReservationChange(ctx context.Context, change *pq.Notification, req *guestProto.GetWaitingReservationRequest, stream guestProto.ReservationWaitlist_GetRealtimeWaitingReservationsServer) {
	if change.Extra == "" {
		return
	}

	var resChange *WaitingReservationChange
	if err := json.Unmarshal([]byte(change.Extra), &resChange); err != nil {
		logger.Default().Error("error unmarshalling reservation change")
		return
	}

	if resChange.Payload.Row.ID == 0 || resChange.Payload.Row.ShiftID != req.ShiftId {
		return
	}

	res, err := r.GetAllReservationWaitlistData(ctx, resChange.Payload.Row.ID)
	if err != nil {
		logger.Default().Errorf("Error retrieving reservation: %v", err)
		return
	}

	convertedDate, err := time.Parse(time.RFC3339, res.Date)
	if err != nil {
		logger.Default().Errorf("error parsing date: %v", err)
		return
	}

	if req.GetDate() == "" || convertedDate.Format(time.DateOnly) == req.GetDate() {
		if err := stream.Send(&guestProto.GetWaitingReservationResponse{
			Result: converters.BuildReservationWaitListResponse(res),
		}); err != nil {
			logger.Default().Errorf("error sending reservation: %v", err)
		}
	}
}
