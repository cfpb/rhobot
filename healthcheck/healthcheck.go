package healthcheck

import (
	"database/sql"
	"errors"
	"io/ioutil"
	"strings"

	log "github.com/Sirupsen/logrus"
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
	Equal    bool
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

// ReadHealthCheckYAMLFromFile loads healthcheck data from a YAML file
func ReadHealthCheckYAMLFromFile(path string) (format Format, err error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return
	}

	err = yaml.Unmarshal(data, &format)
	if err != nil {
		return
	}

	valid := format.ValidateHealthChecks()
	if !valid {
		err = errors.New("Reading Healthcheck file failed")
		return
	}

	return
}

// ValidateHealthChecks validates all healthchecks in specified file
func (healthChecks Format) ValidateHealthChecks() bool {
	for _, test := range healthChecks.Tests {
		if !test.ValidateHealthCheck() {
			return false
		}
	}
	return true
}

// RejectBadHealthChecks validates all healthchecks in specified file
func (healthChecks Format) RejectBadHealthChecks() Format {

	var GoodTests []SQLHealthCheck

	for _, test := range healthChecks.Tests {
		if test.ValidateHealthCheck() {
			GoodTests = append(GoodTests, test)
		}
	}

	healthChecks.Tests = GoodTests
	return healthChecks
}

// RunHealthChecks executes all health checks in the specified file
func (healthChecks Format) RunHealthChecks(cxn *sql.DB) Format {
	for _, test := range healthChecks.Tests {
		test.RunHealthCheck(cxn)
	}
	return healthChecks
}

// EvaluateHealthChecks contains logic for handling the results of RunHealthChecks
func (healthChecks Format) EvaluateHealthChecks() (results []SQLHealthCheck, errors []HCError) {
	for _, test := range healthChecks.Tests {
		results = append(results, test)
		hcErr := test.EvaluateHealthCheck()

		if hcErr.Err != "" {
			errors = append(errors, hcErr)
		}

		if hcErr.Exit {
			break
		}

	}
	return
}

// PreformHealthChecks runs and evaluates healthChecks one at a time
func (healthChecks Format) PreformHealthChecks(cxn *sql.DB) (results []SQLHealthCheck, errors []HCError) {
	for _, test := range healthChecks.Tests {
		test.RunHealthCheck(cxn)
		results = append(results, test)
		hcErr := test.EvaluateHealthCheck()

		if hcErr.Err != "" {
			errors = append(errors, hcErr)
		}

		if hcErr.Exit {
			break
		}

	}
	return
}

// ValidateHealthCheck makes sure a helathcheck has all the fields populated
func (healthCheck SQLHealthCheck) ValidateHealthCheck() bool {

	if len(healthCheck.Expected) == 0 {
		return false
	}

	if len(healthCheck.Query) == 0 {
		return false
	}

	if len(healthCheck.Title) == 0 {
		return false
	}

	if len(healthCheck.Severity) == 0 {
		return false
	}

	return true
}

// RunHealthCheck runs through a single healthcheck and saves the result
func (healthCheck SQLHealthCheck) RunHealthCheck(cxn *sql.DB) {
	answer := ""

	rows, err := cxn.Query(healthCheck.Query)

	if err != nil {
		log.Error(err)
		healthCheck.Passed = false
		healthCheck.Actual = err.Error()
	} else {

		rows.Next()
		rows.Scan(&answer)

		healthCheck.Passed = true
		healthCheck.Equal = healthCheck.Expected == answer
		healthCheck.Actual = answer
	}
}

// EvaluateHealthCheck runs through a single healthcheck and acts on the result
func (healthCheck *SQLHealthCheck) EvaluateHealthCheck() (err HCError) {

	prettyHealthCheck, _ := yaml.Marshal(&healthCheck)
	if !healthCheck.Equal || !healthCheck.Passed {
		switch strings.ToLower(healthCheck.Severity) {

		// When Fatal, return early with an error
		case "fatal":
			log.Errorf("FATAL healthcheck failed - Breaking Away Early\n%s", string(prettyHealthCheck))
			err = HCError{"FATAL healthCheck failure", true}

		// When Error, keep running but add an error
		case "error":
			log.Errorf("Healthcheck failed\n%s", string(prettyHealthCheck))
			err = HCError{"ERROR healthCheck failure", false}

		// When warn or info, print out the result and keep running
		case "warn":
			log.Warnf("Healthcheck failed\n%s", string(prettyHealthCheck))
		case "info":
			log.Infof("Healthcheck failed\n%s", string(prettyHealthCheck))
		default:
			log.Errorf("Undefined severity level:%s\n%s", strings.ToUpper(healthCheck.Severity), string(prettyHealthCheck))
		}
	} else {
		log.Printf("%s healthcheck passed\n%s\n", strings.ToUpper(healthCheck.Severity), string(prettyHealthCheck))
	}

	return err
}

// Implementation of report.Element

// HealthCheckReportHeaders headers used for GetHeaders
var HealthCheckReportHeaders = []string{"Title", "Query", "Passed", "Expected", "Actual", "Equal", "Severity"}

// GetHeaders Implementation for report.Element
func (healthCheck SQLHealthCheck) GetHeaders() []string {
	return HealthCheckReportHeaders[0:]
}

// GetValue Implementation for report.Element
func (healthCheck SQLHealthCheck) GetValue(key string) string {
	switch key {
	case HealthCheckReportHeaders[0]:
		return healthCheck.Title
	case HealthCheckReportHeaders[1]:
		return healthCheck.Query
	case HealthCheckReportHeaders[2]:
		if healthCheck.Passed {
			return "SUCCESS"
		}
		return "FAIL"
	case HealthCheckReportHeaders[3]:
		return healthCheck.Expected
	case HealthCheckReportHeaders[4]:
		return healthCheck.Actual
	case HealthCheckReportHeaders[5]:
		if healthCheck.Equal {
			return "TRUE"
		}
		return "FALSE"
	case HealthCheckReportHeaders[6]:
		return strings.ToUpper(healthCheck.Severity)
	}
	return ""
}
