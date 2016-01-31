package healthcheck

import (
	"log"

	"gopkg.in/yaml.v2"
)

//Hello this should work
func Hello() string {
	return "Hello, healthy world.\n"
}

// SqlHealthCheckMetadata contains control information for a set of SqlHealthChecks
type SqlHealthCheckMetadata struct {
	Distribution []string
}

// SqlHealthCheck is a data type for storing the definition
// and results of a SQL based health check
type SqlHealthCheck struct {
	Name  string `yaml:"name"`
	Query string `yaml:"query"`
	// QueryFile string `yaml:"queryfile"`
	Expected string `yaml:"expected"`
	Actual   string
	Error    bool `yaml:"error"`
	Passed   bool
}

// HealthCheckFormat is for unmarshiling a healthcheck file
type HealthCheckFormat struct {
	Metadata     SqlHealthCheckMetadata `yaml:"metadata"`
	HealthChecks []SqlHealthCheck       `yaml:"healthchecks"`
}

func unmarshalHealthChecks(yamldata string) HealthCheckFormat {

	var data HealthCheckFormat

	// []byte conversion happens here
	err := yaml.Unmarshal([]byte(yamldata), &data)
	if err != nil {
		log.Fatalf("ERROR: %v", err)
	}

	// inflate queryfiles

	// log.Printf("HEALTH CHECKS: %v", data)

	return data
}
