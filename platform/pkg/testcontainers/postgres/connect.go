package postgres

import (
	"context"

	"github.com/go-faster/errors"
	"github.com/jackc/pgx/v5"
)

func connectPostgresClient(ctx context.Context, uri string) (*pgx.Conn, error) {
	conn, err := pgx.Connect(ctx, uri)

	if err != nil {
		return nil, errors.Errorf("failed to parse postgres uri %q: %v", uri, err)
	}

	err = conn.Ping(ctx)
	if err != nil {
		return nil, errors.Errorf("failed to ping postgres uri %q: %v", uri, err)
	}

	return conn, nil
}
