package repository

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/goplaceapp/goplace-common/pkg/logger"
	"github.com/goplaceapp/goplace-common/pkg/meta"
	"github.com/goplaceapp/goplace-common/pkg/utils"
	"github.com/goplaceapp/goplace-guest/internal/services/payment/domain"
	"google.golang.org/grpc/status"
)

func (r *PaymentRepository) SendPaymentWhatsappMessage(ctx context.Context, guests []*domain.PaymentGuest, branchId int32, payment *domain.PaymentRequest) error {
	// Get whatsapp template
	messageTemplate, err := r.reservationRepository.GetClientWhatsappTemplate(ctx, branchId, meta.WhatsappTemplatePaymentRequest)
	if err != nil {
		logger.Default().Errorf("Failed to get whatsapp template: %v", err)
		return nil
	}

	// Get Whatsapp provider integration
	integration, err := r.commonRepository.GetIntegrationBySystemType(ctx, "WhatsApp Platform", branchId)
	if err != nil {
		logger.Default().Errorf("Failed to get whatsapp integration: %v", err)
		return nil
	}

	parseCredentials := []byte(integration.Credentials)
	var whatsappCredentials map[string]interface{}
	if err := json.Unmarshal(parseCredentials, &whatsappCredentials); err != nil {
		logger.Default().Errorf("Failed to unmarshal integration credentials: %v", err)
		return status.Error(http.StatusInternalServerError, err.Error())
	}

	switch strings.ToLower(integration.SystemName) {
	case "wati":
		receivers := make([]map[string]interface{}, 0)

		for _, guest := range guests {
			whatsappNumber := guest.PhoneNumber

			receiverPayload := map[string]interface{}{
				"customParams": []map[string]interface{}{
					{
						"name":  "name",
						"value": guest.FirstName + " " + guest.LastName,
					},
					{
						"name":  "order_number",
						"value": strconv.Itoa(int(payment.ID)),
					},
					{
						"name":  "total_amount",
						"value": strconv.Itoa(int(payment.Invoice.SubTotal)),
					},
					{
						"name":  "inv_number",
						"value": payment.Invoice.InvoiceID,
					},
				},
				"whatsappNumber": whatsappNumber,
			}

			receivers = append(receivers, receiverPayload)

			err = r.reservationRepository.UpdateContactAttributes(ctx, whatsappNumber, branchId, receiverPayload)
			if err != nil {
				logger.Default().Errorf("Failed to update contact attributes: %v", err)
				return err
			}
		}

		messagePayload := map[string]interface{}{
			"broadcast_name": "Payment Request",
			"template_name":  messageTemplate.TemplateName,
			"receivers":      receivers,
		}

		payloadBytes, err := json.Marshal(messagePayload)
		if err != nil {
			logger.Default().Errorf("Failed to marshal message payload: %v", err)
			return status.Error(http.StatusInternalServerError, err.Error())
		}

		req, err := http.NewRequest("POST", integration.BaseURL+"/sendTemplateMessages", bytes.NewBuffer(payloadBytes))
		if err != nil {
			logger.Default().Errorf("Failed to create new request: %v", err)
			return status.Error(http.StatusInternalServerError, err.Error())
		}

		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+utils.Decrypt(whatsappCredentials["api_token"].(string), os.Getenv("AES_ENCRYPTION_KEY")))

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			logger.Default().Errorf("Failed to send whatsapp message: %v", err)
			return status.Error(http.StatusInternalServerError, err.Error())
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			logger.Default().Errorf("Failed to send whatsapp message: %v", err)
			return status.Error(http.StatusInternalServerError, "Failed to send whatsapp message")
		}

		return nil

	default:
		logger.Default().Errorf("Unsupported whatsapp integration: %v", integration.SystemName)
	}

	return nil
}
