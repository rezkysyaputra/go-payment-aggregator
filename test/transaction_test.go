package test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"go-payment-aggregator/internal/helper"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestCreateTransaction(t *testing.T) {
	ClearAll()
	merchant := CreateMerchant()

	requestBody := map[string]interface{}{
		"order_id": "order_123",
		"provider": "midtrans",
		"amount":   1000,
	}

	bodyJson, err := json.Marshal(requestBody)
	assert.Nil(t, err)

	request := httptest.NewRequest("POST", "/v1/transaction/", bytes.NewReader(bodyJson))
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("x-api-key", merchant.ApiKey)

	w := httptest.NewRecorder()
	app.ServeHTTP(w, request)

	assert.Equal(t, http.StatusCreated, w.Code)

	var result helper.Response
	json.Unmarshal(w.Body.Bytes(), &result)
	assert.Equal(t, "transaction created", result.Message)
	assert.NotEmpty(t, result.Data)
}

func TestCreateTransactionError(t *testing.T) {
	ClearAll()
	merchant := CreateMerchant()

	requestBody := map[string]interface{}{
		"order_id": "",
		"provider": "",
		"amount":   0,
	}

	bodyJson, err := json.Marshal(requestBody)
	assert.Nil(t, err)

	request := httptest.NewRequest("POST", "/v1/transaction/", bytes.NewReader(bodyJson))
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("x-api-key", merchant.ApiKey)

	w := httptest.NewRecorder()
	app.ServeHTTP(w, request)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var result helper.Response
	json.Unmarshal(w.Body.Bytes(), &result)
	assert.NotEmpty(t, result.Error)
}

func TestCreateTransactionUnauthorized(t *testing.T) {
	ClearAll()

	requestBody := map[string]interface{}{
		"order_id": "order_123",
		"provider": "midtrans",
		"amount":   1000,
	}

	bodyJson, err := json.Marshal(requestBody)
	assert.Nil(t, err)

	request := httptest.NewRequest("POST", "/v1/transaction/", bytes.NewReader(bodyJson))
	request.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	app.ServeHTTP(w, request)

	assert.Equal(t, http.StatusUnauthorized, w.Code)

	var result helper.Response
	json.Unmarshal(w.Body.Bytes(), &result)
	assert.Equal(t, "Unauthorized - API key missing", result.Error)
}

func TestCreateTransactionInvalidAPIKey(t *testing.T) {
	ClearAll()

	requestBody := map[string]interface{}{
		"order_id": "order_123",
		"provider": "midtrans",
		"amount":   1000,
	}

	bodyJson, err := json.Marshal(requestBody)
	assert.Nil(t, err)

	request := httptest.NewRequest("POST", "/v1/transaction/", bytes.NewReader(bodyJson))
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("x-api-key", "invalid_api_key")

	w := httptest.NewRecorder()
	app.ServeHTTP(w, request)

	assert.Equal(t, http.StatusUnauthorized, w.Code)

	var result helper.Response
	json.Unmarshal(w.Body.Bytes(), &result)
	assert.Equal(t, "Unauthorized - Invalid API key", result.Error)
}

func TestGetTransactionById(t *testing.T) {
	ClearAll()

	merchant := CreateMerchant()
	transaction := CreateTransaction(merchant.ID, 1000, "midtrans")

	fmt.Printf("Created merchant with ID: %s and API Key: %s\n", merchant.ID, merchant.ApiKey)
	request := httptest.NewRequest("GET", "/v1/transaction/"+transaction.ID.String(), nil)
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("x-api-key", merchant.ApiKey)

	w := httptest.NewRecorder()
	app.ServeHTTP(w, request)

	assert.Equal(t, http.StatusOK, w.Code)

	var result helper.Response
	json.Unmarshal(w.Body.Bytes(), &result)
	assert.Equal(t, "transaction found", result.Message)
	assert.NotEmpty(t, result.Data)
}

func TestGetTransactionByIdNotFound(t *testing.T) {
	ClearAll()

	merchant := CreateMerchant()

	fmt.Printf("Created merchant with ID: %s and API Key: %s\n", merchant.ID, merchant.ApiKey)
	request := httptest.NewRequest("GET", "/v1/transaction/"+uuid.New().String(), nil)
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("x-api-key", merchant.ApiKey)

	w := httptest.NewRecorder()
	app.ServeHTTP(w, request)

	assert.Equal(t, http.StatusNotFound, w.Code)

	var result helper.Response
	json.Unmarshal(w.Body.Bytes(), &result)
	assert.Equal(t, "transaction not found", result.Error)
}

func TestGetTransactionByOrderId(t *testing.T) {
	ClearAll()

	merchant := CreateMerchant()
	transaction := CreateTransaction(merchant.ID, 1000, "midtrans")

	fmt.Printf("Created merchant with ID: %s and API Key: %s\n", merchant.ID, merchant.ApiKey)
	request := httptest.NewRequest("GET", "/v1/transaction/order/"+transaction.OrderID, nil)
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("x-api-key", merchant.ApiKey)

	w := httptest.NewRecorder()
	app.ServeHTTP(w, request)

	assert.Equal(t, http.StatusOK, w.Code)

	var result helper.Response
	json.Unmarshal(w.Body.Bytes(), &result)
	assert.Equal(t, "transaction found", result.Message)
	assert.NotEmpty(t, result.Data)
}

func TestGetTransactionByOrderIdNotFound(t *testing.T) {
	ClearAll()

	merchant := CreateMerchant()

	fmt.Printf("Created merchant with ID: %s and API Key: %s\n", merchant.ID, merchant.ApiKey)
	request := httptest.NewRequest("GET", "/v1/transaction/order/wrong_order_id", nil)
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("x-api-key", merchant.ApiKey)

	w := httptest.NewRecorder()
	app.ServeHTTP(w, request)

	assert.Equal(t, http.StatusNotFound, w.Code)

	var result helper.Response
	json.Unmarshal(w.Body.Bytes(), &result)
	assert.Equal(t, "transaction not found", result.Error)
}
