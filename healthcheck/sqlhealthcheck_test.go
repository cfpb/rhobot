package healthcheck

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"reflect"
	"testing"
	"time"

	"gopkg.in/yaml.v2"

	"github.com/cfpb/rhobot/database"
	"github.com/cfpb/rhobot/report"
)

var healthchecks []byte
var host string
var db string
var user string
var pass string
var uri string

func init() {

	host = os.Getenv("PGHOST")
	db = os.Getenv("PGDATABASE")
	user = os.Getenv("PGUSER")
	pass = os.Getenv("PGPASSWORD")
	uri = fmt.Sprintf("postgres://%s:%s@%s/%s", user, pass, host, db)

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

	cxn := database.GetPGConnection(uri)
	healthChecks := unmarshalHealthChecks(healthchecks)
	RunHealthChecks(healthChecks, cxn)

}

func TestEvaluatingBasicChecks(t *testing.T) {

	cxn := database.GetPGConnection(uri)
	healthChecks := unmarshalHealthChecks(healthchecks)
	healthChecks = RunHealthChecks(healthChecks, cxn)
	results, err := EvaluateHealthChecks(healthChecks)

	if err != nil {
		t.Error("healthchecksTest threw an error")
	}
	if len(results) != 3 {
		t.Error("healthchecks results wrong length")
	}
}

func TestEvaluatingErrorsChecks(t *testing.T) {

	cxn := database.GetPGConnection(uri)
	healthChecks := ReadYamlFromFile("healthchecksErrors.yml")
	healthChecks = RunHealthChecks(healthChecks, cxn)
	results, err := EvaluateHealthChecks(healthChecks)

	if err == nil {
		t.Error("healthchecksErrors did not throw an error")
	}
	if len(results) != 2 {
		t.Error("healthchecks results wrong length")
	}
}

func TestEvaluatingFatalChecks(t *testing.T) {

	cxn := database.GetPGConnection(uri)
	healthChecks := ReadYamlFromFile("healthchecksFatal.yml")
	healthChecks = RunHealthChecks(healthChecks, cxn)
	results, err := EvaluateHealthChecks(healthChecks)

	if err == nil {
		t.Error("healthchecksFatal did not throw an error")
	}
	if len(results) != 1 {
		t.Error("healthchecks results wrong length")
	}
}

func TestPreformAllChecks(t *testing.T) {

	cxn := database.GetPGConnection(uri)
	healthChecks := ReadYamlFromFile("healthchecksAll.yml")
	results, err := PreformHealthChecks(healthChecks, cxn)

	if err == nil {
		t.Error("healthchecksAll did not throw an error")
	}
	if len(results) != 5 {
		t.Error("healthchecks results wrong length")
	}
}

func TestSQLHealthCheckReport(t *testing.T) {
    var hcr report.ReportableElement

    hcr = SQLHealthCheck{"true","select (select count(1) from information_schema.tables) > 0;","basic test","FATAL",true,"t"}
    for _, header := range hcr.GetHeaders(){
        fmt.Printf("%s : %s\n", header, hcr.GetValue(header))
    }

    if hcr.GetHeaders() == nil {
		t.Error("no headers in report ReportableElement")
	}

}

func TestHealthcheckPongo2Report(t *testing.T) {
	fmt.Println("TestHealthcheckPongo2Report")
	var rePass, reFail report.ReportableElement
	var rs report.ReportSet
	var prr report.ReportRunner

	rePass = SQLHealthCheck{"true","select (select count(1) from information_schema.tables) > 0;","basic test","FATAL",true,"t"}
	reFail = SQLHealthCheck{"true","select (select count(1) from information_schema.tables) < 0;","basic test","FATAL",false,"f"}
	prr = report.Pongo2ReportRunner{"./TemplateHealthcheck.html"}

	elements := []report.ReportableElement{rePass, reFail}
	metadata := map[string]interface{}{
	    "name":"TestHealthcheckPongo2Report",
        "db_name":"testdb",
        "footer":"Footer",
        "timestamp":time.Now().UTC().String(),
    }
	rs = report.ReportSet{elements, metadata}

	err := prr.WriteReport(rs)

	if err != nil {
		t.Fatalf("error writing report\n%s", err)
	}

}
