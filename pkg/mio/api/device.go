package mio_api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type Device struct {
	DeviceID        string    `json:"deviceId"`
	IMEI            string    `json:"imei"`
	Status          string    `json:"status"`
	Iccid           string    `json:"iccid"`
	ModelNumber     string    `json:"modelNumber"`
	FirmwareVersion string    `json:"firmwareVersion"`
	SerialNumber    string    `json:"serialNumber"`
	CreatedAt       time.Time `json:"createdAt"`
}

type DeviceResponse struct {
	Items []Device `json:"items"`
}

func (c *Client) GetDeviceList() (*DeviceResponse, error) {
	url := fmt.Sprintf("%s/devices", baseURL)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("x-api-key", c.APIKey)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API responded with status code %d", resp.StatusCode)
	}

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var deviceResponse DeviceResponse
	if err := json.Unmarshal(bodyBytes, &deviceResponse); err != nil {
		return nil, err
	}

	return &deviceResponse, nil
}

func (c *Client) GetDevice(deviceId string) (*Device, error) {
	url := fmt.Sprintf("%s/devices/%s", baseURL, deviceId)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("x-api-key", c.APIKey)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API responded with status code %d", resp.StatusCode)
	}

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var device Device
	if err := json.Unmarshal(bodyBytes, &device); err != nil {
		return nil, err
	}

	return &device, nil
}
