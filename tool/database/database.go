package database

import (
	"context"
	"database/sql"
)

type Database interface {
	GetConnection(ctx context.Context) (*sql.Conn, error)
	CloseConnection(ctx context.Context, dbc *sql.Conn) error
	Close() error
}

type database struct {
	db                   *sql.DB
	maxConnectionRetries int
}
