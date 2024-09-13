package repository

import (
	"context"
	"strconv"
	"time"

	domain2 "github.com/goplaceapp/goplace-guest/pkg/reservationservice/domain"

	"github.com/goplaceapp/goplace-common/pkg/logger"
	"github.com/goplaceapp/goplace-common/pkg/utils"
	tableDomain "github.com/goplaceapp/goplace-settings/pkg/tableservice/domain"
)

func (r *ReservationRepository) CreateReservationTables(ctx context.Context, date time.Time, startTime, endTime time.Time, guestsNumber, branchId int) ([]*tableDomain.Table, error) {
	var reservations []*domain2.Reservation
	var tables []*tableDomain.Table

	db := r.GetTenantDBConnection(ctx)

	// Fetch reservations
	if err := db.Table("reservations").
		Joins("JOIN turnover ON least(10, reservations.guests_number) = turnover.guests_number").
		Where("(reservations.date = ? AND reservations.time < ? AND (reservations.time + (turnover.turnover_time * interval '1 minute')) > ? AND reservations.branch_id = ?)",
			date.Format("2006-01-02"),
			endTime.Format("15:04:05"),
			startTime.Format("15:04:05"),
			branchId).
		Distinct().
		Find(&reservations).Error; err != nil {
		return nil, err
	}

	// Fetch tables
	if err := db.Model(&tableDomain.Table{}).
		Joins("LEFT JOIN table_statuses ON table_statuses.table_id = tables.id").
		Where("tables.branch_id = ? AND (table_statuses.date IS NULL OR table_statuses.date <> ?)",
			branchId,
			date.Format("2006-01-02")).
		Distinct().
		Find(&tables).Error; err != nil {
		return nil, err
	}

	availableTables := []*tableDomain.Table{}
	// Check each table for availability
	for _, table := range tables {
		logger.Default().Infof("Table: %v", table)

		if table.CombinedTables != nil && *table.CombinedTables != "" {
			continue
		}

		isReserved := false
		for _, res := range reservations {
			if err := db.
				First(&domain2.ReservationTable{}, "table_id = ? AND reservation_id = ?", table.ID, res.ID).
				Error; err == nil {
				isReserved = true
				break
			}
		}

		if !isReserved && (table.MinPartySize <= guestsNumber && table.MaxPartySize >= guestsNumber) {
			availableTables = append(availableTables, table)
		}
	}

	// Check combined tables for availability
	for _, table := range tables {
		if table.CombinedTables != nil && *table.CombinedTables != "" && table.MinPartySize <= guestsNumber && table.MaxPartySize >= guestsNumber {
			isReserved := false
			combinedTables := utils.ConvertStringToArrayBySeparator(*table.CombinedTables, ",")

			for _, ct := range combinedTables {
				tableId, _ := strconv.ParseInt(ct, 10, 32)

				for _, res := range reservations {
					if err := db.First(&domain2.ReservationTable{}, "table_id = ? AND reservation_id = ?", tableId, res.ID).Error; err == nil {
						isReserved = true
						break
					}
				}
				if isReserved {
					break
				}
			}

			if !isReserved {
				availableTables = append(availableTables, table)
			}
		}
	}

	return availableTables, nil
}
