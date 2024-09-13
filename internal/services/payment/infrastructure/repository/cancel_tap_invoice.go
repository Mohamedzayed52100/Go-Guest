package repository

import (
	"encoding/json"
	"io"
	"net/http"
	"os"

	"github.com/goplaceapp/goplace-common/pkg/logger"
	"github.com/goplaceapp/goplace-common/pkg/utils"
	"github.com/goplaceapp/goplace-guest/internal/services/reservation/domain"
	"google.golang.org/grpc/status"
)

func (r *PaymentRepository) CancelTapInvoice(invoiceId string, integration *domain.Integration) error {
	var tapIntegrationCredentials map[string]interface{}
	
	if err := json.Unmarshal([]byte(integration.Credentials), &tapIntegrationCredentials); err != nil {
		return status.Error(http.StatusInternalServerError, err.Error())
	}

	request, _ := http.NewRequest("DELETE", integration.BaseURL+"/invoices/"+invoiceId, nil)
	request.Header.Set("Authorization", "Bearer "+utils.Decrypt(tapIntegrationCredentials["secret_key"].(string), os.Getenv("AES_ENCRYPTION_KEY")))
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Accept", "application/json")

	res, err := http.DefaultClient.Do(request)
	if err != nil {
		logger.Default().Error(err)
		return status.Error(http.StatusInternalServerError, err.Error())
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		logger.Default().Error(err)
		return status.Error(http.StatusInternalServerError, err.Error())
	}

	var jsonRes map[string]interface{}
	err = json.Unmarshal(body, &jsonRes)
	if err != nil {
		logger.Default().Error(err)
		return status.Error(http.StatusInternalServerError, err.Error())
	}

	if jsonRes["errors"] != nil {
		logger.Default().Error(jsonRes["errors"].([]interface{})[0].(map[string]interface{})["description"].(string))
		return status.Error(http.StatusInternalServerError, jsonRes["errors"].([]interface{})[0].(map[string]interface{})["description"].(string))
	}

	return nil
}
