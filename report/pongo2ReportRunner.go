package report

import (
	"bytes"
	"io"

	"github.com/flosch/pongo2"
)

//Pongo2ReportRunner initilization with template path
type Pongo2ReportRunner struct {
	TemplateFilePath string
}

//ReportReader Implementation for ReportRunner
func (p2rr Pongo2ReportRunner) ReportReader(reportSet Set) (io.Reader, error) {

	var tplExample = pongo2.Must(pongo2.FromFile(p2rr.TemplateFilePath))
	reportBytes, err := tplExample.ExecuteBytes(reportSet.GetReportMap())
	if err != nil {
		panic(err)
	}
	r := bytes.NewReader(reportBytes)
	//fmt.Println(string(reportBytes))
	return r, err
}
