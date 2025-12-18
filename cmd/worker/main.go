package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"go-payment-aggregator/internal/config"
	"net/http"
	"time"

	"github.com/sirupsen/logrus"
)

func main() {
	// 1. Load Config
	viperConfig := config.NewViper()
	logger := config.NewLogger(viperConfig)

	// 2. Connect Redis
	rdb := config.NewRedis(viperConfig, logger)

	fmt.Println("Worker started. Listening for webhooks on 'webhook_queue'...")

	ctx := context.Background()

	for {
		// BLPOP blocks until data is available.
		// It returns a slice: [key_name, value]
		result, err := rdb.BLPop(ctx, 0*time.Second, "webhook_queue").Result()
		if err != nil {
			logger.Errorf("Redis connection error: %v", err)
			time.Sleep(5 * time.Second) // Wait before retrying
			continue
		}

		// Ensure we got data
		if len(result) < 2 {
			continue
		}

		payloadStr := result[1]
		go processWebhook(payloadStr, logger) // Spawn goroutine for each task to be concurrent
	}
}

type WebhookPayload struct {
	TransactionID string  `json:"transaction_id"`
	OrderID       string  `json:"order_id"`
	Status        string  `json:"status"`
	Amount        float64 `json:"amount"`
	Provider      string  `json:"provider"`
	CallbackURL   string  `json:"callback_url"`
}

func processWebhook(raw string, logger *logrus.Logger) {
	var payload WebhookPayload
	if err := json.Unmarshal([]byte(raw), &payload); err != nil {
		logger.Errorf("[ERROR] Invalid JSON payload: %v | Payload: %s", err, raw)
		return
	}

	if payload.CallbackURL == "" {
		logger.Errorf("[SKIP] No callback_url for Order: %s", payload.OrderID)
		return
	}

	logger.Errorf("[PROCESSING] Sending webhook for Order %s to %s", payload.OrderID, payload.CallbackURL)

	// Prepare payload for merchant (remove sensitive internal stuff if any)
	merchantBody := map[string]any{
		"transaction_id": payload.TransactionID,
		"order_id":       payload.OrderID,
		"status":         payload.Status,
		"amount":         payload.Amount,
		"provider":       payload.Provider,
		"timestamp":      time.Now().Unix(),
	}

	jsonBody, _ := json.Marshal(merchantBody)

	// HTTP Request with timeout
	client := http.Client{
		Timeout: 10 * time.Second,
	}

	resp, err := client.Post(payload.CallbackURL, "application/json", bytes.NewBuffer(jsonBody))
	if err != nil {
		logger.Errorf("[FAILED] Order %s: %v", payload.OrderID, err)
		// TODO: Implement Retry Logic:
		// e.g., using specific 'retry_queue' or 'LPUSH' back to 'webhook_queue' with delay
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		logger.Errorf("[SUCCESS] Order %s: Merchant responded %d", payload.OrderID, resp.StatusCode)
	} else {
		logger.Errorf("[FAILED] Order %s: Merchant responded %d", payload.OrderID, resp.StatusCode)
	}
}
