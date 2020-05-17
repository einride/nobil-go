package nobil

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
)

// SearchNearRequest is the input to the SearchNear method.
type SearchNearRequest struct {
	// Coordinate to search around.
	Coordinate LatLng
	// DistanceMetres of radius to search within.
	DistanceMetres int
	// Limit number of results to return.
	Limit int
}

// SearchNearResponse is the output from the SearchNear method.
type SearchNearResponse struct {
	// Results from the search.
	Results []*SearchNearResult
	// Raw JSON search response.
	Raw json.RawMessage
}

// SearchNearResult is a search result from the SearchNear method.
type SearchNearResult struct {
	// ChargingStation in the current search result.
	ChargingStation *ChargingStation
	// DistanceMetres from the searched coordinate.
	DistanceMetres int
}

// SearchNear searches for chargers located within the radius of a coordinate.
func (c *Client) SearchNear(
	ctx context.Context,
	request *SearchNearRequest,
) (_ *SearchNearResponse, err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("search near %dm of %s: %w", request.DistanceMetres, request.Coordinate, err)
		}
	}()
	jsonRequest := struct {
		Action         string  `json:"action"`
		APIKey         string  `json:"apikey"`
		APIVersion     string  `json:"apiversion"`
		Type           string  `json:"type"`
		Latitude       float64 `json:"lat"`
		Longitude      float64 `json:"long"`
		DistanceMetres int     `json:"distance"`
		Limit          int     `json:"limit"`
	}{
		Action:         "search",
		APIKey:         c.apiKey,
		APIVersion:     "3",
		Type:           "near",
		Latitude:       request.Coordinate.Latitude,
		Longitude:      request.Coordinate.Longitude,
		DistanceMetres: request.DistanceMetres,
		Limit:          request.Limit,
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
		Provider         string `json:"Provider"`
		Rights           string `json:"Rights"`
		APIVersion       string `json:"apiver"`
		ChargingStations []*struct {
			Metadata struct {
				jsonMetadata
				DistanceMetres string `json:"distance"`
			} `json:"csmd"`
			jsonAttributes `json:"attr"`
		} `json:"chargerstations"`
	}
	if err := json.Unmarshal(data, &jsonResponse); err != nil {
		return nil, err
	}
	response := SearchNearResponse{
		Results: make([]*SearchNearResult, 0, len(jsonResponse.ChargingStations)),
		Raw:     json.RawMessage(data),
	}
	for _, chargingStation := range jsonResponse.ChargingStations {
		distanceMetres, err := strconv.Atoi(chargingStation.Metadata.DistanceMetres)
		if err != nil {
			return nil, err
		}
		response.Results = append(response.Results, &SearchNearResult{
			DistanceMetres: distanceMetres,
			ChargingStation: (&jsonChargingStation{
				jsonMetadata:   chargingStation.Metadata.jsonMetadata,
				jsonAttributes: chargingStation.jsonAttributes,
			}).ChargingStation(),
		})
	}
	return &response, nil
}
