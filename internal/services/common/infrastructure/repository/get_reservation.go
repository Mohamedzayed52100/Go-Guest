package common

import (
	"context"
	"net/http"

	guestDomain "github.com/goplaceapp/goplace-guest/pkg/guestservice/domain"
	domain2 "github.com/goplaceapp/goplace-guest/pkg/reservationservice/domain"
	"google.golang.org/grpc/status"

	"github.com/goplaceapp/goplace-common/pkg/errorhelper"
	"github.com/goplaceapp/goplace-common/pkg/logger"
	"github.com/goplaceapp/goplace-guest/internal/services/reservation/domain"
	extDomain "github.com/goplaceapp/goplace-guest/pkg/reservationservice/domain"
	seatingAreaDomain "github.com/goplaceapp/goplace-settings/pkg/seatingareaservice/domain"
	tableDomain "github.com/goplaceapp/goplace-settings/pkg/tableservice/domain"
)

func (r *CommonRepository) GetReservationByID(ctx context.Context, reservationId int32) (*domain2.Reservation, error) {
	var (
		reservation *domain2.Reservation
	)

	if err := r.GetTenantDBConnection(ctx).
		Where("id = ?", reservationId).
		First(&reservation).Error; err != nil {
		return nil, status.Error(http.StatusNotFound, "Reservation not found")
	}

	reservation, err := r.GetAllReservationData(ctx, reservation)
	if err != nil {
		return nil, status.Error(http.StatusInternalServerError, err.Error())
	}

	return reservation, nil
}

func (r *CommonRepository) GetAllReservationData(ctx context.Context, res *domain2.Reservation) (*domain2.Reservation, error) {
	var (
		feedback    *domain.SimpleReservationFeedback
		seatingArea *seatingAreaDomain.SeatingArea
		note        *extDomain.ReservationNote
		totalSpent  float32
		err         error
	)

	if err := r.GetTenantDBConnection(ctx).Table("branches").
		Where("id = ?", res.BranchID).
		Select(`"id", "name"`).
		Scan(&res.Branch).Error; err != nil {
		return nil, status.Error(http.StatusNotFound, "Branch not found")
	}

	if err := r.GetTenantDBConnection(ctx).
		Table("shifts").
		Where("id = ?", res.ShiftID).
		Select(`"id", "name"`).Scan(&res.Shift).Error; err != nil {
		return nil, status.Error(http.StatusNotFound, errorhelper.ErrShiftNotFound)
	}

	tables, err := r.GetTablesForReservation(ctx, res.ID)
	if err != nil {
		res.Tables = nil
	} else {
		res.Tables = tables
		for i := range res.Tables {
			res.Tables[i].SeatingArea, err = r.seatingAreaClient.Client.SeatingAreaService.Repository.GetSeatingAreaByID(
				ctx,
				int32(res.Tables[i].SeatingAreaID),
			)
			if err != nil {
				return nil, status.Error(http.StatusInternalServerError, "Error getting seating area")
			}
		}
	}

	if res.SeatingAreaID != 0 {
		seatingArea, err = r.seatingAreaClient.Client.SeatingAreaService.Repository.GetSeatingAreaByID(ctx, res.SeatingAreaID)
		if err != nil {
			logger.Default().Errorf("Error getting seating area for reservation %d: %v", res.ID, err)
			return nil, status.Error(http.StatusInternalServerError, "Error getting seating area")
		}

		res.SeatingArea = seatingArea
	} else if len(res.Tables) > 0 {
		seatingArea, err = r.seatingAreaClient.Client.SeatingAreaService.Repository.GetSeatingAreaByID(
			ctx, int32(res.Tables[0].SeatingAreaID),
		)
		if err != nil {
			logger.Default().Errorf("Error getting seating area for reservation %d: %v", err)
			return nil, status.Error(http.StatusInternalServerError, "Error getting seating area")
		}

		res.SeatingArea = seatingArea
	} else {
		seatingArea = nil
	}

	reservationStatus, err := r.GetReservationStatusByID(ctx, res.StatusID, res.BranchID)
	if err != nil {
		logger.Default().Errorf("Error getting reservation status for reservation %d: %v", res.ID, err)
		return nil, status.Error(http.StatusInternalServerError, "Error getting reservation status")
	}
	res.Status = reservationStatus

	if res.SpecialOccasionID != nil {
		reservationSpecialOccasion, err := r.GetReservationSpecialOccasionByID(ctx, *res.SpecialOccasionID)
		if err != nil {
			res.SpecialOccasion = nil
		} else {
			res.SpecialOccasion = reservationSpecialOccasion
		}
	}

	// Fetch the primary guest
	var primaryGuest *guestDomain.Guest
	if err := r.GetTenantDBConnection(ctx).
		Table("guests").
		Where("id = ?", res.GuestID).
		First(&primaryGuest).Error; err != nil {
		return nil, status.Error(http.StatusNotFound, "Primary guest not found")
	}
	primaryGuest, err = r.GetAllGuestData(ctx, primaryGuest)
	if err != nil {
		return nil, status.Error(http.StatusInternalServerError, err.Error())
	}
	primaryGuest.IsPrimary = true
	res.Guests = []*guestDomain.Guest{primaryGuest}

	r.GetTenantDBConnection(ctx).
		Table("guests").
		Joins("JOIN reservation_visitors ON reservation_visitors.guest_id = guests.id").
		Where("reservation_id = ?", res.ID).
		Find(&res.Guests)

	res.Guests = append(res.Guests, primaryGuest)

	if err := r.GetTenantDBConnection(ctx).Model(&note).
		Where("reservation_id = ?", res.ID).
		Order("created_at desc").
		First(&note).Error; err != nil {
		res.Note = nil
	} else {
		if note.CreatorID != 0 {
			getCreator, err := r.userClient.Client.UserService.Repository.GetUserProfileByID(ctx, note.CreatorID)
			if err != nil {
				logger.Default().Errorf("Error getting creator for note %d: %v", note.ID, err)
				return nil, status.Error(http.StatusInternalServerError, "Error getting creator for note")
			}
			note.Creator = getCreator
		}

		res.Note = note
	}

	tags, err := r.GetReservationTags(ctx, res.ID)
	if err != nil {
		res.Tags = nil
	} else {
		res.Tags = tags
	}

	if res.Tables != nil {
		for _, t := range res.Tables {
			t.SeatingArea = seatingArea
		}
	}

	if err := r.GetTenantDBConnection(ctx).Model(&domain.SimpleReservationFeedback{}).
		Where("reservation_id = ?", res.ID).
		Select("id, rate, description, created_at").
		Scan(&feedback).Error; err != nil {
		res.Feedback = nil
	} else {
		res.Feedback = feedback
	}

	if res.CreatorID != 0 {
		creator, err := r.userClient.Client.UserService.Repository.GetUserProfileByID(ctx, res.CreatorID)
		if err != nil {
			logger.Default().Errorf("Error getting creator for reservation %d: %v", res.ID, err)
			return nil, status.Error(http.StatusInternalServerError, "Error getting creator for reservation")
		}

		res.Creator = creator
	} else {
		res.Creator = nil
	}

	if err := r.GetTenantDBConnection(ctx).
		Model(domain.ReservationOrder{}).
		Where("reservation_id = ?", res.ID).
		Select("COALESCE(SUM(final_total), 0)").
		Scan(&totalSpent).Error; err != nil {
		totalSpent = 0
	}
	res.TotalSpent = totalSpent

	var payStatuses []string
	r.GetTenantDBConnection(ctx).
		Table("payment_requests").
		Joins("join invoices on payment_requests.id = invoices.payment_request_id").
		Where("payment_requests.reservation_id = ?", res.ID).
		Pluck("invoices.status", &payStatuses)

	if len(payStatuses) != 0 {
		var (
			totalPaid   int32
			totalUnPaid int32
		)

		for _, status := range payStatuses {
			if status == "paid" {
				totalPaid++
			} else {
				totalUnPaid++
			}
		}

		res.Payment = &domain2.ReservationPayment{
			TotalPaid:   totalPaid,
			TotalUnPaid: totalUnPaid,
		}

		if totalPaid > 0 && totalUnPaid == 0 {
			res.Payment.Status = "paid"
		} else {
			res.Payment.Status = "unpaid"
		}
	} else {
		res.Payment = nil
	}

	return res, nil
}

func (r *CommonRepository) GetTablesForReservation(ctx context.Context, reservationId int32) ([]*tableDomain.Table, error) {
	var tables []*tableDomain.Table

	// Fetch the required data in a single query using JOIN
	query := `
		SELECT t.id, t.table_number, t.pos_number, t.min_party_size, t.max_party_size, t.seating_area_id, sa.id, sa.name
		FROM reservation_tables rt
		INNER JOIN tables t ON rt.table_id = t.id
		INNER JOIN seating_areas sa ON t.seating_area_id = sa.id
		WHERE rt.reservation_id = ?
	`

	rows, err := r.GetTenantDBConnection(ctx).Raw(query, reservationId).Rows()
	if err != nil {
		return nil, status.Error(http.StatusInternalServerError, "Error fetching tables for reservation")
	}
	defer rows.Close()

	for rows.Next() {
		var table tableDomain.Table
		var seatingArea seatingAreaDomain.SeatingArea

		if err := rows.Scan(
			&table.ID,
			&table.TableNumber,
			&table.PosNumber,
			&table.MinPartySize,
			&table.MaxPartySize,
			&table.SeatingAreaID,
			&seatingArea.ID,
			&seatingArea.Name); err != nil {
			logger.Default().Errorf("Error scanning table data: %v", err)
			return nil, status.Error(http.StatusInternalServerError, "Error scanning table data")
		}

		table.SeatingArea = &seatingArea
		tables = append(tables, &table)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return tables, nil
}
