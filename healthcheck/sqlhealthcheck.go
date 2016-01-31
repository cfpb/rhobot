package healthcheck

import (
	"fmt"
	"io/ioutil"
	"log"

	"database/sql"

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

func unmarshalHealthChecks(yamldata []byte) HealthCheckFormat {

	var data HealthCheckFormat

	err := yaml.Unmarshal(yamldata, &data)
	if err != nil {
		log.Fatalf("ERROR: %v", err)
	}
	// inflate queryfiles
	return data
}

func ReadYamlFromFile(path string) HealthCheckFormat {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatal(err)
	}
	return unmarshalHealthChecks(data)
}

func RunHealthChecks(healthChecks HealthCheckFormat, cxn *sql.DB) HealthCheckFormat {

	for _, healthCheck := range healthChecks.HealthChecks {
		fmt.Println(healthCheck.Query)
		rows, _ := cxn.Query(healthCheck.Query)
		var answer string
		rows.Next()
		rows.Scan(&answer)
		fmt.Println(answer)
		healthCheck.Passed = healthCheck.Expected == answer
		healthCheck.Actual = answer

		fmt.Printf("HEALTH CHECK: %s, Expected: %s, Found:%s\n", healthCheck.Name, healthCheck.Expected, answer)

	}
	return healthChecks
}

func EvaluateHealthChecks(healthChecks HealthCheckFormat) {
	var errors []SqlHealthCheck

	for _, healthCheck := range healthChecks.HealthChecks {
		if !healthCheck.Passed && healthCheck.Error {
			errors = append(errors, healthCheck)
		}
	}

	if len(errors) > 0 {
		log.Fatalf("The folllowing health checks failed", errors)
	}
}
