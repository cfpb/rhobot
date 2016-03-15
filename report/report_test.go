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
	var rs Set
	var jrr Runner

	re = SimpleRE{[]string{"Some", "Thing"}}
	jrr = JSONReportRunner{}

	elements := []ReportableElement{re, re}
	metadata := map[string]interface{}{"test": "json"}
	rs = Set{elements, metadata}

	reader, err := jrr.ReportReader(rs)
	PrintReport(reader)
	if err != nil {
		t.Fatalf("error writing report\n%s", err)
	}

}

func TestPongo2Report(t *testing.T) {
	fmt.Println("TestPongo2Report")
	var re ReportableElement
	var rs Set
	var prr Runner

	re = SimpleRE{[]string{"Some", "Thing"}}
	prr = Pongo2ReportRunner{"./TemplateSimple.html"}

	elements := []ReportableElement{re, re}
	metadata := map[string]interface{}{"test": "pongo2"}
	rs = Set{elements, metadata}

	reader, err := prr.ReportReader(rs)
	PrintReport(reader)
	if err != nil {
		t.Fatalf("error writing report\n%s", err)
	}

}
