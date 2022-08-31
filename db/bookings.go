package db

import (
	"context"
	"fmt"
	"strconv"

	"github.com/jackc/pgx/v4"
)

func (s *pgstorage) Bookings(ctx context.Context, filter BookingsFilter) ([]Booking, error) {
	q := "SELECT id,first_name,last_name,gender,birthday,launchpad_id,destination_id,launch_date FROM bookings "
	var conditions []string
	if !filter.LaunchDate.IsZero() {
		conditions = append(conditions, fmt.Sprintf("CAST(launch_date as DATE) = '%s'", filter.LaunchDate.Format("2006-01-02")))
	}

	var paginationParams []string
	if filter.Limit != 0 {
		paginationParams = append(paginationParams, "LIMIT "+strconv.Itoa(filter.Limit))
	} else {
		paginationParams = append(paginationParams, "LIMIT 100")
	}
	if filter.Offset != 0 {
		paginationParams = append(paginationParams, "OFFSET "+strconv.Itoa(filter.Offset))
	}

	if len(conditions) > 0 {
		q += " WHERE "
		for _, condition := range conditions {
			q += condition + " "
		}
	}

	if len(paginationParams) > 0 {
		for _, paginationParam := range paginationParams {
			q += paginationParam + " "
		}
	}

	rows, err := s.pg.Query(ctx, q)
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

func (s *pgstorage) CreateBooking(ctx context.Context, b Booking) error {
	_, err := s.pg.Exec(ctx, "INSERT INTO bookings "+
		"(first_name, last_name, gender, birthday, launchpad_id, destination_id, launch_date) VALUES "+
		"($1, $2, $3, $4, $5, $6, $7)",
		b.FirstName, b.LastName, b.Gender, b.Birthday, b.LaunchpadID, b.DestinationID, b.LaunchDate)
	return err
}

func (s *pgstorage) BookingExists(ctx context.Context, id int) (bool, error) {
	row := s.pg.QueryRow(ctx, "SELECT id FROM bookings WHERE id = $1", id)
	var val int
	err := row.Scan(&val)
	if err != nil {
		if err == pgx.ErrNoRows {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func (s *pgstorage) BookingDelete(ctx context.Context, id int) error {
	_, err := s.pg.Exec(ctx, "DELETE FROM bookings WHERE id = $1", id)
	return err
}
