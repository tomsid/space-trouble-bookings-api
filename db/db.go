package db

import (
	"context"
	"time"
)

type Storage interface {
	Bookings(ctx context.Context, filter BookingsFilter) ([]Booking, error)
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

type BookingsFilter struct {
	LaunchDate time.Time
	Offset     int
	Limit      int
}

type Destination struct {
	ID   int
	Name string
}
