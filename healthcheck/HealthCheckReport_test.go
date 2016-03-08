package healthcheck

import (
	"testing"
	"fmt"
	"github.com/cfpb/rhobot/report"
)

func TestQuickHCReport(t *testing.T) {
    var hcr report.ReportableElement

    hcr = HealthCheckReport{SQLHealthCheck{"true","select (select count(1) from information_schema.tables) > 0;","basic test","FATAL",true,"t"}}
    for _, header := range hcr.GetHeaders(){
        fmt.Printf("%s : %s\n", header, hcr.GetValue(header))
    }

    if hcr.GetHeaders() == nil {
		t.Error("no headers in report ReportableElement")
	}

}
