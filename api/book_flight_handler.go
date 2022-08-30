package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sort"
	"space-trouble-bookings-api/db"
	"time"
)

type Booking struct {
	FirstName     string `json:"firstName"`
	LastName      string `json:"lastName"`
	Gender        string `json:"gender"`
	Birthday      string `json:"birthday"`
	LaunchpadID   string `json:"launchpadID"`
	DestinationID int    `json:"destinationID"`
	LaunchDate    string `json:"launchDate"`
}

func (a *API) BookFlight(w http.ResponseWriter, r *http.Request) {
	b, err := io.ReadAll(r.Body)
	if err != nil {
		a.log.Error(err)
		a.writeResponse(w, []byte(""))
		w.WriteHeader(http.StatusInternalServerError)
	}

	flightBooking := &Booking{}
	err = json.Unmarshal(b, flightBooking)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		a.writeJSONResponse(w, ErrorResponse{Message: err.Error()})
		return
	}

	launchDate, err := time.Parse("2006-01-02", flightBooking.LaunchDate)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		a.writeJSONResponse(w, ErrorResponse{Message: fmt.Sprintf("Invalid launch date. Shoulld be in formate YYYY-MM-DD: %s", err.Error())})
		return
	}

	// todo validation

	destinations := []db.Destinations{
		{ID: 1, Name: "Moon"},
		{ID: 2, Name: "sdf"},
		{ID: 3, Name: "dfg"},
		{ID: 4, Name: "ddf"},
		{ID: 5, Name: "we"},
		{ID: 6, Name: "dwe"},
		{ID: 7, Name: "sdfs"},
	}

	destinationsMap := make(map[int]string, len(destinations))
	for _, destination := range destinations {
		destinationsMap[destination.ID] = destination.Name
	}

	if _, found := destinationsMap[flightBooking.DestinationID]; !found {
		w.WriteHeader(http.StatusBadRequest)
		a.writeJSONResponse(w, ErrorResponse{Message: fmt.Sprintf("Destination with ID %d not found", flightBooking.DestinationID)})
		return
	}

	//TODO check if spacex have the requested launchpad booked
	//launches, err := a.spacex.GetUpcomingLaunches()
	//if err != nil {
	//	a.log.Error([]byte(err.Error()))
	//	a.internalServerError(w)
	//	return
	//}

	//TODO check if requested day on requested pad a flight with different destination is booked.
	// can happen if the input changed (added/removed lauchpad or destination. Don't book in that case.

	launchPads, err := a.spacex.GetAllLaunchpads()
	if err != nil {
		a.log.Error([]byte(err.Error()))
		a.internalServerError(w)
		return
	}

	var launchPadIDs []string
	var requestedLaunchpadFound bool
	for _, launchPad := range launchPads {
		if launchPad.ID == flightBooking.LaunchpadID {
			requestedLaunchpadFound = true
		}
		launchPadIDs = append(launchPadIDs, launchPad.ID)
	}

	if !requestedLaunchpadFound {
		w.WriteHeader(http.StatusBadRequest)
		a.writeJSONResponse(w, ErrorResponse{Message: fmt.Sprintf("Requested launchpad with ID %q not found", flightBooking.LaunchpadID)})
		return
	}
	sort.Strings(launchPadIDs)

	bookingDay := launchDate.Day()
	pad1Destination := bookingDay%len(destinations) + 1
	fmt.Println(pad1Destination)

	launchpadToDestination := make(map[string]int, len(launchPads))
	for i, id := range launchPadIDs {
		currentPaddestinationID := pad1Destination + i + 1
		if pad1Destination+i+1 <= len(destinations) {
			launchpadToDestination[id] = currentPaddestinationID
		} else {
			launchpadToDestination[id] = currentPaddestinationID - len(destinations)
		}
	}
	a.writeJSONResponse(w, launchpadToDestination)

	if launchpadToDestination[flightBooking.LaunchpadID] != flightBooking.DestinationID {
		errResp := ErrorResponse{
			Message: fmt.Sprintf(
				"No launches available for destinatin %d(%s) on launchpad %s on %s",
				flightBooking.DestinationID, destinationsMap[flightBooking.DestinationID], flightBooking.LaunchpadID, flightBooking.LaunchDate,
			)}
		a.writeJSONResponse(w, errResp)
		return
	}

	//TODO save the booking
}
