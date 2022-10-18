package database

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

const ConnectionFormat = "%s:%s@tcp(%s)/%s?charset=utf8&parseTime=True"

func NewRepository(config Config) (*database, error) {
	// connection
	connectionString := fmt.Sprintf(ConnectionFormat, config.DbUsername, getEnv(config.DbPassword), getEnv(config.DbHost), config.DbName)
	// if config has ConnReadTimeout set, appends readTimeout param
	if config.ConnReadTimeout != nil {
		connectionString = fmt.Sprintf("%s&readTimeout=%s", connectionString, config.ConnReadTimeout.String())
	}
	// if config has ConnWriteTimeout set, appends writeTimeout param
	if config.ConnWriteTimeout != nil {
		connectionString = fmt.Sprintf("%s&writeTimeout=%s", connectionString, config.ConnWriteTimeout.String())
	}
	// if config has ConnTimeout set, appends timeout param
	if config.ConnTimeout != nil {
		connectionString = fmt.Sprintf("%s&timeout=%s", connectionString, config.ConnTimeout.String())
	}
	// open
	db, err := sql.Open(config.Driver, connectionString)
	if err != nil {
		return nil, err
	}
	// test
	err = db.Ping()
	if err != nil {
		return nil, err
	}
	// set connection pool configs
	db.SetMaxIdleConns(config.MaxIdleConns)
	db.SetMaxOpenConns(config.MaxOpenConns)
	db.SetConnMaxLifetime(config.ConnMaxLifetime * time.Second)
	retries := defaultMaxConnectionRetries
	if config.MaxConnectionRetries > 0 {
		retries = config.MaxConnectionRetries
	}
	return &database{
		db:                   db,
		maxConnectionRetries: retries,
	}, nil
}

func getEnv(name string) string {
	value := os.Getenv(name)
	if value != "" {
		return value
	}
	return name
}

func (d *database) GetConnection(ctx context.Context) (*sql.Conn, error) {
	var err error
	var conn *sql.Conn

	for retry := 0; retry < d.maxConnectionRetries; retry++ {
		// Obtain the connection
		conn, err = d.db.Conn(ctx)
		if err != nil {
			continue
		}
		// Test connection
		err = d.TestConnection(ctx)
		if err == nil {
			break
		}
	}
	return conn, err
}

// TestConnection tests the given connection
func (d *database) TestConnection(ctx context.Context) error {
	if ctx != nil {
		return d.db.PingContext(ctx)
	}
	return d.db.Ping()
}

// CloseConnection closes a given connection
func (d *database) CloseConnection(ctx context.Context, dbc *sql.Conn) error {
	if dbc == nil {
		return errors.New("you must send a sql.Conn in order to close")
	}
	// CloseConnection the connection
	err := dbc.Close()
	if err != nil {
		return err
	}
	// Set dbConn to nil
	dbc = nil
	// done
	return nil
}

func (d *database) Close() error {
	return d.db.Close()
}
