package midtrans

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/spf13/viper"
)

type SnapRequest struct {
	TransactionDetails TransactionDetails `json:"transaction_details"`
}

type TransactionDetails struct {
	OrderID     string  `json:"order_id"`
	GrossAmount float64 `json:"gross_amount"`
}

type SnapResponse struct {
	Token       string `json:"token"`
	RedirectURL string `json:"redirect_url"`
}

func CreateTransaction(orderID string, amount float64) (*SnapResponse, error) {
	serverKey := viper.GetString("MIDTRANS_SERVER_KEY")
	if serverKey == "" {
		return nil, fmt.Errorf("midtrans server key is not configured")
	}

	url := "https://app.sandbox.midtrans.com/snap/v1/transactions"
	body := SnapRequest{}
	body.TransactionDetails.OrderID = orderID
	body.TransactionDetails.GrossAmount = amount

	jsonData, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.SetBasicAuth(serverKey, "")

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	var snapRes SnapResponse
	if err := json.NewDecoder(res.Body).Decode(&snapRes); err != nil {
		return nil, err
	}

	return &snapRes, nil
}
