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
