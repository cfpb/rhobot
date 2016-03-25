package report

import (
	"bytes"
	"encoding/json"
	"io"

	log "github.com/Sirupsen/logrus"
)

// JSONReportRunner initilization should contain any variables used for report
type JSONReportRunner struct{}

// ReportReader Implementation for ReportRunner
func (jrr JSONReportRunner) ReportReader(reportSet Set) (io.Reader, error) {
	reportJSON, err := json.MarshalIndent(reportSet.GetReportMap(), "", "    ")
	if err != nil {
		log.Fatal(err)
	}
	r := bytes.NewReader(reportJSON)
	log.Debug(string(reportJSON))
	return r, err
}
