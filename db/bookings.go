package db

import (
	"context"
	"fmt"
	"strconv"
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

	fmt.Println(q)
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
