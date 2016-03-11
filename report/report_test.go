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

func TestJSONHCReport(t *testing.T) {
	fmt.Println("TestJSONHCReport")
	var re ReportableElement
	var jrr ReportRunner

	re = SimpleRE{[]string{"Some", "Thing"}}
	jrr = JSONReportRunner{"./something.json"}

	elements := []ReportableElement{re}
	err := jrr.WriteReport(elements)

	if err != nil {
		t.Fatalf("error writing report\n%s", err)
	}

}
