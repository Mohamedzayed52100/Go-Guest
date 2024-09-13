package repository

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/goplaceapp/goplace-common/pkg/logger"
	"github.com/goplaceapp/goplace-common/pkg/meta"
	"github.com/goplaceapp/goplace-common/pkg/utils"
	guestDomain "github.com/goplaceapp/goplace-guest/pkg/guestservice/domain"
	"github.com/goplaceapp/goplace-guest/pkg/reservationservice/domain"
	"google.golang.org/grpc/status"
)

func (r *ReservationRepository) SendBulkReservationWhatsappFeedback(ctx context.Context, reservationIds []int32, branchId int32) (bool, error) {
	var (
		sent bool
		err  error
	)

	receivers := make([]map[string]interface{}, 0)
	reservations := make([]*domain.Reservation, 0)

	// Get reservation details
	for _, reservationId := range reservationIds {
		getReservation, err := r.CommonRepo.GetReservationByID(ctx, reservationId)
		if err != nil {
			logger.Default().Errorf("Failed to get reservation by id: %v", err)
			return false, err
		}

		var primaryGuest *guestDomain.Guest
		r.GetTenantDBConnection(ctx).First(&primaryGuest, "id = ?", getReservation.GuestID)

		// Update contact attributes with new reservation details
		receiverPayload := map[string]interface{}{
			"customParams": []map[string]interface{}{
				{
					"name":  "name",
					"value": primaryGuest.FirstName + " " + primaryGuest.LastName,
				},
				{
					"name":  "reservation_number",
					"value": strconv.Itoa(int(getReservation.ID)),
				},
				{
					"name":  "date",
					"value": getReservation.Date.Format("02 Jan 2006"),
				},
				{
					"name":  "time",
					"value": getReservation.Time,
				},
				{
					"name":  "restaurant_branch",
					"value": getReservation.Branch.Name,
				},
				{
					"name":  "person_count",
					"value": strconv.Itoa(int(getReservation.GuestsNumber)),
				},
			},
			"whatsappNumber": primaryGuest.PhoneNumber,
		}

		receivers = append(receivers, receiverPayload)
		reservations = append(reservations, getReservation)
	}

	// Get whatsapp template
	messageTemplate, err := r.GetClientWhatsappTemplate(ctx, branchId, meta.WhatsappTemplateTypeFeedbackReservation)
	if err != nil {
		logger.Default().Errorf("Failed to get whatsapp template: %v", err)
		return false, nil
	}

	// Get Whatsapp provider integration
	integration, err := r.CommonRepo.GetIntegrationBySystemType(ctx, "WhatsApp Platform", branchId)
	if err != nil {
		logger.Default().Errorf("Failed to get whatsapp integration: %v", err)
		return false, nil
	}

	parseCredentials := []byte(integration.Credentials)
	var revelIntegrationCredentials map[string]interface{}
	if err := json.Unmarshal(parseCredentials, &revelIntegrationCredentials); err != nil {
		logger.Default().Errorf("Failed to unmarshal integration credentials: %v", err)
		return false, status.Error(http.StatusInternalServerError, err.Error())
	}

	for _, res := range reservations {
		reservationTime, err := time.Parse("15:04:05", res.Time)
		if err != nil {
			logger.Default().Errorf("Failed to parse reservation time: %v", err)
			return false, status.Error(http.StatusInternalServerError, err.Error())
		}

		var primaryGuest *guestDomain.Guest
		r.GetTenantDBConnection(ctx).First(&primaryGuest, "id = ?", res.GuestID)

		// Update contact attributes with new reservation details
		updateContactPayload := map[string]interface{}{
			"customParams": []map[string]interface{}{
				{
					"name":  "name",
					"value": primaryGuest.FirstName + " " + primaryGuest.LastName,
				},
				{
					"name":  "reservation_number",
					"value": strconv.Itoa(int(res.ID)),
				},
				{
					"name":  "date",
					"value": res.Date.Format("02 Jan 2006"),
				},
				{
					"name":  "time",
					"value": reservationTime.Format("03:04 PM"),
				},
				{
					"name":  "restaurant_branch",
					"value": res.Branch.Name,
				},
				{
					"name":  "person_count",
					"value": strconv.Itoa(int(res.GuestsNumber)),
				},
			},
		}

		err = r.UpdateContactAttributes(ctx, primaryGuest.PhoneNumber, res.BranchID, updateContactPayload)
		if err != nil {
			logger.Default().Errorf("Failed to update contact attributes: %v", err)
			continue
		}
	}

	// Send whatsapp message
	messagePayload := map[string]interface{}{
		"broadcast_name": "Feedback Reservation",
		"template_name":  messageTemplate.TemplateName,
		"receivers":      receivers,
	}

	payloadBytes, err := json.Marshal(messagePayload)
	if err != nil {
		logger.Default().Errorf("Failed to marshal message payload: %v", err)
		return false, status.Error(http.StatusInternalServerError, err.Error())
	}

	req, err := http.NewRequest("POST", integration.BaseURL+"/sendTemplateMessages", bytes.NewBuffer(payloadBytes))
	if err != nil {
		logger.Default().Errorf("Failed to create http request: %v", err)
		return false, status.Error(http.StatusInternalServerError, err.Error())
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+utils.Decrypt(revelIntegrationCredentials["api_token"].(string), os.Getenv("AES_ENCRYPTION_KEY")))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		logger.Default().Errorf("Failed to send http request: %v", err)
		return false, status.Error(http.StatusInternalServerError, err.Error())
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		logger.Default().Errorf("Failed to send whatsapp message: %v", err)
		return false, status.Error(http.StatusInternalServerError, "Failed to send whatsapp message")
	}

	sent = true

	return sent, nil
}

func (r *ReservationRepository) UpdateContactAttributes(ctx context.Context, phoneNumber string, branchId int32, customParams map[string]interface{}) error {
	// Get Whatsapp provider integration
	integration, err := r.CommonRepo.GetIntegrationBySystemType(ctx, "WhatsApp Platform", branchId)
	if err != nil {
		logger.Default().Errorf("Failed to get whatsapp integration: %v", err)
		return status.Error(http.StatusInternalServerError, "WhatsApp integration not found")
	}

	parseCredentials := []byte(integration.Credentials)
	var revelIntegrationCredentials map[string]interface{}
	if err := json.Unmarshal(parseCredentials, &revelIntegrationCredentials); err != nil {
		logger.Default().Errorf("Failed to unmarshal integration credentials: %v", err)
		return status.Error(http.StatusInternalServerError, "Failed to unmarshal integration credentials")
	}

	customParamsJson, err := json.Marshal(customParams)
	if err != nil {
		logger.Default().Errorf("Failed to marshal custom params: %v", err)
		return status.Error(http.StatusInternalServerError, "Failed to marshal custom params")
	}
	logger.Default().Errorf("customParamsJson: %v", customParamsJson)
	req, err := http.NewRequest("POST", integration.BaseURL+"/updateContactAttributes/"+phoneNumber, bytes.NewBuffer(customParamsJson))
	if err != nil {
		logger.Default().Errorf("Failed to create new request: %v", err)
		return status.Error(http.StatusInternalServerError, "Failed to create new request")
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+utils.Decrypt(revelIntegrationCredentials["api_token"].(string), os.Getenv("AES_ENCRYPTION_KEY")))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		logger.Default().Errorf("Failed to send http request: %v", err)
		return status.Error(http.StatusInternalServerError, "Failed to send http request")
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return status.Error(http.StatusInternalServerError, "Failed to update contact attributes")
	}

	return nil
}
