package repository

import (
	"context"
	"net/http"
	"strconv"
	"time"

	domain2 "github.com/goplaceapp/goplace-guest/pkg/reservationservice/domain"
	tableDomain "github.com/goplaceapp/goplace-settings/pkg/tableservice/domain"

	"github.com/goplaceapp/goplace-common/pkg/errorhelper"
	guestProto "github.com/goplaceapp/goplace-guest/api/v1"
	reservationLogDomain "github.com/goplaceapp/goplace-guest/internal/services/reservation-log/domain"
	"github.com/goplaceapp/goplace-guest/pkg/guestservice/converters"
	"google.golang.org/grpc/status"
)

func (r *ReservationRepository) UpdateReservationTable(ctx context.Context, req *guestProto.UpdateReservationTableRequest) (*guestProto.UpdateReservationTableResponse, error) {
	var (
		res              *domain2.Reservation
		logs             = []*reservationLogDomain.ReservationLog{}
		newSeatingAreaId int32
	)

	if len(req.GetTables()) == 0 {
		return nil, status.Error(http.StatusBadRequest, errorhelper.ErrTableIDIsRequired)
	}

	if err := r.GetTenantDBConnection(ctx).
		First(&res, "id = ?", req.GetReservationId()).
		Error; err != nil {
		return nil, status.Error(http.StatusInternalServerError, errorhelper.ErrReservationNotFound)
	}

	loggedInUser, err := r.userClient.Client.UserService.Repository.GetLoggedInUser(ctx)
	if err != nil {
		return nil, status.Error(http.StatusInternalServerError, err.Error())
	}

	oldTables, err := r.CommonRepo.GetTablesForReservation(ctx, req.GetReservationId())
	if err != nil {
		return nil, status.Error(http.StatusInternalServerError, err.Error())
	}

	for _, i := range oldTables {
		if err := r.GetTenantDBConnection(ctx).
			Delete(&domain2.ReservationTable{}, "reservation_id = ? AND table_id = ?",
				req.GetReservationId(),
				i.ID).
			Error; err != nil {
			return nil, status.Error(http.StatusInternalServerError, err.Error())
		}
	}

	for _, i := range req.GetTables() {
		var table *tableDomain.Table

		if err := r.GetTenantDBConnection(ctx).
			First(&table, "id = ?", i).
			Error; err != nil {
			return nil, status.Error(http.StatusInternalServerError, "Table with id "+strconv.Itoa(int(i))+" not found")
		}

		newSeatingAreaId = int32(table.SeatingAreaID)

		if err := r.GetTenantDBConnection(ctx).
			Create(&domain2.ReservationTable{
				ReservationID: int(req.GetReservationId()),
				TableID:       int(i),
			}).
			Scan(&domain2.ReservationTable{}).Error; err != nil {
			return nil, status.Error(http.StatusInternalServerError, err.Error())
		}

		existsInOld := false
		for _, oldTable := range oldTables {
			if int32(oldTable.ID) == i {
				existsInOld = true
				break
			}
		}

		if !existsInOld {
			logs = append(logs, &reservationLogDomain.ReservationLog{
				ReservationID: req.GetReservationId(),
				CreatorID:     loggedInUser.ID,
				FieldName:     "tables",
				OldValue:      "",
				NewValue:      table.TableNumber,
				Action:        "update",
			})
		}
	}

	if len(oldTables) > 0 {
		for _, oldTable := range oldTables {
			existsInNew := false
			for _, newTable := range req.GetTables() {
				if int32(oldTable.ID) == newTable {
					existsInNew = true
					break
				}
			}

			if !existsInNew {
				logs = append(logs, &reservationLogDomain.ReservationLog{
					ReservationID: req.GetReservationId(),
					CreatorID:     loggedInUser.ID,
					FieldName:     "tables",
					OldValue:      strconv.Itoa(oldTable.ID),
					NewValue:      "",
					Action:        "delete",
				})
			}
		}
	}

	if len(logs) > 0 {
		for _, log := range logs {
			if _, err := r.CreateReservationLogs(ctx, &reservationLogDomain.ReservationLog{
				ReservationID: req.GetReservationId(),
				CreatorID:     log.CreatorID,
				MadeBy:        log.MadeBy,
				FieldName:     log.FieldName,
				OldValue:      log.OldValue,
				NewValue:      log.NewValue,
				Action:        log.Action,
				CreatedAt:     log.CreatedAt,
				UpdatedAt:     log.UpdatedAt,
			}); err != nil {
				return nil, status.Error(http.StatusInternalServerError, err.Error())
			}
		}
	}

	if err := r.GetTenantDBConnection(ctx).
		Model(&res).
		Update("updated_at", time.Now()).
		Error; err != nil {
		return nil, status.Error(http.StatusInternalServerError, err.Error())
	}

	if res.SeatingAreaID != newSeatingAreaId {
		if err := r.GetTenantDBConnection(ctx).
			Model(&res).
			Update("seating_area_id", newSeatingAreaId).
			Error; err != nil {
			return nil, status.Error(http.StatusInternalServerError, err.Error())
		}
	}

	res, err = r.CommonRepo.GetAllReservationData(ctx, res)
	if err != nil {
		return nil, status.Error(http.StatusInternalServerError, err.Error())
	}

	return &guestProto.UpdateReservationTableResponse{
		Result: converters.BuildReservationProto(res),
	}, nil
}
