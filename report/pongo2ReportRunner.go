package report

import (
	"bytes"
	"io"

	"github.com/flosch/pongo2"
)

// NewPongo2ReportRunnerFromFile constructor with template file
func NewPongo2ReportRunnerFromFile(TemplateFilePath string) *Pongo2ReportRunner {
	var template = pongo2.Must(pongo2.FromFile(TemplateFilePath))
	return &Pongo2ReportRunner{
		Template: *template,
	}
}

// NewPongo2ReportRunnerFromString constructor with template string
func NewPongo2ReportRunnerFromString(TemplateString string) *Pongo2ReportRunner {
	var template = pongo2.Must(pongo2.FromString(TemplateString))
	return &Pongo2ReportRunner{
		Template: *template,
	}
}

// Pongo2ReportRunner initilization with template object
type Pongo2ReportRunner struct {
	Template pongo2.Template
}

// ReportReader Implementation for ReportRunner
func (p2rr Pongo2ReportRunner) ReportReader(reportSet Set) (io.Reader, error) {

	reportBytes, err := p2rr.Template.ExecuteBytes(reportSet.GetReportMap())
	if err != nil {
		panic(err)
	}
	r := bytes.NewReader(reportBytes)
	// fmt.Println(string(reportBytes))
	return r, err
}
