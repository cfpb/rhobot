package healthcheck

import (
	"database/sql"
	"errors"
	"fmt"
	"io/ioutil"
	"strings"

	log "github.com/Sirupsen/logrus"
	"gopkg.in/yaml.v2"
)

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
func (healthChecks *Format) ValidateHealthChecks() bool {
	for _, test := range healthChecks.Tests {
		if !test.ValidateHealthCheck() {
			return false
		}
	}
	return true
}

// RejectBadHealthChecks validates all healthchecks in specified file
func (healthChecks *Format) RejectBadHealthChecks() {

	var GoodTests []SQLHealthCheck

	for _, test := range healthChecks.Tests {
		if test.ValidateHealthCheck() {
			GoodTests = append(GoodTests, test)
		}
	}

	healthChecks.Tests = GoodTests
}

// RunHealthChecks executes all health checks in the specified file
func (healthChecks *Format) RunHealthChecks(cxn *sql.DB) {
	for i := 0; i < len(healthChecks.Tests); i++ {
		healthChecks.Tests[i].RunHealthCheck(cxn)
	}
}

// PreformHealthChecks runs and evaluates healthChecks one at a time
func (healthChecks *Format) PreformHealthChecks(cxn *sql.DB) (results []SQLHealthCheck, errors []HCError) {
	for i, test := range healthChecks.Tests {
		if cxn != nil {
			test.RunHealthCheck(cxn)
		}
		results = append(results, test)
		hcErr := test.EvaluateHealthCheck()

		if hcErr.Err != "" {
			errors = append(errors, hcErr)
		}

		if hcErr.Exit {
			//add unfinished healthchecks
			for _, unfinished := range healthChecks.Tests[i+1:] {
				results = append(results, unfinished)
			}
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
func (healthCheck *SQLHealthCheck) RunHealthCheck(cxn *sql.DB) {
	answer := ""

	rows, err := cxn.Query(healthCheck.Query)
	if err != nil {
		log.Error(err)
		healthCheck.Passed = false
		healthCheck.Actual = err.Error()
	} else {
		defer rows.Close()
		rows.Next()
		rows.Scan(&answer)

		compResult := false
		switch strings.ToLower(healthCheck.Operation) {
		case "eq":
			compResult = healthCheck.Expected == answer
		case "ne":
			compResult = healthCheck.Expected != answer
		case "lt":
			compResult = healthCheck.Expected < answer
		case "le":
			compResult = healthCheck.Expected <= answer
		case "gt":
			compResult = healthCheck.Expected > answer
		case "ge":
			compResult = healthCheck.Expected >= answer
		default:
			log.Info("opperation not specified, checking if equals")
			compResult = healthCheck.Expected == answer
		}

		healthCheck.Passed = true
		healthCheck.Actual = answer
		healthCheck.Equal = compResult
	}

}

// EvaluateHCErrors given a slice of HCErrors, determine if error or early exit
func EvaluateHCErrors(hcerrors []HCError) (int, int, bool) {
	numErrors := 0
	numWarnings := 0
	fatal := false
	for _, hcerr := range hcerrors {
		if strings.Contains(strings.ToUpper(hcerr.Err), "FATAL") {
			fatal = true
		}
		if strings.Contains(strings.ToUpper(hcerr.Err), "ERROR") {
			numErrors = numErrors + 1
		}
		if strings.Contains(strings.ToUpper(hcerr.Err), "WARN") {
			numWarnings = numWarnings + 1
		}
	}
	return numErrors, numWarnings, fatal
}

// EvaluateHealthCheck runs through a single healthcheck and acts on the result
func (healthCheck *SQLHealthCheck) EvaluateHealthCheck() (err HCError) {

	prettyHealthCheck, _ := yaml.Marshal(&healthCheck)
	severity := strings.ToUpper(healthCheck.Severity)
	if !healthCheck.Equal || !healthCheck.Passed {
		earlyExit := false
		errorMsg := fmt.Sprintf("%s - healthcheck failed \n%s",
			severity, string(prettyHealthCheck))

		switch strings.ToUpper(healthCheck.Severity) {
		case "FATAL":
			log.Errorf("Breaking Away Early \n%s ", errorMsg)
			earlyExit = true
		case "ERROR":
			log.Error(errorMsg)
		case "WARN":
			log.Warn(errorMsg)
		case "INFO":
			log.Info(errorMsg)
		case "DEBUG":
			log.Debug(errorMsg)
		default:
			log.Errorf("Breaking Away Early %s\n%s ", severity, errorMsg)
		}
		err = HCError{errorMsg, earlyExit}
	} else {
		happyMsg := fmt.Sprintf("%s - healthcheck passed \n%s",
			severity, string(prettyHealthCheck))
		log.Print(happyMsg)
	}
	return err
}

// Implementation of report.Element

// HealthCheckReportHeaders headers used for GetHeaders
var HealthCheckReportHeaders = []string{"Title", "Query", "Passed", "Expected", "Actual", "Equal", "Severity", "Operation"}

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
	case HealthCheckReportHeaders[7]:
		return strings.ToUpper(healthCheck.Operation)
	}
	return ""
}
