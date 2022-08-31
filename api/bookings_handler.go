package api

import (
	"context"
	"net/http"
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
	bookings, err := a.db.Bookings(ctx)
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
