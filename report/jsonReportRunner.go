package report

import (
	"encoding/json"
	"bytes"
	"io"
)

type JSONReportRunner struct {
}

func (jrr JSONReportRunner) ReportReader(reportSet ReportSet) (io.Reader,error) {

	reportJSON, err := json.MarshalIndent(reportSet.GetReportMap(), "", "    ")
	if err != nil {
		panic(err)
	}
	r := bytes.NewReader(reportJSON)
	//err = ioutil.WriteFile(jrr.OutputFilePath, reportJSON, 0666)
	return r, err
}
