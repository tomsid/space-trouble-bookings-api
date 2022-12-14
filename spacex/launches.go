package spacex

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type Launch struct {
	Rocket          string `json:"rocket"`
	Launchpad       string `json:"launchpad"`
	FlightNumber    int    `json:"flight_number"`
	Name            string `json:"name"`
	DateUTC         string `json:"date_utc"`
	DateUnix        int    `json:"date_unix"`
	DateLocal       string `json:"date_local"`
	DatePrecision   string `json:"date_precision"`
	Upcoming        bool   `json:"upcoming"`
	LaunchLibraryID string `json:"launch_library_id"`
	ID              string `json:"id"`
}

func (c *client) GetUpcomingLaunches(ctx context.Context) ([]Launch, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, APIBaseURL+"v5/launches/upcoming", nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("response code no OK: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var launches []Launch
	err = json.Unmarshal(body, &launches)
	if err != nil {
		return nil, err
	}

	return launches, nil
}
