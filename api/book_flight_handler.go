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

	flightBooking := BookingRequest{}
	err = json.Unmarshal(b, &flightBooking)
	if err != nil {
		a.writeBadRequest(w, ErrorResponse{Message: err.Error()})
		return
	}

	launchDate, err := time.Parse("2006-01-02", flightBooking.LaunchDate)
	if err != nil {
		a.writeBadRequest(w, ErrorResponse{Message: fmt.Sprintf("Invalid launch date. Should be in format YYYY-MM-DD: %s", err.Error())})
		return
	}

	birthday, err := time.Parse("2006-01-02", flightBooking.Birthday)
	if err != nil {
		a.writeBadRequest(w, ErrorResponse{Message: fmt.Sprintf("Invalid birthday date. Should be in format YYYY-MM-DD: %s", err.Error())})
		return
	}

	if flightBooking.Gender != "male" && flightBooking.Gender != "female" {
		a.writeBadRequest(w, ErrorResponse{Message: "Gender should be male or female"})
		return
	}

	if len(flightBooking.FirstName) == 0 {
		a.writeBadRequest(w, ErrorResponse{Message: "field first_name can't be empty"})
		return
	}

	if len(flightBooking.LastName) == 0 {
		a.writeBadRequest(w, ErrorResponse{Message: "field last_name can't be empty"})
		return
	}

	err = a.flightSchedulable(ctx, flightBooking)
	if err != nil {
		if _, ok := err.(ScheduleError); ok {
			a.writeBadRequest(w, ErrorResponse{Message: err.Error()})
			return
		}

		a.log.Error(err)
		a.internalServerError(w)
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

func sameDay(t1 time.Time, t2 time.Time) bool {
	return t1.Day() == t2.Day() && t1.Month() == t2.Month() && t1.Year() == t2.Year()
}

func (a *API) flightSchedulable(ctx context.Context, flightBooking BookingRequest) error {
	launchDate, err := time.Parse("2006-01-02", flightBooking.LaunchDate)
	if err != nil {
		return err
	}

	destinations, err := a.getDestinationsMap(ctx)
	if err != nil {
		return err
	}

	if _, found := destinations[flightBooking.DestinationID]; !found {
		return ScheduleError{Reason: fmt.Sprintf("Destination with ID %d not found", flightBooking.DestinationID)}
	}

	busy, err := a.launchpadBusy(ctx, launchDate, flightBooking.LaunchpadID)
	if err != nil {
		return err
	}
	if busy {
		return ScheduleError{Reason: "SpaceX uses the launchpad on that day"}
	}

	launchDateBookings, err := a.db.Bookings(ctx, db.BookingsFilter{LaunchDate: launchDate})
	if err != nil {
		return err
	}
	if len(launchDateBookings) > 0 && launchDateBookings[0].DestinationID != flightBooking.DestinationID && launchDateBookings[0].LaunchpadID != flightBooking.LaunchpadID {
		return ScheduleError{Reason: fmt.Sprintf("On that day bookings only for destination %d are allowed", launchDateBookings[0].DestinationID)}
	}

	launchpadToDestination, err := a.getScheduleForDay(ctx, launchDate, flightBooking, destinations)
	if err != nil {
		return err
	}

	// if the launchpad's destination matches the client's requested booking destination
	if launchpadToDestination[flightBooking.LaunchpadID] != flightBooking.DestinationID {
		// or if there is a different destination on that day (overridden or the schedule shifted because of added/removed launchpads or destinations)
		// and user requested it
		if len(launchDateBookings) > 0 && launchDateBookings[0].DestinationID == flightBooking.DestinationID &&
			launchDateBookings[0].LaunchpadID == flightBooking.LaunchpadID {
			a.log.Info("According to timetable the flight shouldn't be scheduled, but scheduling anyway since on that day there are booking with that destination already")
		} else {
			return ScheduleError{fmt.Sprintf(
				"No launches available for destination %d(%s) on launchpad %s on %s",
				flightBooking.DestinationID, destinations[flightBooking.DestinationID],
				flightBooking.LaunchpadID, flightBooking.LaunchDate,
			)}
		}
	}

	return nil
}

func (a *API) getDestinationsMap(ctx context.Context) (map[int]string, error) {
	destinations, err := a.db.Destinations(ctx)
	if err != nil {
		return nil, err
	}

	destinationsMap := make(map[int]string, len(destinations))
	for _, destination := range destinations {
		destinationsMap[destination.ID] = destination.Name
	}

	return destinationsMap, nil
}

func (a *API) launchpadBusy(ctx context.Context, launchDate time.Time, launchpadID string) (bool, error) {
	upcomingLaunches, err := a.spacex.GetUpcomingLaunches(ctx)
	if err != nil {
		return false, err
	}
	for _, upcomingLaunch := range upcomingLaunches {
		t, err := time.Parse(time.RFC3339, upcomingLaunch.DateUTC)
		if err != nil {
			a.log.Errorf("failed to parse upcoming launch time: %s", err.Error())
		}
		if sameDay(t, launchDate) && upcomingLaunch.Launchpad == launchpadID {
			return true, nil
		}
	}

	return false, nil
}

func (a *API) getScheduleForDay(ctx context.Context, launchDate time.Time, flightBooking BookingRequest, destinations map[int]string) (map[string]int, error) {
	launchPads, err := a.spacex.GetAllLaunchpads(ctx)
	if err != nil {
		return nil, err
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
		return nil, ScheduleError{Reason: fmt.Sprintf("Requested launchpad with ID %q not found", flightBooking.LaunchpadID)}
	}
	sort.Strings(launchPadIDs)

	bookingDay := launchDate.YearDay()
	pad1Destination := bookingDay%len(destinations) + 1

	launchpadToDestination := make(map[string]int, len(launchPads))
	for i, id := range launchPadIDs {
		currentPadDestinationID := pad1Destination + i + 1
		if pad1Destination+i+1 <= len(destinations) {
			launchpadToDestination[id] = currentPadDestinationID
		} else {
			launchpadToDestination[id] = currentPadDestinationID - len(destinations)
		}
	}

	return launchpadToDestination, nil
}
