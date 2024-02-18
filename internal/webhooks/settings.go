package wh

import (
	"bytes"
	"encoding/json"
	"errors"
	"go.uber.org/zap"
	"io"
	"net/http"
	"np_consumer/internal/models"
	"time"
)

type webhookMessage struct {
	TrackNumber string    `json:"trackNumber"`
	EventType   string    `json:"eventType"`
	EventTime   time.Time `json:"eventTime"`
}

func SendNotification(url string, ms models.Payload, logger *zap.Logger) error {
	data := webhookMessage{
		TrackNumber: ms.TrackingNumber,
		EventType:   ms.EventType,
		EventTime:   ms.EventTime,
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		logger.Error("failed to marshal ms", zap.Error(err))
		return err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		logger.Error("failed to create request", zap.Error(err))
		return err
	}

	req.Header.Set("content-type", "application/json")

	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		logger.Error("failed to send HTTP request")
		return err
	}
	defer func(Body io.ReadCloser) {
		err = Body.Close()
		if err != nil {
			logger.Error("failed to close response body", zap.Error(err))
			return
		}
	}(resp.Body)

	if resp.StatusCode != 200 {
		logger.Error("webhooks status code is not 200", zap.Int("statusCode", resp.StatusCode))
		return errors.New("status code is not 200")
	}

	return nil
}
