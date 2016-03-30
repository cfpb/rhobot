package report

import (
	"io/ioutil"

	log "github.com/Sirupsen/logrus"

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

func unmarshalDistributionFormat(yamldata []byte) DistributionFormat {

	var data DistributionFormat
	err := yaml.Unmarshal(yamldata, &data)
	if err != nil {
		log.Fatal(err)
	}
	return data
}

// ReadDistributionFormatYamlFromFile loads DistributionFormat data from a YAML file
func ReadDistributionFormatYamlFromFile(path string) DistributionFormat {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatal(err)
	}
	return unmarshalDistributionFormat(data)
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
