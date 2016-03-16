package report

import (
	"bytes"
	"encoding/json"
	"io"
)

// JSONReportRunner initilization should contain any variables used for report
type JSONReportRunner struct {
}

// ReportReader Implementation for ReportRunner
func (jrr JSONReportRunner) ReportReader(reportSet Set) (io.Reader, error) {

	reportJSON, err := json.MarshalIndent(reportSet.GetReportMap(), "", "    ")
	if err != nil {
		panic(err)
	}
	r := bytes.NewReader(reportJSON)
	return r, err
}
