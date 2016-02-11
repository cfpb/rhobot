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

// SqlHealthCheck is a data type for storing the definition
// and results of a SQL based health check
type SqlHealthCheck struct {
	Expected string `yaml:"expected"`
	Query    string `yaml:"query"`
	Title    string `yaml:"title"`
	Error    bool   `yaml:"error"`
	Passed   bool
	Actual   string
}

// HealthCheckFormat is for unmarshiling a healthcheck file
// and contains control information for a set of SqlHealthChecks
type HealthCheckFormat struct {
	Name         string           `yaml:"name"`
	Distribution []string         `yaml:"distribution"`
	Tests        []SqlHealthCheck `yaml:"tests"`
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

func EvaluateHealthChecks(healthChecks HealthCheckFormat) {
	var errors []SqlHealthCheck

	for _, healthCheck := range healthChecks.Tests {
		if !healthCheck.Passed && healthCheck.Error {
			errors = append(errors, healthCheck)
		}
	}

	if len(errors) > 0 {
		log.Fatalf("The folllowing health checks failed", errors)
	}
}
