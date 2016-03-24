package healthcheck

import (
	"bytes"
	"io"
	"os"
	"reflect"
	"testing"
	"time"

	"gopkg.in/yaml.v2"

	log "github.com/Sirupsen/logrus"
	"github.com/cfpb/rhobot/config"
	"github.com/cfpb/rhobot/database"
	"github.com/cfpb/rhobot/report"
)

var healthchecks []byte
var conf *config.Config

func init() {
	conf = config.NewConfig()

	buf := bytes.NewBuffer(nil)
	f, _ := os.Open("healthchecksTest.yml")
	io.Copy(buf, f)
	f.Close()
	healthchecks = buf.Bytes()
}

func TestUnmarshal(t *testing.T) {
	unmarshalHealthChecks(healthchecks)
}

// TestUnmarshalFidelityLoss checks that data can be reserielized without fidelity loss
func TestUnmarshalFidelityLoss(t *testing.T) {
	data := unmarshalHealthChecks(healthchecks)
	healthchecks2, _ := yaml.Marshal(data)
	data2 := unmarshalHealthChecks(healthchecks2)
	if !reflect.DeepEqual(data, data2) {
		t.Error("not the same")
	}
}

func TestRunningBasicChecks(t *testing.T) {
	cxn := database.GetPGConnection(conf.DBURI())
	healthChecks := unmarshalHealthChecks(healthchecks)
	RunHealthChecks(healthChecks, cxn)
}

func TestEvaluatingBasicChecks(t *testing.T) {
	cxn := database.GetPGConnection(conf.DBURI())
	healthChecks := unmarshalHealthChecks(healthchecks)
	healthChecks = RunHealthChecks(healthChecks, cxn)
	results, err := EvaluateHealthChecks(healthChecks)

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
	healthChecks := ReadYamlFromFile("healthchecksErrors.yml")
	healthChecks = RunHealthChecks(healthChecks, cxn)
	results, err := EvaluateHealthChecks(healthChecks)

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
	healthChecks := ReadYamlFromFile("healthchecksFatal.yml")
	healthChecks = RunHealthChecks(healthChecks, cxn)
	results, err := EvaluateHealthChecks(healthChecks)

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
	healthChecks := ReadYamlFromFile("healthchecksAll.yml")
	results, err := PreformHealthChecks(healthChecks, cxn)

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
