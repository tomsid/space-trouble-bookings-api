package api

import (
	"context"
	"net/http"
	"net/http/httptest"
	"space-trouble-bookings-api/db"
	"space-trouble-bookings-api/spacex"
	"strings"
	"testing"
	"time"

	"go.uber.org/zap"
)

func TestAPI_BookFlight(t *testing.T) {
	destinations := []db.Destination{
		{ID: 1, Name: "Mars"},
		{ID: 2, Name: "Moon"},
		{ID: 3, Name: "Pluto"},
		{ID: 4, Name: "Asteroid Belt"},
		{ID: 5, Name: "Europa"},
		{ID: 6, Name: "Titan"},
		{ID: 7, Name: "Ganymede"},
	}
	testCases := []struct {
		name             string
		body             string
		launchPads       []spacex.Launchpad
		upcomingLaunches []spacex.Launch
		existingBookings []db.Booking
		expectedStatus   int
		expectedBody     string
	}{
		{
			name:           "invalid json",
			body:           "invalid json",
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "{\"message\":\"invalid character 'i' looking for beginning of value\"}",
		},
		{
			name:           "invalid launch date",
			body:           `{"launch_date": "invaliddate"}`,
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `{"message":"Invalid launch date. Should be in format YYYY-MM-DD: parsing time \"invaliddate\" as \"2006-01-02\": cannot parse \"invaliddate\" as \"2006\""}`,
		},
		{
			name:           "invalid birthday",
			body:           `{"launch_date": "2022-10-03", "birthday": "invalid"}`,
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `{"message":"Invalid birthday date. Should be in format YYYY-MM-DD: parsing time \"invalid\" as \"2006-01-02\": cannot parse \"invalid\" as \"2006\""}`,
		},
		{
			name:           "destination not found",
			body:           `{"launch_date": "2022-10-03", "birthday": "1993-04-18", "first_name": "fname", "last_name": "lname", "gender": "male", "destination_id": 9, "launchpad_id": "jwojeoijwfj"}`,
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `{"message":"Flight can't be booked: Destination with ID 9 not found"}`,
		},
		{
			name: "lauchpad is busy",
			body: `{"launch_date": "2022-10-03", "birthday": "1993-04-18", "first_name": "fname", "last_name": "lname", "gender": "male", "destination_id": 3, "launchpad_id": "jwojeoijwfj"}`,
			upcomingLaunches: []spacex.Launch{
				{
					Launchpad: "jwojeoijwfj", //we want this launchpad
					Name:      "",
					DateUTC:   "2022-10-03T05:40:00.000Z", //on the same day
				},
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `{"message":"Flight can't be booked: SpaceX uses the launchpad on that day"}`,
		},
		{
			name: "successfully book a ticket",
			body: `{"launch_date": "2022-10-08", "birthday": "1993-04-18", "first_name": "fname", "last_name": "lname", "gender": "male", "destination_id": 3, "launchpad_id": "jwojeoijwfj"}`,
			upcomingLaunches: []spacex.Launch{
				{
					Launchpad: "jwojeoijwfj", //we want this launchpad
					Name:      "",
					DateUTC:   "2022-10-25T05:40:00.000Z", //on a different day
				},
			},
			launchPads: []spacex.Launchpad{
				{
					ID: "jwojeoijwfj",
				},
			},
			expectedStatus: http.StatusCreated,
			expectedBody:   ``,
		},
		{
			name: "can't book a ticket for destination on a particular launchpad",
			body: `{"launch_date": "2022-10-03", "birthday": "1993-04-18", "first_name": "fname", "last_name": "lname", "gender": "male", "destination_id": 3, "launchpad_id": "jwojeoijwfj"}`,
			upcomingLaunches: []spacex.Launch{
				{
					Launchpad: "jwojeoijwfj", //we want this launchpad
					Name:      "",
					DateUTC:   "2022-10-25T05:40:00.000Z", //on a different day
				},
			},
			launchPads: []spacex.Launchpad{
				{
					ID: "jwojeoijwfj",
				},
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `{"message":"Flight can't be booked: No launches available for destination 3(Pluto) on launchpad jwojeoijwfj on 2022-10-03"}`,
		},
		{
			name: "launchpad not found",
			body: `{"launch_date": "2022-10-03", "birthday": "1993-04-18", "first_name": "fname", "last_name": "lname", "gender": "male", "destination_id": 3, "launchpad_id": "nonexisting"}`,
			upcomingLaunches: []spacex.Launch{
				{
					Launchpad: "jwojeoijwfj",
					Name:      "",
					DateUTC:   "2022-10-25T05:40:00.000Z",
				},
			},
			launchPads: []spacex.Launchpad{
				{
					ID: "jwojeoijwfj",
				},
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `{"message":"Flight can't be booked: Requested launchpad with ID \"nonexisting\" not found"}`,
		},
		{
			name: "schedule even if not according to schedule, but the destination is already scheduled before",
			body: `{"launch_date": "2022-10-03", "birthday": "1993-04-18", "first_name": "fname", "last_name": "lname", "gender": "male", "destination_id": 3, "launchpad_id": "jwojeoijwfj"}`,
			upcomingLaunches: []spacex.Launch{
				{
					Launchpad: "jwojeoijwfj",
					Name:      "",
					DateUTC:   "2022-10-25T05:40:00.000Z",
				},
			},
			launchPads: []spacex.Launchpad{
				{
					ID: "jwojeoijwfj",
				},
			},
			existingBookings: []db.Booking{
				{
					ID:            1,
					FirstName:     "terst",
					LastName:      "test_l",
					Gender:        "male",
					Birthday:      time.Now(),
					LaunchpadID:   "jwojeoijwfj",
					DestinationID: 3,
					LaunchDate:    time.Date(2022, 10, 3, 15, 34, 0, 0, time.UTC),
				},
			},
			expectedStatus: http.StatusCreated,
			expectedBody:   ``,
		},
	}

	for _, tc := range testCases {
		t.Log(tc.name)

		a := &API{
			spacex: &spacexMock{launchpads: tc.launchPads, upcomingLaunches: tc.upcomingLaunches},
			log:    zap.NewNop().Sugar(),
			db: &dbMock{
				destinations: destinations,
				bookings:     tc.existingBookings,
			},
		}

		resp := httptest.NewRecorder()
		a.BookFlight(resp, httptest.NewRequest("POST", "/booking", strings.NewReader(tc.body)))
		if tc.expectedStatus != resp.Code {
			t.Logf("unexpected status code. Got %d, want %d", resp.Code, tc.expectedStatus)
			t.Fail()
		}
		if tc.expectedBody != resp.Body.String() {
			t.Logf("unexpected body. Got %s, want %s", resp.Body.String(), tc.expectedBody)
			t.Fail()
		}
	}
}

type spacexMock struct {
	launchpads       []spacex.Launchpad
	upcomingLaunches []spacex.Launch
}

func (s *spacexMock) GetUpcomingLaunches(ctx context.Context) ([]spacex.Launch, error) {
	return s.upcomingLaunches, nil
}

func (s *spacexMock) GetAllLaunchpads(ctx context.Context) ([]spacex.Launchpad, error) {
	return s.launchpads, nil
}

type dbMock struct {
	destinations []db.Destination
	bookings     []db.Booking
}

func (m *dbMock) Bookings(ctx context.Context, filter db.BookingsFilter) ([]db.Booking, error) {
	return m.bookings, nil
}

func (m *dbMock) CreateBooking(ctx context.Context, booking db.Booking) error {
	m.bookings = append(m.bookings, db.Booking{
		ID:            len(m.bookings) + 1,
		FirstName:     booking.FirstName,
		LastName:      booking.LastName,
		Gender:        booking.Gender,
		Birthday:      time.Now(),
		LaunchpadID:   booking.LaunchpadID,
		DestinationID: booking.DestinationID,
		LaunchDate:    time.Now(),
	})
	return nil
}

func (m *dbMock) Destinations(ctx context.Context) ([]db.Destination, error) {
	return m.destinations, nil
}

func (m *dbMock) BookingExists(ctx context.Context, id int) (bool, error) {
	for _, booking := range m.bookings {
		if booking.ID == id {
			return true, nil
		}
	}

	return false, nil
}

func (m *dbMock) BookingDelete(ctx context.Context, id int) error {
	for i, booking := range m.bookings {
		if booking.ID == i {
			m.bookings = append(m.bookings[:i], m.bookings[i+1:]...)
			return nil
		}
	}

	return nil
}
