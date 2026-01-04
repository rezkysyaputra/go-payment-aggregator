package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"go-payment-aggregator/internal/config"
	"net/http"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
)

func main() {
	viperConfig := config.NewViper()
	logger := config.NewLogger(viperConfig)

	rdb := config.NewRedis(viperConfig, logger)

	fmt.Println("Worker started. Listening for webhooks on 'webhook_queue'...")

	ctx := context.Background()

	for {
		result, err := rdb.BLPop(ctx, 0*time.Second, "webhook_queue").Result()
		if err != nil {
			logger.Errorf("Redis connection error: %v", err)
			time.Sleep(5 * time.Second)
			continue
		}

		if len(result) < 2 {
			continue
		}

		payloadStr := result[1]
		go processWebhook(payloadStr, logger, rdb)
	}
}

type WebhookPayload struct {
	TransactionID string  `json:"transaction_id"`
	OrderID       string  `json:"order_id"`
	Status        string  `json:"status"`
	Amount        float64 `json:"amount"`
	Provider      string  `json:"provider"`
	CallbackURL   string  `json:"callback_url"`
	RetryCount    int     `json:"retry_count"`
}

func processWebhook(raw string, logger *logrus.Logger, rdb *redis.Client) {
	var payload WebhookPayload
	if err := json.Unmarshal([]byte(raw), &payload); err != nil {
		logger.Errorf("[ERROR] Invalid JSON payload: %v | Payload: %s", err, raw)
		return
	}

	if payload.CallbackURL == "" {
		logger.Errorf("[SKIP] No callback_url for Order: %s", payload.OrderID)
		return
	}

	logger.Infof("[PROCESSING] Sending webhook for Order %s to %s", payload.OrderID, payload.CallbackURL)

	merchantBody := map[string]any{
		"transaction_id": payload.TransactionID,
		"order_id":       payload.OrderID,
		"status":         payload.Status,
		"amount":         payload.Amount,
		"provider":       payload.Provider,
		"timestamp":      time.Now().Unix(),
	}

	jsonBody, _ := json.Marshal(merchantBody)

	client := http.Client{
		Timeout: 10 * time.Second,
	}

	resp, err := client.Post(payload.CallbackURL, "application/json", bytes.NewBuffer(jsonBody))
	if err != nil {
		logger.Errorf("[FAILED] Order %s: %v", payload.OrderID, err)
		retry(payload, logger, rdb)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		logger.Infof("[SUCCESS] Order %s: Merchant responded %d", payload.OrderID, resp.StatusCode)
	} else {
		logger.Errorf("[FAILED] Order %s: Merchant responded %d", payload.OrderID, resp.StatusCode)
		retry(payload, logger, rdb)
	}
}

func retry(payload WebhookPayload, logger *logrus.Logger, rdb *redis.Client) {
	maxRetry := 5
	if payload.RetryCount >= maxRetry {
		logger.Errorf("[GIVE UP] Max retry reached for Order %s", payload.OrderID)
		return
	}

	payload.RetryCount++
	newBody, _ := json.Marshal(payload)

	waitTime := time.Duration(payload.RetryCount*5) * time.Second

	logger.Warnf("[RETRY] Rescheduling Order %s in %v (Attempt %d/%d)", payload.OrderID, waitTime, payload.RetryCount, maxRetry)

	go func() {
		time.Sleep(waitTime)
		if err := rdb.RPush(context.Background(), "webhook_queue", newBody).Err(); err != nil {
			logger.Errorf("Failed to re-queue task: %v", err)
		}
	}()
}
