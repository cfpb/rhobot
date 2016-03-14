package report

import (
	"fmt"
	"github.com/flosch/pongo2"
)

type Pongo2ReportRunner struct {
	TemplateFilePath string
}

func (p2rr Pongo2ReportRunner) WriteReport(reportSet ReportSet) error {

	var tplExample = pongo2.Must(pongo2.FromFile(p2rr.TemplateFilePath))
	out, err := tplExample.Execute(reportSet.GetReportMap())
	if err != nil {
		panic(err)
	}
	fmt.Println(out)
	return nil
}
