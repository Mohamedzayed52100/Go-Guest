package repository

import (
	"context"
	"errors"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/goplaceapp/goplace-common/pkg/utils"
	guestProto "github.com/goplaceapp/goplace-guest/api/v1"
	"github.com/goplaceapp/goplace-guest/internal/services/payment/adapters/converters"
	"github.com/goplaceapp/goplace-guest/internal/services/payment/domain"
	guestDomain "github.com/goplaceapp/goplace-guest/pkg/guestservice/domain"
	reservationDomain "github.com/goplaceapp/goplace-guest/pkg/reservationservice/domain"
	settingsProto "github.com/goplaceapp/goplace-settings/api/v1"
	itemDomain "github.com/goplaceapp/goplace-settings/pkg/itemservice/domain"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"
)

func (r *PaymentRepository) SendPaymentRequest(ctx context.Context, req *guestProto.PaymentRequest) (*guestProto.PaymentResponse, error) {
	var (
		reservation     *reservationDomain.Reservation
		currentBranchId = r.userClient.Client.UserService.Repository.GetCurrentBranchId(ctx)
	)

	if err := r.GetTenantDBConnection(ctx).First(&reservation, "id = ?", req.GetReservationId()).Error; err != nil {
		return nil, status.Error(http.StatusNotFound, "Reservation not found")
	}

	payment := &domain.PaymentRequest{
		ReservationID: req.GetReservationId(),
		BranchID:      currentBranchId,
		Delivery:      req.GetDelivery(),
		Date:          time.Now().Format(time.DateOnly),
		Invoice:       &domain.Invoice{},
	}

	branch, err := r.userClient.Client.UserService.Repository.GetBranchByID(ctx, currentBranchId)
	if err != nil {
		return nil, status.Error(http.StatusInternalServerError, err.Error())
	}

	payment.Branch = &domain.PaymentBranch{
		Name:          branch.Name,
		Address:       branch.Address,
		VatPercent:    branch.VatPercent,
		ServiceCharge: branch.ServiceCharge,
		CrNumber:      branch.CrNumber,
		VatRegNumber:  branch.VatRegNumber,
	}

	var guest *guestDomain.Guest
	if err := r.GetTenantDBConnection(ctx).
		Table("guests").
		Joins("JOIN reservations ON reservations.guest_id = guests.id").
		Where("reservations.id = ?", req.GetReservationId()).
		Select("guests.id, guests.first_name, guests.last_name, guests.phone_number, guests.email, guests.address").
		Scan(&guest).Error; err != nil {
		return nil, status.Error(http.StatusInternalServerError, err.Error())
	}

	payment.Guest = &domain.PaymentGuest{
		ID:          guest.ID,
		FirstName:   guest.FirstName,
		LastName:    guest.LastName,
		PhoneNumber: guest.PhoneNumber,
		Address:     guest.Address,
		Email:       guest.Email,
	}

	var existItems []*itemDomain.RestaurantItem
	if err := r.GetTenantDBConnection(ctx).
		Table("restaurant_items").
		Select("id, name, price").
		Scan(&existItems).Error; err != nil {
		return nil, status.Error(http.StatusInternalServerError, err.Error())
	}

	if err := r.GetTenantDBConnection(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Table("payment_requests").Create(&payment).Error; err != nil {
			return status.Error(http.StatusInternalServerError, err.Error())
		}

		for _, item := range req.GetItems() {
			foundItem, err := r.findOrCreateItem(item, existItems, ctx)
			if err != nil {
				return err
			}

			foundItem["payment_id"] = payment.ID
			if err := savePaymentItemAssignment(tx, foundItem); err != nil {
				return err
			}

			payment.Items = append(payment.Items, &domain.InvoiceItem{
				ID:       item.Id,
				Name:     item.Name,
				Price:    item.Price,
				Quantity: item.Quantity,
			})
		}

		return nil
	}); err != nil {
		return nil, err
	}

	var finalTotal float32 = 0
	for _, item := range payment.Items {
		finalTotal += item.Price * float32(item.Quantity)
	}

	payment.Invoice.SubTotal = finalTotal

	for _, contactID := range utils.ConvertStringToArrayBySeparator(req.GetContacts(), ",") {
		var contact *domain.PaymentRequestContact

		parsedContactID, err := strconv.ParseInt(contactID, 10, 32)
		if err != nil {
			return nil, status.Error(http.StatusInternalServerError, err.Error())
		}

		if err := r.GetTenantDBConnection(ctx).
			Model(&domain.PaymentRequestContact{}).
			Create(&domain.PaymentRequestContact{
				ID:               uuid.New().String(),
				PaymentRequestID: payment.ID,
				ContactID:        int32(parsedContactID),
			}).
			Scan(&contact).Error; err != nil {
			return nil, status.Error(http.StatusInternalServerError, err.Error())
		}

		payment.Contacts = append(payment.Contacts, contact)
	}

	// Get Payment gateway provider integration
	integration, err := r.commonRepository.GetIntegrationBySystemType(ctx, "Payment Gateway", reservation.BranchID)
	if err != nil {
		return nil, status.Error(http.StatusInternalServerError, "Integration not found")
	}

	switch strings.ToLower(integration.SystemName) {
	case "tap":
		tapPayment, tapRes, err := r.CreateTapPaymentInvoice(ctx, payment, integration)
		if err != nil {
			return nil, status.Error(http.StatusInternalServerError, err.Error())
		}

		var contactGuests []*domain.PaymentGuest
		for _, contact := range payment.Contacts {
			var contactGuest *domain.PaymentGuest

			if err := r.GetTenantDBConnection(ctx).
				Table("guests").
				Find(&contactGuest, "id = ?", contact.ContactID).Error; err != nil {
				if errors.Is(err, gorm.ErrRecordNotFound) || errors.Is(err, gorm.ErrDuplicatedKey) {
					continue
				}

				return nil, status.Error(http.StatusInternalServerError, err.Error())
			}

			contactGuests = append(contactGuests, contactGuest)
		}

		if isAvailableDelivery(req.Delivery, "whatsapp") {
			if err := r.SendPaymentWhatsappMessage(ctx, contactGuests, int32(payment.BranchID), tapPayment); err != nil {
				return nil, status.Error(http.StatusInternalServerError, err.Error())
			}
		}

		if isAvailableDelivery(req.Delivery, "email") {
			if err := r.SendPaymentEmail(contactGuests, tapRes["url"].(string), tapPayment.Branch.Name); err != nil {
				return nil, status.Error(http.StatusInternalServerError, err.Error())
			}
		}

		if isAvailableDelivery(req.Delivery, "sms") {
			// TODO: send sms notification
		}

		return converters.BuildPaymentResponse(tapPayment), nil
	default:
		return nil, status.Error(http.StatusInternalServerError, "Payment gateway provider not found")
	}
}

// findOrCreateItem tries to find the item in availableItems; if not found, it creates a new restaurant item.
func (r *PaymentRepository) findOrCreateItem(item *guestProto.PaymentItem, availableItems []*itemDomain.RestaurantItem, ctx context.Context) (map[string]interface{}, error) {
	for _, availableItem := range availableItems {
		if availableItem.ID == item.Id ||
			(availableItem.Name == item.Name && availableItem.Price == item.Price) {
			return map[string]interface{}{
				"price":    availableItem.Price,
				"item_id":  availableItem.ID,
				"quantity": item.Quantity,
			}, nil
		}
	}
	return r.createRestaurantItem(ctx, item)
}

// createRestaurantItem creates a new restaurant item and returns its details.
func (r *PaymentRepository) createRestaurantItem(ctx context.Context, item *guestProto.PaymentItem) (map[string]interface{}, error) {
	resItem, err := r.itemClient.Client.ItemService.CreateRestaurantItem(ctx, &settingsProto.CreateRestaurantItemRequest{
		Params: &settingsProto.RestaurantItemParams{
			Name:  item.Name,
			Price: item.Price,
			Code:  strconv.Itoa(rand.Intn(9999-1000) + 1000),
		},
	})
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"item_id":  resItem.GetResult().GetId(),
		"price":    item.Price,
		"quantity": item.Quantity,
	}, nil
}

// savePaymentItemAssignment saves the payment item assignment to the database.
func savePaymentItemAssignment(tx *gorm.DB, foundItem map[string]interface{}) error {
	return tx.Table("payment_item_assignments").Create(&foundItem).Error
}

// isAvailableDelivery checks if the given delivery method is available.
func isAvailableDelivery(list string, method string) bool {
	for _, s := range utils.ConvertStringToArrayBySeparator(list, ",") {
		if strings.EqualFold(s, method) {
			return true
		}
	}
	return false
}
