package healthcheck

import (
	"testing"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/cfpb/rhobot/config"
	"github.com/cfpb/rhobot/database"
	"github.com/cfpb/rhobot/report"
)

var conf *config.Config

func init() {
	conf = config.NewConfig()
}

func TestUnmarshal(t *testing.T) {
	format, err := ReadHealthCheckYAMLFromFile("healthchecksTest.yml")
	if err != nil || format.Name != "rhobot healthcheck test" {
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
}

func TestEvaluatingBasicChecks(t *testing.T) {
	cxn := database.GetPGConnection(conf.DBURI())
	healthChecks, _ := ReadHealthCheckYAMLFromFile("healthchecksTest.yml")
	healthChecks.RunHealthChecks(cxn)
	results, err := healthChecks.EvaluateHealthChecks()

	if err != nil {
		log.Error(err)
		t.Fail()
	}
	if len(results) != 3 {
		log.Error("Healthchecks results had the wrong length")
		t.Fail()
	}
}

func TestEvaluatingErrorsChecks(t *testing.T) {
	cxn := database.GetPGConnection(conf.DBURI())
	healthChecks, _ := ReadHealthCheckYAMLFromFile("healthchecksErrors.yml")
	healthChecks.RunHealthChecks(cxn)
	results, err := healthChecks.EvaluateHealthChecks()

	if err == nil {
		log.Error("Healthchecks did not throw an error, but should have")
		t.Fail()
	}
	if len(results) != 2 {
		log.Error("Healthcheck results had the wrong length")
		t.Fail()
	}
}

func TestEvaluatingFatalChecks(t *testing.T) {
	cxn := database.GetPGConnection(conf.DBURI())
	healthChecks, _ := ReadHealthCheckYAMLFromFile("healthchecksFatal.yml")
	healthChecks.RunHealthChecks(cxn)
	results, err := healthChecks.EvaluateHealthChecks()

	if err == nil {
		log.Error("Healthchecks did not throw an error, but should have")
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
	if len(results) != 5 {
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
		true,
		"t",
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

	rePass = SQLHealthCheck{"true", "select (select count(1) from information_schema.tables) > 0;", "basic test", "FATAL", true, "t"}
	reFail = SQLHealthCheck{"true", "select (select count(1) from information_schema.tables) < 0;", "basic test", "FATAL", false, "f"}
	prr = report.NewPongo2ReportRunnerFromString(TemplateHealthcheck)
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
