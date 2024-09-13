package repository

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"os"

	"github.com/goplaceapp/goplace-common/pkg/logger"
	"github.com/goplaceapp/goplace-common/pkg/meta"
	commonUtils "github.com/goplaceapp/goplace-common/pkg/utils"
	guestProto "github.com/goplaceapp/goplace-guest/api/v1"
	"github.com/goplaceapp/goplace-guest/internal/services/reservation/domain"
	"github.com/goplaceapp/goplace-guest/utils"
	"google.golang.org/grpc/status"
)

func (r *ReservationRepository) RequestReservationWebhook(ctx context.Context, req *guestProto.RequestReservationWebhookRequest) (*guestProto.RequestReservationWebhookResponse, error) {
	var (
		randomBranchID int32
		integration    domain.Integration
		err            error
		waTemplate     domain.WhatsappTemplate
	)

	req.PhoneNumber = utils.RemovePlusSign(req.PhoneNumber)

	if err := r.GetTenantDBConnection(ctx).
		Table("branches").
		Order("RANDOM()").
		Select("id").
		Limit(1).
		Scan(&randomBranchID).Error; err != nil {
		return nil, status.Error(http.StatusInternalServerError, err.Error())
	}

	// Get Whatsapp provider integration
	if err := r.GetTenantDBConnection(ctx).
		Model(&integration).
		Where("system_type = ?", "WhatsApp Platform").
		First(&integration).Error; err != nil {
		return nil, err
	}

	// Get whatsapp template
	if err := r.GetTenantDBConnection(ctx).
		Model(&waTemplate).
		Where("template_type = ?", meta.WhatsappTemplateRequestReservation).
		First(&waTemplate).
		Error; err != nil {
		return nil, err
	}

	parseCredentials := []byte(integration.Credentials)
	var revelIntegrationCredentials map[string]interface{}
	if err := json.Unmarshal(parseCredentials, &revelIntegrationCredentials); err != nil {
		return nil, status.Error(http.StatusInternalServerError, err.Error())
	}

	// Send whatsapp message
	messagePayload := map[string]interface{}{
		"broadcast_name": "Request Reservation",
		"template_name":  waTemplate.TemplateName,
	}

	payloadBytes, err := json.Marshal(messagePayload)
	if err != nil {
		return nil, status.Error(http.StatusInternalServerError, err.Error())
	}

	request, err := http.NewRequest("POST", integration.BaseURL+"/sendTemplateMessage?whatsappNumber="+req.PhoneNumber, bytes.NewBuffer(payloadBytes))
	if err != nil {
		return nil, status.Error(http.StatusInternalServerError, err.Error())
	}

	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Authorization", "Bearer "+commonUtils.Decrypt(revelIntegrationCredentials["api_token"].(string), os.Getenv("AES_ENCRYPTION_KEY")))

	client := &http.Client{}
	resp, err := client.Do(request)
	if err != nil {
		return nil, status.Error(http.StatusInternalServerError, err.Error())
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		logger.Default().Errorf("Failed to send whatsapp message: %v", resp)
		return nil, status.Error(http.StatusInternalServerError, "Failed to send whatsapp message")
	}

	return &guestProto.RequestReservationWebhookResponse{
		Code:    http.StatusOK,
		Message: "Reservation Request has been sent successfully",
	}, nil
}
