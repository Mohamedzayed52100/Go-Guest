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
	"github.com/goplaceapp/goplace-guest/internal/services/reservation/domain"
	guestDomain "github.com/goplaceapp/goplace-guest/pkg/guestservice/domain"
	reservationDomain "github.com/goplaceapp/goplace-guest/pkg/reservationservice/domain"
	"google.golang.org/grpc/status"
)

func (r *ReservationRepository) SendReservationWhatsappDetails(ctx context.Context, reservation *reservationDomain.Reservation) (bool, error) {
	var (
		sent bool
		err  error
	)

	// Get Whatsapp provider integration
	integration, err := r.CommonRepo.GetIntegrationBySystemType(ctx, "WhatsApp Platform", reservation.BranchID)
	if err != nil {
		return false, nil
	}

	// Get whatsapp template
	messageTemplate, err := r.GetClientWhatsappTemplate(ctx, reservation.BranchID, meta.WhatsappTemplateTypeCreateReservation)
	if err != nil {
		return false, nil
	}

	parseCredentials := []byte(integration.Credentials)
	var revelIntegrationCredentials map[string]interface{}
	if err := json.Unmarshal(parseCredentials, &revelIntegrationCredentials); err != nil {
		return false, status.Error(http.StatusInternalServerError, err.Error())
	}

	reservationTime, err := time.Parse("15:04:05", reservation.Time)
	if err != nil {
		return false, status.Error(http.StatusInternalServerError, err.Error())
	}

	// Define the cutoff time as 12 PM (00:00 to 11:59 belongs to the previous day)
	cutoffTime := time.Date(reservationTime.Year(), reservationTime.Month(), reservationTime.Day(), 12, 0, 0, 0, time.UTC)

	// If the reservation time is before the cutoff time, adjust the reservation date for the message
	if reservationTime.Before(cutoffTime) {
		reservation.Date = reservation.Date.AddDate(0, 0, 1)
	}

	formated12HourTime := reservationTime.Format("03:04 PM")
	dayName := reservation.Date.Weekday().String()
	arabicDay := meta.ArabicDays[dayName]

	var primaryGuest *guestDomain.Guest
	r.GetTenantDBConnection(ctx).First(&primaryGuest, "id = ?", reservation.GuestID)

	// Update contact attributes with new reservation details
	updateContactPayload := map[string]interface{}{
		"customParams": []map[string]interface{}{
			{
				"name":  "name",
				"value": primaryGuest.FirstName + " " + primaryGuest.LastName,
			},
			{
				"name":  "reservation_number",
				"value": strconv.Itoa(int(reservation.ID)),
			},
			{
				"name":  "date",
				"value": reservation.Date.Format("02 Jan 2006"),
			},
			{
				"name":  "day_ar",
				"value": arabicDay,
			},
			{
				"name":  "time",
				"value": formated12HourTime,
			},
			{
				"name":  "restaurant_branch",
				"value": reservation.Branch.Name,
			},
			{
				"name":  "person_count",
				"value": strconv.Itoa(int(reservation.GuestsNumber)),
			},
		},
	}

	customParamsJson, err := json.Marshal(updateContactPayload)
	if err != nil {
		return false, status.Error(http.StatusInternalServerError, err.Error())
	}

	req, err := http.NewRequest("POST", integration.BaseURL+"/updateContactAttributes/"+primaryGuest.PhoneNumber, bytes.NewBuffer(customParamsJson))
	if err != nil {
		return false, status.Error(http.StatusInternalServerError, err.Error())
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+utils.Decrypt(revelIntegrationCredentials["api_token"].(string), os.Getenv("AES_ENCRYPTION_KEY")))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return false, status.Error(http.StatusInternalServerError, err.Error())
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		logger.Default().Errorf("Failed to update contact attributes: %v", resp)
		return false, status.Error(http.StatusInternalServerError, "Failed to update contact attributes")
	}

	// Send whatsapp message
	messagePayload := map[string]interface{}{
		"broadcast_name": "Reservation",
		"template_name":  messageTemplate.TemplateName,
		"parameters":     updateContactPayload["customParams"],
	}

	payloadBytes, err := json.Marshal(messagePayload)
	if err != nil {
		return false, status.Error(http.StatusInternalServerError, err.Error())
	}

	req, err = http.NewRequest("POST", integration.BaseURL+"/sendTemplateMessage?whatsappNumber="+primaryGuest.PhoneNumber, bytes.NewBuffer(payloadBytes))
	if err != nil {
		return false, status.Error(http.StatusInternalServerError, err.Error())
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+utils.Decrypt(revelIntegrationCredentials["api_token"].(string), os.Getenv("AES_ENCRYPTION_KEY")))

	client = &http.Client{}
	resp, err = client.Do(req)
	if err != nil {
		return false, status.Error(http.StatusInternalServerError, err.Error())
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		logger.Default().Errorf("Failed to send whatsapp message: %v", resp)
		return false, status.Error(http.StatusInternalServerError, "Failed to send whatsapp message")
	}

	sent = true

	return sent, nil
}

func (r *ReservationRepository) GetClientWhatsappTemplate(ctx context.Context, branchId int32, templateType string) (*domain.WhatsappTemplate, error) {
	var template domain.WhatsappTemplate

	if err := r.GetTenantDBConnection(ctx).
		Model(&template).
		Where("template_type = ? AND branch_id = ?", templateType, branchId).
		First(&template).
		Error; err != nil {
		return nil, err
	}

	return &template, nil
}
