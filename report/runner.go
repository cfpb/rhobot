package report

import (
	"bytes"
	"encoding/json"
	"io"
	"strings"

	log "github.com/Sirupsen/logrus"
	"github.com/flosch/pongo2"
)

func init() {
	pongo2.RegisterFilter("addquote", filterAddquote)
}

// JSONReportRunner initilization should contain any variables used for report
type JSONReportRunner struct{}

// ReportReader Implementation for JSONReportRunner
func (jrr JSONReportRunner) ReportReader(reportSet Set) (io.Reader, error) {
	reportJSON, err := json.MarshalIndent(reportSet.GetReportMap(), "", "    ")
	if err != nil {
		log.Fatal(err)
	}
	r := bytes.NewReader(reportJSON)
	log.Debug(string(reportJSON))
	return r, err
}

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

// ReportReader Implementation for Pongo2ReportRunner
func (p2rr Pongo2ReportRunner) ReportReader(reportSet Set) (io.Reader, error) {
	reportBytes, err := p2rr.Template.ExecuteBytes(reportSet.GetReportMap())
	if err != nil {
		log.Fatal(err)
	}
	r := bytes.NewReader(reportBytes)
	log.Debug(string(reportBytes))
	return r, err
}

// filterAddquote pongo2 filter for adding an extra quote
func filterAddquote(in *pongo2.Value, param *pongo2.Value) (*pongo2.Value, *pongo2.Error) {
	output := strings.Replace(in.String(), "'", "''", -1)
	return pongo2.AsValue(output), nil
}
