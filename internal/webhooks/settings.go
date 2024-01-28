package wh

import (
	"bytes"
	"encoding/json"
	"go.uber.org/zap"
	"io"
	"net/http"
	"np_consumer/internal/models"
	"strconv"
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

	logger.Info("webhooks status code: " + strconv.Itoa(resp.StatusCode))
	return nil
}
