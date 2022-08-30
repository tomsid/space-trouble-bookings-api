package spacex

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type Launchpad struct {
	Images          Images   `json:"images"`
	Name            string   `json:"name"`
	FullName        string   `json:"full_name"`
	Locality        string   `json:"locality"`
	Region          string   `json:"region"`
	Latitude        float64  `json:"latitude"`
	Longitude       float64  `json:"longitude"`
	LaunchAttempts  int      `json:"launch_attempts"`
	LaunchSuccesses int      `json:"launch_successes"`
	Rockets         []string `json:"rockets"`
	Timezone        string   `json:"timezone"`
	Launches        []string `json:"launches"`
	Status          string   `json:"status"`
	Details         string   `json:"details"`
	Id              string   `json:"id"`
}

func (c *client) GetAllLaunchpads() ([]Launchpad, error) {
	req, err := http.NewRequest(http.MethodGet, APIBaseURL+"v4/launchpads", nil)
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

	var launchpads []Launchpad
	err = json.Unmarshal(body, &launchpads)
	if err != nil {
		return nil, err
	}

	return launchpads, nil
}
