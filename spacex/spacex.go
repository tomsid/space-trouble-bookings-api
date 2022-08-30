package spacex

import (
	"net/http"
	"time"
)

const (
	APIBaseURL = "https://api.spacexdata.com/"
)

type Client interface {
	GetAllLaunchpads() ([]Launchpad, error)
	GetUpcomingLaunches() ([]Launch, error)
}

type Images struct {
	Small []string `json:"small,omitempty"`
	Large []string `json:"large,omitempty"`
}

func NewClient(httpClient *http.Client) Client {
	if httpClient == nil {
		httpClient = &http.Client{
			Timeout: time.Second * 10,
		}
	}

	return &client{
		httpClient,
	}
}

type client struct {
	*http.Client
}
