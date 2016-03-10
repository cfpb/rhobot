package report

import (
	"fmt"
	"testing"
)

type SimpleRE struct {
	SimpleHeaders []string
}

func (sre SimpleRE) GetHeaders() []string {
	return sre.SimpleHeaders[0:]
}

func (sre SimpleRE) GetValue(key string) string {
	return "simple"
}

func TestJSONReport(t *testing.T) {
	fmt.Println("TestJSONReport")
	var re ReportableElement
	var rs ReportSet
	var jrr ReportRunner

	re = SimpleRE{[]string{"Some", "Thing"}}
	jrr = JSONReportRunner{"./something.json"}

	elements := []ReportableElement{re, re}
	metadata := map[string]interface{}{"test": "json"}
	rs = ReportSet{elements, metadata}

	err := jrr.WriteReport(rs)

	if err != nil {
		t.Fatalf("error writing report\n%s", err)
	}

}

func TestPongo2Report(t *testing.T) {
	fmt.Println("TestPongo2Report")
	var re ReportableElement
	var rs ReportSet
	var prr ReportRunner

	re = SimpleRE{[]string{"Some", "Thing"}}
	prr = Pongo2ReportRunner{"./simpleTemplate.html"}

	elements := []ReportableElement{re, re}
	metadata := map[string]interface{}{"test": "pongo2"}
	rs = ReportSet{elements, metadata}

	err := prr.WriteReport(rs)

	if err != nil {
		t.Fatalf("error writing report\n%s", err)
	}

}
