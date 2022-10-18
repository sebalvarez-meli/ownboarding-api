package database

import "time"

// default values
const (
	defaultMaxConnectionRetries = 3
)

// Config database service config
type Config struct {
	Driver               string         `yaml:"driver"`
	DbHost               string         `yaml:"host"`
	DbName               string         `yaml:"name"`
	DbUsername           string         `yaml:"user"`
	DbPassword           string         `yaml:"password"`
	ConnMaxLifetime      time.Duration  `yaml:"connection_max_life_time_seconds"`
	ConnReadTimeout      *time.Duration `yaml:"connection_read_timeout"`
	ConnWriteTimeout     *time.Duration `yaml:"connection_write_timeout"`
	ConnTimeout          *time.Duration `yaml:"connection_timeout"`
	MaxConnectionRetries int            `yaml:"max_connection_retries"`
	MaxIdleConns         int            `yaml:"max_idle_connections_per_host"`
	MaxOpenConns         int            `yaml:"max_open_connections"`
}
