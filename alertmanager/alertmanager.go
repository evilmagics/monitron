package alertmanager

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/go-resty/resty/v2"
)
// Alert represents a single alert sent to Alertmanager
type Alert struct {
	Labels      map[string]string `json:"labels"`
	Annotations map[string]string `json:"annotations"`
	StartsAt    time.Time         `json:"startsAt"`
	EndsAt      time.Time         `json:"endsAt"`
	Status      string            `json:"status"` // Added Status field
}

// AlertmanagerPayload is the structure for sending alerts to Alertmanager
type AlertmanagerPayload struct {
	Version     string            `json:"version"`
	GroupKey    string            `json:"groupKey"`
	TruncatedAlerts int               `json:"truncatedAlerts"`
	Status      string            `json:"status"`
	Receiver    string            `json:"receiver"`
	GroupLabels map[string]string `json:"groupLabels"`
	CommonLabels  map[string]string `json:"commonLabels"`
	CommonAnnotations map[string]string `json:"commonAnnotations"`
	ExternalURL string            `json:"externalURL"`
	Alerts      []Alert           `json:"alerts"`
}

// SendAlert sends a single alert to Alertmanager
func SendAlert(alertmanagerURL string, alert Alert) error {
	payload := AlertmanagerPayload{
		Version:     "4", // Alertmanager webhook API version
		GroupKey:    "<generated>",
		Status:      "firing",
		Receiver:    "monitron-receiver", // This should match a receiver in your Alertmanager config
		GroupLabels: map[string]string{"alertname": alert.Labels["alertname"]},
		CommonLabels:  alert.Labels,
		CommonAnnotations: alert.Annotations,
		Alerts:      []Alert{alert},
	}

	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal alert payload: %w", err)
	}

	client := resty.New()
	resp, err := client.R().
		SetHeader("Content-Type", "application/json").
		SetBody(payload).
		Post(alertmanagerURL)
	if err != nil {
		return fmt.Errorf("failed to send alert to Alertmanager: %w", err)
	}

	if resp.StatusCode() != http.StatusOK {
		return fmt.Errorf("Alertmanager returned non-OK status: %s", resp.Status())
	}

	log.Printf("Alert sent to Alertmanager: %s", alert.Labels["alertname"])
	return nil
}

// ReceiveWebhook handles incoming webhooks from Alertmanager (e.g., for notifications)
// This function is intended to be used as a Fiber handler, so it doesn't need to import fiber directly.
// The router will pass the fiber.Ctx object.
func ReceiveWebhook(c interface{}) error { // Changed type to interface{} to avoid direct fiber import
	// In a real scenario, you would cast c to *fiber.Ctx and then use its methods.
	// For now, we'll just log a message.
	log.Printf("Received Alertmanager webhook. (Fiber context not directly used here)")
	return nil
}


