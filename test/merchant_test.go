package test

import (
	"net/http/httptest"
	"testing"
)

func TestRegisterMerchant(t *testing.T) {
	// body req
	req := httptest.NewRequest("POST", "/merchants", nil)
}
