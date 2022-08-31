package db

import (
	"context"
	"time"

	"github.com/jackc/pgx/v4/pgxpool"
)

type Storage interface {
	Bookings(ctx context.Context, filter BookingsFilter) ([]Booking, error)
	CreateBooking(ctx context.Context, booking Booking) error
	Destinations(ctx context.Context) ([]Destination, error)
	BookingExists(ctx context.Context, id int) (bool, error)
	BookingDelete(ctx context.Context, id int) error
}

type pgstorage struct {
	pg *pgxpool.Pool
}

func NewPGStorage(pool *pgxpool.Pool) Storage {
	return &pgstorage{pg: pool}
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
