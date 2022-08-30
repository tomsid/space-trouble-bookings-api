package api

type Booking struct {
	FirstName     string `json:"firstName"`
	LastName      string `json:"lastName"`
	Gender        string `json:"gender"`
	Birthday      string `json:"birthday"`
	LaunchpadID   int    `json:"launchpadID"`
	DestinationID int    `json:"destinationID"`
	LaunchDate    string `json:"launchDate"`
}
