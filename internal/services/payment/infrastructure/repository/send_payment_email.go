package repository

import (
	"os"

	"github.com/goplaceapp/goplace-guest/internal/services/payment/domain"
	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

func (r *PaymentRepository) SendPaymentEmail(guests []*domain.PaymentGuest, paymentLink, branchName string) error {
	for _, guest := range guests {
		from := mail.NewEmail("GoPlace", os.Getenv("SENDGRID_NO_REPLY_EMAIL"))
		subject := "Payment Request"
		to := mail.NewEmail(guest.FirstName+" "+guest.LastName, *guest.Email)

		plainTextContent := "Dear " + guest.FirstName + " " + guest.LastName + ",\n\n" +
			"Thank you for choosing " + branchName + " for your recent stay.\n\n" +
			"We hope you had a pleasant experience.\nWe kindly request you to complete the payment for your stay using the secure link below: " +
			paymentLink + "\n\n" +
			"If you have any questions or need assistance, please contact us.\nThank you for your prompt attention to this matter."

		htmlContent := "Dear " + guest.FirstName + " " + guest.LastName + ",<br><br>" +
			"Thank you for choosing " + branchName + " for your recent stay.<br><br>" +
			"We hope you had a pleasant experience.<br>We kindly request you to complete the payment for your stay using the secure link below: " + "<a href = " +
			paymentLink + ">Payment Link</a><br><br>" +
			"If you have any questions or need assistance, please contact us.<br>Thank you for your prompt attention to this matter."

		message := mail.NewSingleEmail(from, subject, to, plainTextContent, htmlContent)

		client := sendgrid.NewSendClient(os.Getenv("SENDGRID_API_KEY"))
		_, err := client.Send(message)
		if err != nil {
			return err
		}
	}
	return nil
}
