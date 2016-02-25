package healthcheck

import (
	"io/ioutil"
	"log"
	"strings"

	"database/sql"

	"gopkg.in/yaml.v2"
)

// SQLHealthCheck is a data type for storing the definition
// and results of a SQL based health check
type SQLHealthCheck struct {
	Expected string `yaml:"expected"`
	Query    string `yaml:"query"`
	Title    string `yaml:"title"`
	Severity string `yaml:"severity"`
	Passed   bool
	Actual   string
}

// Format is for unmarshiling a healthcheck file
// and contains control information for a set of SQLHealthChecks
type Format struct {
	Name         string           `yaml:"name"`
	Distribution []string         `yaml:"distribution"`
	Tests        []SQLHealthCheck `yaml:"tests"`
}

func unmarshalHealthChecks(yamldata []byte) Format {

	var data Format

	err := yaml.Unmarshal(yamldata, &data)
	if err != nil {
		log.Fatalf("ERROR: %v", err)
	}

	// inflate queryfiles
	return data
}

// ReadYamlFromFile loads healthcheck data from a YAML file
func ReadYamlFromFile(path string) Format {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatal(err)
	}
	return unmarshalHealthChecks(data)
}

// RunHealthChecks executes all health checks in the specified file
func RunHealthChecks(healthChecks Format, cxn *sql.DB) Format {

	for _, healthCheck := range healthChecks.Tests {
		healthCheck = RunHealthCheck(healthCheck, cxn)
	}
	return healthChecks
}

// EvaluateHealthChecks contains logic for handling the results of RunHealthChecks
func EvaluateHealthChecks(healthChecks Format) {
	var errors []SQLHealthCheck

	for _, healthCheck := range healthChecks.Tests {

		if !healthCheck.Passed {

			failed = append(failed, healthCheck)
			prettyHealthCheck, _ := yaml.Marshal(&healthCheck)

			switch strings.ToLower(healthCheck.Severity) {
			case "fatal":
				prettyHealthCheck, _ := yaml.Marshal(&healthCheck)
				log.Printf("FATAL healthcheck failed\nBreaking Away Early\n%s\n\n", string(prettyHealthCheck))
				break
			case "error", "warn", "info":
				prettyHealthCheck, _ := yaml.Marshal(&healthCheck)
				log.Printf("%s healthcheck failed\n %s\n\n", strings.ToUpper(healthCheck.Severity), string(prettyHealthCheck))
			default:
				log.Printf("undefined severity level:%s\n%s\n\n", strings.ToUpper(healthCheck.Severity), string(prettyHealthCheck))
			}
		}

		//TODO: Replace with reporting module to send back results of healthchecks
		if len(errors) > 0 {
			prettyFailed, _ := yaml.Marshal(&failed)
			log.Printf("The folllowing health checks failed\n %s\n\n", string(prettyFailed))
		}
	}
}

func RunHealthCheck(healthCheck SqlHealthCheck, cxn *sql.DB) SqlHealthCheck {
	rows, _ := cxn.Query(healthCheck.Query)
	var answer string
	rows.Next()
	rows.Scan(&answer)
	healthCheck.Passed = healthCheck.Expected == answer
	healthCheck.Actual = answer
	return healthCheck
}
