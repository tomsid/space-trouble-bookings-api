package db

import (
	"context"
	"time"
)

type Storage interface {
	Bookings(ctx context.Context) ([]Booking, error)
	CreateBooking(ctx context.Context, booking Booking) error
	Destinations(ctx context.Context) ([]Destination, error)
}

type Booking struct {
	ID            int
	FirstName     string
	LastName      string
	Gender        string
	Birthday      time.Time
	LaunchpadID   string
	DestinationID int
	LaunchDate    time.Time
}

type Destination struct {
	ID   int
	Name string
}
