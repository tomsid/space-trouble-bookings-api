package api

import (
	"context"
	"net/http"
	"space-trouble-bookings-api/db"
	"strconv"
	"time"
)

type BookingsResponse struct {
	Bookings []Booking `json:"bookings"`
}

type Booking struct {
	ID            int    `json:"id"`
	FirstName     string `json:"first_name"`
	LastName      string `json:"last_name"`
	Gender        string `json:"gender"`
	Birthday      string `json:"birthday"`
	LaunchpadID   string `json:"launchpad_id"`
	DestinationID int    `json:"destination_id"`
	LaunchDate    string `json:"launch_date"`
}

func (a *API) Bookings(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()

	bookingsFilter := db.BookingsFilter{}
	q := r.URL.Query()
	if q.Has("launch_date") {
		launchDay, err := time.Parse(dateFormat, q.Get("launch_date"))
		if err != nil {
			a.writeBadRequest(w, ErrorResponse{Message: "launch_date should be in format YYYY-MM-DD"})
			return
		}
		bookingsFilter.LaunchDate = launchDay
	}

	if q.Has("offset") {
		offset, err := strconv.Atoi(q.Get("offset"))
		if err != nil || offset < 0 {
			a.writeBadRequest(w, ErrorResponse{Message: "offset should be an integer and be >=0"})
			return
		}
		bookingsFilter.Offset = offset
	}

	if q.Has("limit") {
		limit, err := strconv.Atoi(q.Get("limit"))
		if err != nil || limit < 1 || limit > 300 {
			a.writeBadRequest(w, ErrorResponse{Message: "limit should be an integer and be > 0 and <= 300"})
			return
		}
		bookingsFilter.Limit = limit
	}

	bookings, err := a.db.Bookings(ctx, bookingsFilter)
	if err != nil {
		a.log.Error(err)
		a.internalServerError(w)
		return
	}

	respBookings := make([]Booking, 0, len(bookings))
	for _, booking := range bookings {
		respBookings = append(respBookings, Booking{
			ID:            booking.ID,
			FirstName:     booking.FirstName,
			LastName:      booking.LastName,
			Gender:        booking.Gender,
			Birthday:      booking.Birthday.Format(dateFormat),
			LaunchpadID:   booking.LaunchpadID,
			DestinationID: booking.DestinationID,
			LaunchDate:    booking.LaunchDate.Format(dateFormat),
		})
	}

	a.writeJSONResponse(w, BookingsResponse{Bookings: respBookings})
}
