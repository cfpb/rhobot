package config

import (
	"fmt"
	"io"
	"net/url"
	"os"
	"regexp"

	log "github.com/Sirupsen/logrus"
)

// Config represents Rhobot's full set of configuration options
type Config struct {
	logFormatter log.Formatter
	logOutput    io.Writer
	logLevel     log.Level

	PgHost     string
	PgPort     string
	PgDatabase string
	PgUser     string
	PgPassword string
	PgSSLMode  string

	GOCDHost     string
	GOCDPort     string
	GOCDUser     string
	GOCDPassword string
	GOCDTimeout  string

	SMTPHost  string
	SMTPPort  string
	SMTPEmail string
	SMTPName  string
}

// NewDefaultConfig creates a new configuration object with default settings
func NewDefaultConfig() *Config {
	return &Config{
		logFormatter: &log.TextFormatter{},
		logOutput:    os.Stderr,
		logLevel:     log.InfoLevel,
		PgHost:       "localhost",
		PgPort:       "5432",
		PgDatabase:   "postgres",
		PgUser:       "postgres",
		PgPassword:   "password",
		PgSSLMode:    "disable",
		GOCDHost:     "http://localhost",
		GOCDPort:     "8153",
		GOCDUser:     "",
		GOCDPassword: "",
		GOCDTimeout:  "120",
		SMTPHost:     "localhost",
		SMTPPort:     "25",
		SMTPEmail:    "admin@localhost",
		SMTPName:     "admin",
	}
}

// NewConfig creates a new configuration object from environment variables
func NewConfig() (config *Config) {
	config = NewDefaultConfig()
	config.UpdateLogger()
	log.Debug("Created default configuration.")

	log.Debug("Loading settings from environment variables, when appropriate.")
	if os.Getenv("PGHOST") != "" {
		log.Debug("Retrieving value from PGHOST environment variable.")
		config.PgHost = os.Getenv("PGHOST")
	}

	if os.Getenv("PGPORT") != "" {
		log.Debug("Retrieving value from PGPORT environment variable.")
		config.PgPort = os.Getenv("PGPORT")
	}

	if os.Getenv("PGDATABASE") != "" {
		log.Debug("Retrieving value from PGDATABASE environment variable.")
		config.PgDatabase = os.Getenv("PGDATABASE")
	}

	if os.Getenv("PGUSER") != "" {
		log.Debug("Retrieving value from PGUSER environment variable.")
		config.PgUser = os.Getenv("PGUSER")
	}

	if os.Getenv("PGPASSWORD") != "" {
		log.Debug("Retrieving value from PGPASSWORD environment variable.")
		config.PgPassword = os.Getenv("PGPASSWORD")
	}

	if os.Getenv("PGSSLMODE") != "" {
		log.Debug("Retrieving value from PGSSLMODE environment variable.")
		config.PgSSLMode = os.Getenv("PGSSLMODE")
	}

	if os.Getenv("GOCDHOST") != "" {
		log.Debug("Retrieving value from GOCDHOST environment variable.")
		config.GOCDHost = os.Getenv("GOCDHOST")
	}

	if os.Getenv("GOCDPORT") != "" {
		log.Debug("Retrieving value from GOCDPORT environment variable.")
		config.GOCDPort = os.Getenv("GOCDPORT")
	}

	log.Debug("Retrieving value from GOCDUSER environment variable.")
	config.GOCDUser = os.Getenv("GOCDUSER")

	log.Debug("Retrieving value from GOCDPASSWORD environment variable.")
	config.GOCDPassword = os.Getenv("GOCDPASSWORD")

	if os.Getenv("GOCDTIMEOUT") != "" {
		log.Debug("Retrieving value from GOCDTIMEOUT environment variable.")
		config.GOCDTimeout = os.Getenv("GOCDTIMEOUT")
	}

	if os.Getenv("SMTPHOST") != "" {
		log.Debug("Retrieving value from SMTPHOST environment variable.")
		config.SMTPHost = os.Getenv("SMTPHOST")
	}

	if os.Getenv("SMTPPORT") != "" {
		log.Debug("Retrieving value from SMTPPORT environment variable.")
		config.SMTPPort = os.Getenv("SMTPPORT")
	}

	if os.Getenv("SMTPEMAIL") != "" {
		log.Debug("Retrieving value from SMTPEMAIL environment variable.")
		config.SMTPEmail = os.Getenv("SMTPEMAIL")
	}

	if os.Getenv("SMTPNAME") != "" {
		log.Debug("Retrieving value from SMTPNAME environment variable.")
		config.SMTPName = os.Getenv("SMTPNAME")
	}

	return
}

// UpdateLogger sets the logger configurtaion based on the Rhobot config
func (config *Config) UpdateLogger() {
	log.SetFormatter(config.logFormatter)
	log.SetOutput(config.logOutput)
	log.SetLevel(config.logLevel)
}

// SetLogLevel sets the global log level
func (config *Config) SetLogLevel(level string) {
	var err error
	if config.logLevel, err = log.ParseLevel(level); err != nil {
		log.Error(err)
		return
	}
	config.UpdateLogger()
}

// SetDBURI extracts Postgres connection variables from a DB URI
func (config *Config) SetDBURI(dbURI string) {
	dbURIRegex := regexp.MustCompile(`postgres://(?P<pg_user>\w+):(?P<pg_password>.+?)@(?P<pg_host>\w+):(?P<pg_port>\w+)/(?P<pg_database>\w+).*`)

	match := dbURIRegex.FindStringSubmatch(dbURI)
	if match == nil {
		log.Error("Unable to set DB connection parameters, invalid DB URI!")
	} else if len(match) < 4 {
		log.Error("Unable to set DB connection parameters, too few regex matches with DB URI!")
	}

	config.PgUser = match[1]
	config.PgPassword = match[2]
	config.PgHost = match[3]
	config.PgDatabase = match[4]
}

// DBURI generates a DB URI from the proper configruation options
func (config *Config) DBURI() (dbURI string) {
	parsedURI := make(map[string]interface{})

	//URL Encode Postgres variables
	parsedURI["PgUser"] = url.QueryEscape(config.PgUser)
	parsedURI["PgPassword"] = url.QueryEscape(config.PgPassword)
	parsedURI["PgHost"] = url.QueryEscape(config.PgHost)
	parsedURI["PgPort"] = url.QueryEscape(config.PgPort)
	parsedURI["PgDatabase"] = url.QueryEscape(config.PgDatabase)
	parsedURI["sslmode"] = config.PgSSLMode

	dbURI = fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
		parsedURI["PgUser"],
		parsedURI["PgPassword"],
		parsedURI["PgHost"],
		parsedURI["PgPort"],
		parsedURI["PgDatabase"],
		parsedURI["sslmode"])

	return
}

// SetGoCDHost sets the host value of the GoCD server
func (config *Config) SetGoCDHost(host string) {
	config.GOCDHost = host
}
