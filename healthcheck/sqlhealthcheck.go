package healthcheck

import (
	"io/ioutil"
	"log"
	"strings"

	"database/sql"

	"gopkg.in/yaml.v2"
)

// SqlHealthCheck is a data type for storing the definition
// and results of a SQL based health check
type SqlHealthCheck struct {
	Expected string `yaml:"expected"`
	Query    string `yaml:"query"`
	Title    string `yaml:"title"`
	Severity string `yaml:"severity"`
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
		healthCheck = RunHealthCheck(healthCheck, cxn)
	}
	return healthChecks
}

func EvaluateHealthChecks(healthChecks HealthCheckFormat) {
	var failed []SqlHealthCheck

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

		//         //TODO: Replace with reporting module to send back results of healthchecks
		// 		if len(failed) > 0 {
		// 			prettyFailed, _ := yaml.Marshal(&failed)
		// 			log.Printf("The folllowing health checks failed\n %s\n\n", string(prettyFailed))
		// 		}
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
