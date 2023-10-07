package mio_api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type ShipmentAddress struct {
	FirstName   string `json:"firstName"`
	LastName    string `json:"lastName"`
	Country     string `json:"country"`
	ZipCode     string `json:"zipCode"`
	PhoneNumber string `json:"phoneNumber"`
	City        string `json:"city"`
	State       string `json:"state"`
	AddressLine string `json:"addressLine"`
}

type ShipmentItem struct {
	DeviceId string `json:"id"`
	Count    uint   `json:"count"`
}

type ShipmentRequest struct {
	Address ShipmentAddress `json:"address"`
	Items   []ShipmentItem  `json:"items"`
}

func (c *Client) ShipItem(req ShipmentRequest) error {
	url := fmt.Sprintf("%s/shipments", baseURL)

	reqBytes, err := json.Marshal(req)
	if err != nil {
		return err
	}

	httpReq, err := http.NewRequest("POST", url, bytes.NewReader(reqBytes))
	if err != nil {
		return err
	}
	httpReq.Header.Add("x-api-key", c.APIKey)

	resp, err := http.DefaultClient.Do(httpReq)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("API responded with status code %d", resp.StatusCode)
	}

	return nil
}
