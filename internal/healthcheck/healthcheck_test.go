package healthcheck

import (
	"testing"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/cfpb/rhobot/internal/config"
	"github.com/cfpb/rhobot/internal/database"
	"github.com/cfpb/rhobot/internal/report"
)

var conf *config.Config

func init() {
	conf = config.NewConfig()
}

func TestUnmarshal(t *testing.T) {
	format, err := ReadHealthCheckYAMLFromFile("healthchecksTest.yml")
	if err != nil || format.Name != "rhobot healthcheck TEST" {
		t.Error("could not read file")
	}
}

func TestFileNotFound(t *testing.T) {
	_, err := ReadHealthCheckYAMLFromFile("should_not_exist.yml")
	if err == nil {
		t.Error("did not throw error as should have")
	}
}

func TestRunningBasicChecks(t *testing.T) {
	cxn := database.GetPGConnection(conf.DBURI())
	healthChecks, _ := ReadHealthCheckYAMLFromFile("healthchecksTest.yml")
	healthChecks.RunHealthChecks(cxn)

	for _, hc := range healthChecks.Tests {
		if !hc.Passed {
			log.Error("Basic HealthCheck Test Failed")
			t.Fail()
		}
	}
}

func TestEvaluatingBasicChecks(t *testing.T) {
	cxn := database.GetPGConnection(conf.DBURI())
	healthChecks, _ := ReadHealthCheckYAMLFromFile("healthchecksTest.yml")
	results, hcerrs := healthChecks.PreformHealthChecks(cxn)
	numErrors, numWarnings, _ := EvaluateHCErrors(hcerrs)

	if numWarnings != 1 {
		log.Error("numWarnings had the wrong length")
		t.Fail()
	}
	if numErrors > 0 {
		log.Error("numErrors had the wrong length")
		t.Fail()
	}
	if len(results) != 3 {
		log.Error("Healthchecks results had the wrong length")
		t.Fail()
	}
	for _, result := range results {
		if !result.Passed {
			log.Error("Basic HealthCheck Test Failed")
			t.Fail()
		}
	}
}

func TestEvaluatingErrorsChecks(t *testing.T) {
	cxn := database.GetPGConnection(conf.DBURI())
	healthChecks, _ := ReadHealthCheckYAMLFromFile("healthchecksErrors.yml")
	results, err := healthChecks.PreformHealthChecks(cxn)

	if err == nil {
		log.Error("Healthchecks did not throw an error, but should have")
		t.Fail()
	}
	if len(err) != 3 {
		log.Error("Healthcheck errors had the wrong length")
		t.Fail()
	}
	if len(results) != 5 {
		log.Error("Healthcheck results had the wrong length")
		t.Fail()
	}
}

func TestEvaluatingFatalChecks(t *testing.T) {
	cxn := database.GetPGConnection(conf.DBURI())
	healthChecks, _ := ReadHealthCheckYAMLFromFile("healthchecksFatal.yml")
	results, err := healthChecks.PreformHealthChecks(cxn)

	if err == nil {
		log.Error("Healthchecks did not throw an error, but should have")
		t.Fail()
	}

	if len(results) != 2 {
		log.Error("Healthcheck results had the wrong length")
		t.Fail()
	}
}

func TestEvaluatingIncompleteChecks(t *testing.T) {
	cxn := database.GetPGConnection(conf.DBURI())
	healthChecks, ferr := ReadHealthCheckYAMLFromFile("healthchecksIncomplete.yml")

	if ferr == nil {
		log.Error("Reading Healthcheck file did not throw an error, but should have")
		t.Fail()
	}

	healthChecks.RejectBadHealthChecks()
	results, err := healthChecks.PreformHealthChecks(cxn)

	if err != nil {
		log.Error("Evaluating Healthchecks threw an error")
		t.Fail()
	}

	if len(results) != 1 {
		log.Error("Healthcheck results had the wrong length")
		t.Fail()
	}
}

func TestPreformAllChecks(t *testing.T) {
	cxn := database.GetPGConnection(conf.DBURI())
	healthChecks, _ := ReadHealthCheckYAMLFromFile("healthchecksAll.yml")
	results, err := healthChecks.PreformHealthChecks(cxn)

	if err == nil {
		log.Error("Healthchecks did not throw an error, but should have")
		t.Fail()
	}
	if len(results) != 6 {
		log.Error("Healthcheck results had the wrong length")
		t.Fail()
	}
}

func TestPreformOperationsChecks(t *testing.T) {
	cxn := database.GetPGConnection(conf.DBURI())
	healthChecks, _ := ReadHealthCheckYAMLFromFile("healthchecksOperations.yml")
	results, err := healthChecks.PreformHealthChecks(cxn)

	if err != nil {
		if len(err) != 3 {
			log.Error("3 Errors weere expected")
			t.Fail()
		}
	}

	if len(results) != 10 {
		log.Error("Healthcheck results had the wrong length")
		t.Fail()
	}
}

func TestEvaluatingInvalidChecks(t *testing.T) {
	cxn := database.GetPGConnection(conf.DBURI())
	healthChecks, _ := ReadHealthCheckYAMLFromFile("healthchecksInvalid.yml")
	results, err := healthChecks.PreformHealthChecks(cxn)

	if err == nil {
		log.Error("Healthchecks did not throw an error, but should have")
		t.Fail()
	}
	if len(results) != 2 {
		log.Error("Healthcheck results had the wrong length")
		t.Fail()
	}
}

func TestSQLHealthCheckReportableElement(t *testing.T) {
	var hcr report.Element

	hcr = SQLHealthCheck{
		"true",
		"select (select count(1) from information_schema.tables) > 0;",
		"basic test", "FATAL",
		"equal",
		true,
		"t",
		true,
	}

	for _, header := range hcr.GetHeaders() {
		log.Debugf("%s : %s\n", header, hcr.GetValue(header))
	}

	if hcr.GetHeaders() == nil {
		log.Error("No headers in report ReportableElement")
		t.Fail()
	}
}

func TestHealthcheckPongo2Report(t *testing.T) {
	var rePass, reFail report.Element
	var rs report.Set
	var prr report.Runner
	var phr report.Handler

	rePass = SQLHealthCheck{"true", "select (select count(1) from information_schema.tables) > 0;", "basic test", "equal", "FATAL", true, "t", true}
	reFail = SQLHealthCheck{"true", "select (select count(1) from information_schema.tables) < 0;", "basic test", "equal", "FATAL", false, "f", true}
	prr = report.NewPongo2ReportRunnerFromString(TemplateHealthcheckHTML, true)
	phr = report.PrintHandler{}

	elements := []report.Element{rePass, reFail}
	metadata := map[string]interface{}{
		"name":      "TestHealthcheckPongo2Report",
		"db_name":   "testdb",
		"footer":    FooterHealthcheck,
		"timestamp": time.Now().UTC().String(),
	}
	rs = report.Set{Elements: elements, Metadata: metadata}

	reader, err := prr.ReportReader(rs)
	err = phr.HandleReport(reader)
	if err != nil {
		log.Errorf("Error writing report: %v", err)
		t.FailNow()
	}
}
