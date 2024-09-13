package repository

import (
	"encoding/json"
	"fmt"
	domain2 "github.com/goplaceapp/goplace-guest/pkg/reservationservice/domain"
	"net/http"
	"os"
	"time"

	"github.com/goplaceapp/goplace-common/pkg/utils"
	"github.com/goplaceapp/goplace-guest/internal/services/reservation/domain"
	"google.golang.org/grpc/status"
)

func (r *ReservationRepository) GetRevelOrderDetails(tableNumber int, revelIntegration *domain.Integration, reservation *domain2.Reservation, currentTime time.Time) (map[string]interface{}, error) {
	parseCredentials := []byte(revelIntegration.Credentials)
	var revelIntegrationCredentials map[string]interface{}
	if err := json.Unmarshal(parseCredentials, &revelIntegrationCredentials); err != nil {
		return nil, status.Error(http.StatusInternalServerError, err.Error())
	}

	if revelIntegrationCredentials["api_key"] == nil || revelIntegrationCredentials["api_secret"] == nil || revelIntegrationCredentials["establishment_id"] == nil {
		return nil, status.Error(http.StatusInternalServerError, "Revel integration credentials not found")
	}

	if reservation == nil || reservation.CheckIn == nil || reservation.CheckOut == nil {
		return nil, status.Error(http.StatusBadRequest, "Invalid reservation details")
	}

	if tableNumber == 0 {
		return nil, status.Error(http.StatusBadRequest, "Invalid table number")
	}

	params := map[string]interface{}{
		"api_key":       utils.Decrypt(revelIntegrationCredentials["api_key"].(string), os.Getenv("AES_ENCRYPTION_KEY")),
		"api_secret":    utils.Decrypt(revelIntegrationCredentials["api_secret"].(string), os.Getenv("AES_ENCRYPTION_KEY")),
		"establishment": utils.Decrypt(revelIntegrationCredentials["establishment_id"].(string), os.Getenv("AES_ENCRYPTION_KEY")),
		"created_date__range": fmt.Sprintf("%s,%s",
			reservation.CheckIn.Add(-20*time.Minute).Format("2006-01-02T15:04:05"),
			reservation.CheckOut.Format("2006-01-02T15:04:05"),
		),
		"order_by": "-created_date",
		"fields":   "id,table,final_total,prevailing_tax,tax,discount_amount,discount_reason,subtotal",
		"limit":    200,
	}

	req, err := http.NewRequest("GET", revelIntegration.BaseURL+"/resources/Order?"+utils.EncodeParams(params), nil)
	if err != nil {
		return nil, status.Error(http.StatusInternalServerError, err.Error())
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, status.Error(http.StatusInternalServerError, err.Error())
	}

	defer resp.Body.Close()
	var revelOrderDetails map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&revelOrderDetails); err != nil {
		return nil, status.Error(http.StatusInternalServerError, err.Error())
	}

	var revelOrderDetailsFiltered []map[string]interface{}

	for _, order := range revelOrderDetails["objects"].([]interface{}) {
		if order.(map[string]interface{})["table"] == nil {
			continue
		}

		if order.(map[string]interface{})["table"].(string) == fmt.Sprintf("/resources/Table/%d/", tableNumber) {
			revelOrderDetailsFiltered = append(revelOrderDetailsFiltered, order.(map[string]interface{}))
		}
	}

	if len(revelOrderDetailsFiltered) == 0 {
		return nil, status.Error(http.StatusNotFound, "No order found")
	}

	return revelOrderDetailsFiltered[0], nil
}

func (r *ReservationRepository) GetRevelOrderItems(orderId int, revelIntegration *domain.Integration) ([]map[string]interface{}, error) {
	parseCredentials := []byte(revelIntegration.Credentials)
	var revelIntegrationCredentials map[string]interface{}
	if err := json.Unmarshal(parseCredentials, &revelIntegrationCredentials); err != nil {
		return nil, status.Error(http.StatusInternalServerError, err.Error())
	}

	params := map[string]interface{}{
		"api_key":       utils.Decrypt(revelIntegrationCredentials["api_key"].(string), os.Getenv("AES_ENCRYPTION_KEY")),
		"api_secret":    utils.Decrypt(revelIntegrationCredentials["api_secret"].(string), os.Getenv("AES_ENCRYPTION_KEY")),
		"establishment": utils.Decrypt(revelIntegrationCredentials["establishment_id"].(string), os.Getenv("AES_ENCRYPTION_KEY")),
		"order":         orderId,
		"fields":        "id,product_name_override,price,quantity",
	}

	req, err := http.NewRequest("GET", revelIntegration.BaseURL+"/resources/OrderItem?"+utils.EncodeParams(params), nil)
	if err != nil {
		return nil, status.Error(http.StatusInternalServerError, err.Error())
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, status.Error(http.StatusInternalServerError, err.Error())
	}

	defer resp.Body.Close()
	var revelOrderDetails map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&revelOrderDetails); err != nil {
		return nil, status.Error(http.StatusInternalServerError, err.Error())
	}

	var revelOrderItems []map[string]interface{}
	for _, item := range revelOrderDetails["objects"].([]interface{}) {
		revelOrderItems = append(revelOrderItems, item.(map[string]interface{}))
	}

	return revelOrderItems, nil
}
