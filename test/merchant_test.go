package test

import (
	"bytes"
	"encoding/json"
	"go-payment-aggregator/internal/helper"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRegister(t *testing.T) {
	ClearAll()
	requestBody := map[string]interface{}{
		"name":         "test",
		"callback_url": "http://localhost:8080",
	}

	bodyJson, err := json.Marshal(requestBody)
	assert.Nil(t, err)

	request := httptest.NewRequest("POST", "/v1/merchant/register", bytes.NewReader(bodyJson))
	request.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	app.ServeHTTP(w, request)

	assert.Equal(t, http.StatusCreated, w.Code)

	var result helper.Response
	json.Unmarshal(w.Body.Bytes(), &result)
	assert.Equal(t, "merchant created", result.Message)
	assert.NotEmpty(t, result.Data)

}

func TestRegisterError(t *testing.T) {
	ClearAll()
	requestBody := map[string]interface{}{
		"name":         "",
		"callback_url": "",
	}

	bodyJson, err := json.Marshal(requestBody)
	assert.Nil(t, err)

	request := httptest.NewRequest("POST", "/v1/merchant/register", bytes.NewReader(bodyJson))
	request.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	app.ServeHTTP(w, request)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var result helper.Response
	json.Unmarshal(w.Body.Bytes(), &result)
	assert.NotEmpty(t, result.Error)
}

func TestRegisterDuplicate(t *testing.T) {
	ClearAll()
	requestBody := map[string]interface{}{
		"name":         "test",
		"callback_url": "http://localhost:8080",
	}

	bodyJson, err := json.Marshal(requestBody)
	assert.Nil(t, err)

	request := httptest.NewRequest("POST", "/v1/merchant/register", bytes.NewReader(bodyJson))
	request.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	app.ServeHTTP(w, request)

	assert.Equal(t, http.StatusCreated, w.Code)

	var result helper.Response
	json.Unmarshal(w.Body.Bytes(), &result)
	assert.Equal(t, "merchant created", result.Message)
	assert.NotEmpty(t, result.Data)

	w = httptest.NewRecorder()
	app.ServeHTTP(w, request)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	json.Unmarshal(w.Body.Bytes(), &result)
	assert.NotEmpty(t, result.Error)
}
