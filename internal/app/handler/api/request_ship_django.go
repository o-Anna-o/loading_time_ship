package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"loading_time/internal/app/ds"

	"github.com/sirupsen/logrus"
)

const (
	DJANGO_URL   = "http://localhost:8000"
	DJANGO_TOKEN = "12345678"
)

// sendToDjangoLoadingTime — отправка асинхронного запроса в Django
func (h *RequestShipHandler) sendToDjangoLoadingTime(rs *ds.RequestShip) error {

	reqBody := map[string]interface{}{
		"request_ship_id": rs.RequestShipID,
		"containers_20ft": rs.Containers20ftCount,
		"containers_40ft": rs.Containers40ftCount,
	}

	client := &http.Client{Timeout: 30 * time.Second}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return err
	}

	req, err := http.NewRequest(
		http.MethodPost,
		DJANGO_URL+"/calculate_loading_time/",
		bytes.NewBuffer(jsonData),
	)
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", DJANGO_TOKEN)

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusAccepted {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("django error: %s", body)
	}

	logrus.Info("Async request sent to Django")
	return nil
}
