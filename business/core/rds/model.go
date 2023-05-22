package rds

// Config represents information about an Config database config.
type Config struct {
	Username string
	Password string
	Host     string
	Port     int
	Region   string
}

// NewConfig contains information needed to create a new Config database configuration.
type NewConfig struct {
	Username string
	Password string
	Host     string
	Port     int
	Region   string
}
