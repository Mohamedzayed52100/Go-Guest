package converters

import (
	"fmt"
	"strconv"

	"github.com/goplaceapp/goplace-common/pkg/utils"
	guestProto "github.com/goplaceapp/goplace-guest/api/v1"
	"github.com/goplaceapp/goplace-guest/internal/services/payment/domain"
)

func BuildAllPaymentsResponse(payments []*domain.PaymentRequest) []*guestProto.PaymentResponse {
	var result []*guestProto.PaymentResponse
	for _, payment := range payments {
		result = append(result, BuildPaymentResponse(payment))
	}

	return result
}

func BuildPaymentResponse(payment *domain.PaymentRequest) *guestProto.PaymentResponse {
	res := &guestProto.PaymentResponse{
		Id:     payment.ID,
		Status: payment.Invoice.Status,
		Guest: &guestProto.PaymentGuest{
			FirstName:   payment.Guest.FirstName,
			LastName:    payment.Guest.LastName,
			PhoneNumber: payment.Guest.PhoneNumber,
			Address:     payment.Guest.Address,
		},
		Delivery: utils.ConvertStringToArrayBySeparator(payment.Delivery, ","),
		Invoice: &guestProto.Invoice{
			InvoiceId:  payment.Invoice.InvoiceID,
			InvoiceRef: fmt.Sprintf("INV_%06d", payment.Invoice.ID),
			Date:       payment.Date,
			Waiter:     "",
			Items:      buildPaymentItemsResponse(payment.Items),
		},
		Branch: &guestProto.PaymentBranch{
			Name:          payment.Branch.Name,
			Address:       payment.Branch.Address,
			VatPercent:    payment.Branch.VatPercent,
			ServiceCharge: payment.Branch.ServiceCharge,
			CrNumber:      payment.Branch.CrNumber,
			VatRegNumber:  payment.Branch.VatRegNumber,
		},
		Contacts: int32(len(payment.Contacts)),
	}

	if payment.Invoice.Status == "paid" {
		res.Card = &guestProto.PaymentCard{
			CardType:       payment.Invoice.CardType,
			LastFourDigits: payment.Invoice.LastFourDigits,
			CardExpireDate: payment.Invoice.ExpDate,
		}
	}

	for _, item := range payment.Items {
		res.Invoice.SubTotal += item.Price * float32(item.Quantity)
	}

	total, _ := strconv.ParseFloat(fmt.Sprintf("%.2f", res.Invoice.SubTotal), 32)
	res.Invoice.SubTotal = float32(total)

	return res
}

func buildPaymentItemsResponse(items []*domain.InvoiceItem) []*guestProto.PaymentItem {
	var result []*guestProto.PaymentItem
	for _, item := range items {
		result = append(result, &guestProto.PaymentItem{
			Id:       item.ID,
			Name:     item.Name,
			Quantity: item.Quantity,
			Price:    item.Price,
		})
	}
	return result
}
