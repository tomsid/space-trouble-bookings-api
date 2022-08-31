package db

import (
	"context"

	"github.com/jackc/pgx/v4/pgxpool"
)

type pgstorage struct {
	pg *pgxpool.Pool
}

func NewPGStorage(pool *pgxpool.Pool) Storage {
	return &pgstorage{pg: pool}
}

func (s *pgstorage) CreateBooking(ctx context.Context, b Booking) error {
	_, err := s.pg.Exec(ctx, "INSERT INTO bookings ("+
		"first_name, last_name, gender, birthday, launchpad_id, destination_id, launch_date) VALUES "+
		"($1, $2, $3, $4, $5, $6, $7)",
		b.FirstName, b.LastName, b.Gender, b.Birthday, b.LaunchpadID, b.DestinationID, b.LaunchDate)
	return err
}

func (s *pgstorage) Destinations(ctx context.Context) ([]Destination, error) {
	rows, err := s.pg.Query(ctx, "SELECT id,name FROM destinations")
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var destinations []Destination
	for rows.Next() {
		var destination Destination
		err = rows.Scan(&destination.ID, &destination.Name)
		if err != nil {
			return nil, err
		}
		destinations = append(destinations, destination)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return destinations, nil
}

func (s *pgstorage) Bookings(ctx context.Context) ([]Booking, error) {
	rows, err := s.pg.Query(ctx, "SELECT id,first_name,last_name,gender,birthday,launchpad_id,destination_id,launch_date FROM bookings")
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var bookings []Booking
	for rows.Next() {
		var b Booking
		err = rows.Scan(&b.ID, &b.FirstName, &b.LastName, &b.Gender, &b.Birthday, &b.LaunchpadID, &b.DestinationID, &b.LaunchDate)
		if err != nil {
			return nil, err
		}
		bookings = append(bookings, b)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return bookings, nil
}
