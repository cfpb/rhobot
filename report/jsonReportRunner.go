package report

import (
	"encoding/json"
	"io/ioutil"
)

type JSONReportRunner struct {
	OutputFilePath string
}

func (jrr JSONReportRunner) WriteReport(reportSet ReportSet) error {

	reportJSON, err := json.MarshalIndent(reportSet.GetReportMap(), "", "    ")
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(jrr.OutputFilePath, reportJSON, 0666)
	return err
}
