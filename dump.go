package nobil

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strconv"
)

// DumpRequest is the input to the Dump method.
type DumpRequest struct {
	// CountryCode to filter on (all countries if not specified).
	CountryCode string
	// FromDate to include in dump.
	FromDate Date
	// Format of the dump (XML if not specified).
	Format Format
	// NoRealTimeChargers filters out real-time chargers from the result set.
	NoRealTimeChargers bool
	// NoRealTimeData filters out chargers that have been updated by realtime data.
	NoRealTimeUpdates bool
}

// Dump all chargers in the database.
func (c *Client) Dump(ctx context.Context, request *DumpRequest) (_ io.ReadCloser, err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("dump: %w", err)
		}
	}()
	httpRequest, err := http.NewRequest(http.MethodGet, c.dumpURL, nil)
	if err != nil {
		return nil, err
	}
	query := httpRequest.URL.Query()
	query.Set("apikey", c.apiKey)
	query.Set("file", "false")
	if request.CountryCode != "" {
		query.Set("countrycode", request.CountryCode)
	}
	if request.FromDate != (Date{}) {
		query.Set("fromdate", request.FromDate.String())
	}
	if request.Format != "" {
		query.Set("format", string(request.Format))
	}
	query.Set("norealtime", strconv.FormatBool(request.NoRealTimeChargers))
	query.Set("nonimupdate", strconv.FormatBool(request.NoRealTimeUpdates))
	httpRequest.URL.RawQuery = query.Encode()
	response, err := c.httpClient.Do(httpRequest.WithContext(ctx))
	if err != nil {
		return nil, err
	}
	if response.StatusCode != http.StatusOK {
		_ = response.Body.Close()
		return nil, fmt.Errorf("non-successful status code: %v", response.StatusCode)
	}
	return response.Body, nil
}
