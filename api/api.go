package api

import (
	"fmt"
	"net/http"
	"space-trouble-bookings-api/spacex"
)

type API struct {
	spacex spacex.Client
}

func NewAPI(spacexClient spacex.Client) *API {
	return &API{
		spacex: spacexClient,
	}
}

func (a *API) BookFlight(w http.ResponseWriter, r *http.Request) {
	launches, err := a.spacex.GetUpcomingLaunches()
	if err != nil {
		w.Write([]byte(err.Error()))
		return
	}

	w.Write([]byte(fmt.Sprintf("%#v", launches)))

}

func (a *API) Bookings(w http.ResponseWriter, r *http.Request) {
	launchpads, err := a.spacex.GetAllLaunchpads()
	if err != nil {
		w.Write([]byte(err.Error()))
		return
	}
	w.Write([]byte(fmt.Sprintf("%#v", launchpads)))
}
