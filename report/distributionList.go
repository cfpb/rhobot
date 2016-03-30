package report

import (
	"io/ioutil"

	"github.com/davecgh/go-spew/spew"
	"gopkg.in/yaml.v2"
)

// DistributionFormat is for unmarshiling a email distributionList file
type DistributionFormat struct {
	Severity SeverityDistribution `yaml:"severity"`
}

// SeverityDistribution list of emails ordered by severity
type SeverityDistribution struct {
	Debug []string `yaml:"debug,omitempty"`
	Info  []string `yaml:"info,omitempty"`
	Warn  []string `yaml:"warn,omitempty"`
	Error []string `yaml:"error,omitempty"`
	Fatal []string `yaml:"fatal,omitempty"`
}

// ReadDistributionFormatYAMLFromFile loads DistributionFormat data from a YAML file
func ReadDistributionFormatYAMLFromFile(path string) (format DistributionFormat, err error) {
	if data, err := ioutil.ReadFile(path); err == nil {
		err = yaml.Unmarshal(data, &format)
	}
	return
}

// Print the DistributionFormat
func (df DistributionFormat) Print() {
	spew.Dump(df)
}

// GetEmails returns list of emails based on log level
func (df DistributionFormat) GetEmails(level string) []string {
	switch level {
	case LogLevelArray[0]:
		return df.Severity.Debug
	case LogLevelArray[1]:
		return df.Severity.Info
	case LogLevelArray[2]:
		return df.Severity.Warn
	case LogLevelArray[3]:
		return df.Severity.Error
	case LogLevelArray[4]:
		return df.Severity.Fatal
	}
	return nil
}
