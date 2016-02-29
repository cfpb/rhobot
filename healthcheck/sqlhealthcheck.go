package healthcheck

import (
	"database/sql"
	"errors"
	"io/ioutil"
	"log"
	"strings"

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

// HCError is a error helper for knowing to exit early on a failed healthcheck
type HCError struct {
	Err  string
	Exit bool
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
func EvaluateHealthChecks(healthChecks Format) (results []SQLHealthCheck, err error) {

	for _, healthCheck := range healthChecks.Tests {
		results = append(results, healthCheck)
		hcErr := EvaluateHealthCheck(healthCheck)

		if hcErr.Err != "" {
			err = errors.New(hcErr.Err)
		}

		if hcErr.Exit {
			return results, err
		}

	}
	return results, err
}

// PreformHealthChecks runs and evaluates healthChecks one at a time
func PreformHealthChecks(healthChecks Format, cxn *sql.DB) (results []SQLHealthCheck, err error) {
	for _, healthCheck := range healthChecks.Tests {
		healthCheck = RunHealthCheck(healthCheck, cxn)
		results = append(results, healthCheck)
		hcErr := EvaluateHealthCheck(healthCheck)

		if hcErr.Err != "" {
			err = errors.New(hcErr.Err)
		}

		if hcErr.Exit {
			return results, err
		}

	}
	return results, err
}

// RunHealthCheck runs through a single healthcheck and saves the result
func RunHealthCheck(healthCheck SQLHealthCheck, cxn *sql.DB) SQLHealthCheck {
	rows, _ := cxn.Query(healthCheck.Query)
	var answer string
	rows.Next()
	rows.Scan(&answer)
	healthCheck.Passed = healthCheck.Expected == answer
	healthCheck.Actual = answer
	return healthCheck
}

// EvaluateHealthCheck runs through a single healthcheck and acts on the result
func EvaluateHealthCheck(healthCheck SQLHealthCheck) (err HCError) {

	if !healthCheck.Passed {
		prettyHealthCheck, _ := yaml.Marshal(&healthCheck)

		switch strings.ToLower(healthCheck.Severity) {

		// When Fatal, return early with an error
		case "fatal":
			prettyHealthCheck, _ := yaml.Marshal(&healthCheck)
			log.Printf("FATAL healthcheck failed\nBreaking Away Early\n%s\n\n", string(prettyHealthCheck))
			err = HCError{"FATAL healthCheck failure", true}

		// When Error, keep running but add an error
		case "error":
			prettyHealthCheck, _ := yaml.Marshal(&healthCheck)
			log.Printf("%s healthcheck failed\n %s\n\n", strings.ToUpper(healthCheck.Severity), string(prettyHealthCheck))
			err = HCError{"ERROR healthCheck failure", false}

		// When warn or info, print out the result and keep running
		case "warn", "info":
			prettyHealthCheck, _ := yaml.Marshal(&healthCheck)
			log.Printf("%s healthcheck failed\n %s\n\n", strings.ToUpper(healthCheck.Severity), string(prettyHealthCheck))
		default:
			log.Printf("undefined severity level:%s\n%s\n\n", strings.ToUpper(healthCheck.Severity), string(prettyHealthCheck))
		}
	}

	return err

}
