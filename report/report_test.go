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
	var re Element
	var rs Set
	var jrr Runner
	var phr Handler

	re = SimpleRE{[]string{"Some", "Thing"}}
	jrr = JSONReportRunner{}
	phr = PrintHandler{}

	elements := []Element{re, re}
	metadata := map[string]interface{}{"test": "json"}
	rs = Set{elements, metadata}

	reader, err := jrr.ReportReader(rs)
	err = phr.HandleReport(reader)
	if err != nil {
		t.Fatalf("error writing report\n%s", err)
	}

}

func TestJSONReportToFile(t *testing.T) {
	fmt.Println("TestJSONReportToFile")
	var re Element
	var rs Set
	var jrr Runner
	var fhr Handler

	re = SimpleRE{[]string{"Some", "Thing"}}
	jrr = JSONReportRunner{}
	fhr = FileHandler{"something.json"}

	elements := []Element{re, re}
	metadata := map[string]interface{}{"test": "json"}
	rs = Set{elements, metadata}

	reader, err := jrr.ReportReader(rs)
	err = fhr.HandleReport(reader)
	if err != nil {
		t.Fatalf("error writing report\n%s", err)
	}

}

func TestPongo2Report(t *testing.T) {
	fmt.Println("TestPongo2Report")
	var re Element
	var rs Set
	var prr Runner
	var phr Handler

	re = SimpleRE{[]string{"Some", "Thing"}}
	prr = NewPongo2ReportRunnerFromString(TemplateSimple)
	phr = PrintHandler{}

	elements := []Element{re, re}
	metadata := map[string]interface{}{"test": "pongo2"}
	rs = Set{elements, metadata}

	reader, err := prr.ReportReader(rs)
	err = phr.HandleReport(reader)
	if err != nil {
		t.Fatalf("error writing report\n%s", err)
	}

}
