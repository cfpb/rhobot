package config

import (
	"fmt"
	"io"
	"os"
	"regexp"

	log "github.com/Sirupsen/logrus"
)

// Config represents Rhobot's full set of configuration options
type Config struct {
	logFormatter log.Formatter
	logOutput    io.Writer
	logLevel     log.Level

	pgHost     string
	pgPort     string
	PgDatabase string
	pgUser     string
	pgPassword string

	gocdHost     string
	gocdPort     string
	gocdUser     string
	gocdPassword string

	SMTPHost string
	SMTPPort string
}

// NewDefaultConfig creates a new configuration object with default settings
func NewDefaultConfig() *Config {
	return &Config{
		logFormatter: &log.TextFormatter{},
		logOutput:    os.Stderr,
		logLevel:     log.DebugLevel, // TODO: Will be log.WarnLevel
		pgHost:       "localhost",
		pgPort:       "5432",
		PgDatabase:   "postgres",
		pgUser:       "postgres",
		pgPassword:   "password",
		gocdHost:     "http://localhost",
		gocdPort:     "8153",
		gocdUser:     "admin",
		gocdPassword: "password",
	}
}

// NewConfig creates a new configuration object from environment variables
func NewConfig() (config *Config) {
	log.Debug("Creating default configuration.")
	config = NewDefaultConfig()

	config.UpdateLogger()

	log.Debug("Loading settings from environment variables, when appropriate.")
	if os.Getenv("PGHOST") != "" {
		log.Debug("Retrieving value from PGHOST environment variable.")
		config.pgHost = os.Getenv("PGHOST")
	}

	if os.Getenv("PGPORT") != "" {
		log.Debug("Retrieving value from PGPORT environment variable.")
		config.pgPort = os.Getenv("PGPORT")
	}

	if os.Getenv("PGDATABASE") != "" {
		log.Debug("Retrieving value from PGDATABASE environment variable.")
		config.PgDatabase = os.Getenv("PGDATABASE")
	}

	if os.Getenv("PGUSER") != "" {
		log.Debug("Retrieving value from PGUSER environment variable.")
		config.pgUser = os.Getenv("PGUSER")
	}

	if os.Getenv("PGPASSWORD") != "" {
		log.Debug("Retrieving value from PGPASSWORD environment variable.")
		config.pgPassword = os.Getenv("PGPASSWORD")
	}

	if os.Getenv("GOCDHOST") != "" {
		log.Debug("Retrieving value from GOCDHOST environment variable.")
		config.gocdHost = os.Getenv("GOCDHOST")
	}

	if os.Getenv("GOCDPORT") != "" {
		log.Debug("Retrieving value from GOCDPORT environment variable.")
		config.gocdPort = os.Getenv("GOCDPORT")
	}

	if os.Getenv("GOCDUSER") != "" {
		log.Debug("Retrieving value from GOCDUSER environment variable.")
		config.gocdUser = os.Getenv("GOCDUSER")
	}

	if os.Getenv("GOCDPASSWORD") != "" {
		log.Debug("Retrieving value from GOCDPASSWORD environment variable.")
		config.gocdPassword = os.Getenv("GOCDPASSWORD")
	}

	if os.Getenv("SMTPHOST") != "" {
		config.SMTPHost = os.Getenv("SMTPHOST")
	}

	if os.Getenv("SMTPPORT") != "" {
		config.SMTPPort = os.Getenv("SMTPPORT")
	}

	return
}

// UpdateLogger sets the logger configurtaion based on the Rhobot config
func (config *Config) UpdateLogger() {
	log.SetFormatter(config.logFormatter)
	log.SetOutput(config.logOutput)
	log.SetLevel(config.logLevel)
}

// SetDBURI extracts Postgres connection variables from a DB URI
func (config *Config) SetDBURI(dbURI string) {
	dbURIRegex := regexp.MustCompile(`postgres://(?P<pg_user>\w+):(?P<pg_password>\w+)@(?P<pg_host>\w+):(?P<pg_port>\w+)/(?P<pg_database>\w+).*`)

	match := dbURIRegex.FindStringSubmatch(dbURI)
	if match == nil {
		log.Error("Unable to set DB connection parameters, invalid DB URI!")
	} else if len(match) < 4 {
		log.Error("Unable to set DB connection parameters, too few regex matches with DB URI!")
	}

	config.pgUser = match[1]
	config.pgPassword = match[2]
	config.pgHost = match[3]
	config.PgDatabase = match[4]
}

// DBURI generates a DB URI from the proper configruation options
func (config *Config) DBURI() (dbURI string) {
	dbURI = fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=require",
		config.pgUser,
		config.pgPassword,
		config.pgHost,
		config.pgPort,
		config.PgDatabase)

	return
}

// SetGoCDHost sets the host value of the GoCD server
func (config *Config) SetGoCDHost(host string) {
	config.gocdHost = host
}

// GoCDURL returns the host of the GoCD server
func (config *Config) GoCDURL() string {
	return fmt.Sprintf("%s:%s", config.gocdHost, config.gocdPort)
}
