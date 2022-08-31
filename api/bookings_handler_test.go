package api

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"space-trouble-bookings-api/db"
	"testing"
	"time"

	"go.uber.org/zap"
)

func TestAPI_Bookings(t *testing.T) {
	testCases := []struct {
		name           string
		queryParams    url.Values
		bookings       []db.Booking
		expectedStatus int
		expectedBody   string
	}{
		{
			name:        "invalid offset query param",
			queryParams: url.Values{"offset": []string{""}},
			bookings: []db.Booking{
				{
					ID:            1,
					FirstName:     "asd",
					LastName:      "dsd",
					Gender:        "male",
					Birthday:      time.Now(),
					LaunchpadID:   "saffsdf",
					DestinationID: 2,
					LaunchDate:    time.Now(),
				},
				{
					ID:            2,
					FirstName:     "dsf",
					LastName:      "tyyy",
					Gender:        "male",
					Birthday:      time.Now(),
					LaunchpadID:   "saffsdf",
					DestinationID: 5,
					LaunchDate:    time.Now(),
				},
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `{"message":"offset should be an integer and be more than 0"}`,
		},
		{
			name:        "invalid limit query param",
			queryParams: url.Values{"offset": []string{"3"}, "limit": []string{"invalid"}},
			bookings: []db.Booking{
				{
					ID:            1,
					FirstName:     "asd",
					LastName:      "dsd",
					Gender:        "male",
					Birthday:      time.Now(),
					LaunchpadID:   "saffsdf",
					DestinationID: 2,
					LaunchDate:    time.Now(),
				},
				{
					ID:            2,
					FirstName:     "dsf",
					LastName:      "tyyy",
					Gender:        "male",
					Birthday:      time.Now(),
					LaunchpadID:   "saffsdf",
					DestinationID: 5,
					LaunchDate:    time.Now(),
				},
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `{"message":"limit should be an integer and be more that 0 and less or equal 300"}`,
		},
		{
			name:        "success",
			queryParams: url.Values{},
			bookings: []db.Booking{
				{
					ID:            1,
					FirstName:     "asd",
					LastName:      "dsd",
					Gender:        "male",
					Birthday:      time.Now(),
					LaunchpadID:   "saffsdf",
					DestinationID: 2,
					LaunchDate:    time.Now(),
				},
				{
					ID:            2,
					FirstName:     "dsf",
					LastName:      "tyyy",
					Gender:        "male",
					Birthday:      time.Now(),
					LaunchpadID:   "saffsdf",
					DestinationID: 5,
					LaunchDate:    time.Now(),
				},
			},
			expectedStatus: http.StatusOK,
			expectedBody:   `{"bookings":[{"id":1,"first_name":"asd","last_name":"dsd","gender":"male","birthday":"2022-08-31","launchpad_id":"saffsdf","destination_id":2,"launch_date":"2022-08-31"},{"id":2,"first_name":"dsf","last_name":"tyyy","gender":"male","birthday":"2022-08-31","launchpad_id":"saffsdf","destination_id":5,"launch_date":"2022-08-31"}]}`,
		},
	}

	for _, tc := range testCases {
		t.Log(tc.name)

		a := &API{
			log: zap.NewNop().Sugar(),
			db: &dbMock{
				bookings: tc.bookings,
			},
		}

		resp := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "/booking", nil)
		req.URL.RawQuery = tc.queryParams.Encode()
		a.Bookings(resp, req)
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
