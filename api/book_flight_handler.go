package api

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sort"
	"space-trouble-bookings-api/db"
	"time"
)

const dateFormat = "2006-01-02"

type BookingRequest struct {
	FirstName     string `json:"first_name"`
	LastName      string `json:"last_name"`
	Gender        string `json:"gender"`
	Birthday      string `json:"birthday"`
	LaunchpadID   string `json:"launchpad_id"`
	DestinationID int    `json:"destination_id"`
	LaunchDate    string `json:"launch_date"`
}

func (a *API) BookFlight(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	b, err := io.ReadAll(r.Body)
	if err != nil {
		a.log.Error(err)
		a.writeResponse(w, []byte(""))
		w.WriteHeader(http.StatusInternalServerError)
	}

	flightBooking := &BookingRequest{}
	err = json.Unmarshal(b, flightBooking)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		a.writeJSONResponse(w, ErrorResponse{Message: err.Error()})
		return
	}

	launchDate, err := time.Parse("2006-01-02", flightBooking.LaunchDate)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		a.writeJSONResponse(w, ErrorResponse{Message: fmt.Sprintf("Invalid launch date. Shoulld be in format YYYY-MM-DD: %s", err.Error())})
		return
	}

	birthday, err := time.Parse("2006-01-02", flightBooking.Birthday)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		a.writeJSONResponse(w, ErrorResponse{Message: fmt.Sprintf("Invalid birthday date. Shoulld be in format YYYY-MM-DD: %s", err.Error())})
		return
	}

	if flightBooking.Gender != "male" && flightBooking.Gender != "female" {
		w.WriteHeader(http.StatusBadRequest)
		a.writeJSONResponse(w, ErrorResponse{Message: fmt.Sprintf("Gender should be male or female")})
		return
	}

	if len(flightBooking.FirstName) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		a.writeJSONResponse(w, ErrorResponse{Message: fmt.Sprintf("field first_name can't be empty")})
		return
	}

	if len(flightBooking.LastName) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		a.writeJSONResponse(w, ErrorResponse{Message: fmt.Sprintf("field last_name can't be empty")})
		return
	}

	destinations, err := a.db.Destinations(ctx)
	if err != nil {
		a.log.Error(err)
		a.internalServerError(w)
		return
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

	upcomingLaunches, err := a.spacex.GetUpcomingLaunches()
	if err != nil {
		a.log.Error([]byte(err.Error()))
		a.internalServerError(w)
		return
	}
	for _, upcomingLaunch := range upcomingLaunches {
		t, err := time.Parse(time.RFC3339, upcomingLaunch.DateUTC)
		if err != nil {
			a.log.Errorf("failed to pares upcoming launch time: %s", err.Error())
		}
		if t.Day() == launchDate.Day() && t.Month() == launchDate.Month() && t.Year() == launchDate.Year() {
			w.WriteHeader(http.StatusBadRequest)
			a.writeJSONResponse(w, ErrorResponse{Message: "SpaceX uses the launchpad on that dayÏ"})
			return
		}
	}

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

	bookingDay := launchDate.YearDay()
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

	err = a.db.CreateBooking(ctx, db.Booking{
		FirstName:     flightBooking.FirstName,
		LastName:      flightBooking.LastName,
		DestinationID: flightBooking.DestinationID,
		LaunchpadID:   flightBooking.LaunchpadID,
		Gender:        flightBooking.Gender,
		LaunchDate:    launchDate,
		Birthday:      birthday,
	})

	if err != nil {
		a.log.Error(err)
		a.internalServerError(w)
		return
	}

	w.WriteHeader(http.StatusCreated)
}
