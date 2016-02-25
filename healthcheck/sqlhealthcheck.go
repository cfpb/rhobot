package healthcheck

import (
	"fmt"
	"io/ioutil"
	"log"

	"database/sql"

	"gopkg.in/yaml.v2"
)

// SQLHealthCheck is a data type for storing the definition
// and results of a SQL based health check
type SQLHealthCheck struct {
	Expected string `yaml:"expected"`
	Query    string `yaml:"query"`
	Title    string `yaml:"title"`
	Error    bool   `yaml:"error"`
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
		fmt.Println(healthCheck.Query)
		rows, _ := cxn.Query(healthCheck.Query)
		var answer string
		rows.Next()
		rows.Scan(&answer)
		fmt.Println(answer)
		healthCheck.Passed = healthCheck.Expected == answer
		healthCheck.Actual = answer

		fmt.Printf("HEALTH CHECK: %s, Expected: %s, Found:%s\n", healthCheck.Title, healthCheck.Expected, answer)

	}
	return healthChecks
}

// EvaluateHealthChecks contains logic for handling the results of RunHealthChecks
func EvaluateHealthChecks(healthChecks Format) {
	var errors []SQLHealthCheck

	for _, healthCheck := range healthChecks.Tests {
		if !healthCheck.Passed && healthCheck.Error {
			errors = append(errors, healthCheck)
		}
	}

	if len(errors) > 0 {
		log.Fatalf("The folllowing health checks failed: %v", errors)
	}
}
