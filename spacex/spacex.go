package spacex

import (
	"context"
	"net/http"
	"time"
)

const (
	APIBaseURL = "https://api.spacexdata.com/"
)

type Client interface {
	GetAllLaunchpads(ctx context.Context) ([]Launchpad, error)
	GetUpcomingLaunches(ctx context.Context) ([]Launch, error)
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
