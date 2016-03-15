package report

import (
    "bytes"
	"github.com/flosch/pongo2"
	"io"
)

type Pongo2ReportRunner struct {
	TemplateFilePath string
}

func (p2rr Pongo2ReportRunner) ReportReader(reportSet ReportSet) (io.Reader,error)  {

	var tplExample = pongo2.Must(pongo2.FromFile(p2rr.TemplateFilePath))
	reportBytes, err := tplExample.ExecuteBytes(reportSet.GetReportMap())
	if err != nil {
		panic(err)
	}
	r := bytes.NewReader(reportBytes)
	//fmt.Println(string(reportBytes))
	return r, err
}
