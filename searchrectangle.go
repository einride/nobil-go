package nobil

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

// SearchRectangleRequest is the input to the SearchRectangle method.
type SearchRectangleRequest struct {
	// SouthWest coordinate of the rectangle.
	SouthWest LatLng
	// NothEast coordinate of the rectangle.
	NorthEast LatLng
	// ExistingIDs of chargers to not include in the search results.
	ExistingIDs []string
}

// SearchRectangleResponse is the output from the SearchRectangle method.
type SearchRectangleResponse struct {
	// ChargingStations returned from the search.
	ChargingStations []*ChargingStation
	// Raw JSON search response.
	Raw json.RawMessage
}

// SearchRectangle searches for chargers located within a rectangle.
func (c *Client) SearchRectangle(
	ctx context.Context,
	request *SearchRectangleRequest,
) (_ *SearchRectangleResponse, err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("search rectangle SW%s NE%s: %w", request.SouthWest, request.NorthEast, err)
		}
	}()
	jsonRequest := struct {
		Action      string                `json:"action"`
		APIKey      string                `json:"apikey"`
		APIVersion  string                `json:"apiversion"`
		Type        string                `json:"type"`
		ExistingIDs commaSeparatedStrings `json:"existingids"`
		DataType    string                `json:"dataType"`
		NorthEast   LatLng                `json:"northeast"`
		SouthWest   LatLng                `json:"southwest"`
		Format      string                `json:"format"`
	}{
		Action:      "search",
		APIKey:      c.apiKey,
		APIVersion:  "3",
		Type:        "rectangle",
		NorthEast:   request.NorthEast,
		SouthWest:   request.SouthWest,
		Format:      "json",
		ExistingIDs: request.ExistingIDs,
	}
	var body bytes.Buffer
	body.Grow(300) // expected request size
	if err := json.NewEncoder(&body).Encode(&jsonRequest); err != nil {
		return nil, err
	}
	httpRequest, err := http.NewRequest(http.MethodPost, c.searchURL, &body)
	if err != nil {
		return nil, err
	}
	httpRequest.Header.Set("Content-Type", "application/json")
	httpResponse, err := c.httpClient.Do(httpRequest.WithContext(ctx))
	if err != nil {
		return nil, err
	}
	data, err := ioutil.ReadAll(httpResponse.Body)
	if err != nil {
		_ = httpResponse.Body.Close()
		return nil, err
	}
	if err := httpResponse.Body.Close(); err != nil {
		return nil, err
	}
	var jsonResponse struct {
		Provider         string                 `json:"Provider"`
		Rights           string                 `json:"Rights"`
		APIVersion       string                 `json:"apiver"`
		ChargingStations []*jsonChargingStation `json:"chargerstations"`
	}
	if err := json.Unmarshal(data, &jsonResponse); err != nil {
		return nil, err
	}
	response := SearchRectangleResponse{
		ChargingStations: make([]*ChargingStation, 0, len(jsonResponse.ChargingStations)),
		Raw:              json.RawMessage(data),
	}
	for _, chargingStation := range jsonResponse.ChargingStations {
		response.ChargingStations = append(response.ChargingStations, chargingStation.ChargingStation())
	}
	return &response, nil
}
