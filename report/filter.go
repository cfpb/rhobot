package report

import (
	"io/ioutil"
	"strings"

	log "github.com/Sirupsen/logrus"
	"github.com/davecgh/go-spew/spew"
	"gopkg.in/yaml.v2"
)

// LogLevelArray string denoting levels of logs
var LogLevelArray = []string{"Debug", "Info", "Warn", "Error", "Fatal"}

// LogLevelMap a map that holds the integer for a log level
var LogLevelMap = map[string]int{
	"debug": 0,
	"info":  1,
	"warn":  2,
	"error": 3,
	"fatal": 4,
}

// FilterReportSet by the logLevel
func FilterReportSet(rs Set, logLevel string) Set {

	elements := rs.GetElementArray()
	filteredElements := make([]Element, 0, 0)
	copy(filteredElements, elements)

	for _, elm := range elements {
		if logLevelIncludes(elm, logLevel) {
			filteredElements = append(filteredElements, elm)
		}
	}

	filteredSet := Set{Elements: filteredElements, Metadata: rs.Metadata}
	return filteredSet
}

// logLevelIncludes utility function to know if one loglevel includes another
func logLevelIncludes(elm Element, logLevel string) bool {

	elmSeverity := strings.ToLower(elm.GetValue("Severity"))

	if elmSeverity == "" {
		log.Warn("severity field not found in element")
		return false
	}

	elmIndex, ok := LogLevelMap[elmSeverity]
	if !ok {
		log.Error("severity level not found in element")
		return false
	}

	logSeverity := strings.ToLower(logLevel)
	logIndex, ok := LogLevelMap[logSeverity]
	if !ok {
		log.Error("severity level not found in argument")
		return false
	}

	log.Debugf("logLevelIncludes elem:%v , log:%v, bool:%v", elmSeverity, logSeverity, logIndex <= elmIndex)
	return logIndex <= elmIndex

}

// DistributionFormat is for unmarshiling a email distributionList file
type DistributionFormat struct {
	Severity struct {
		Debug []string `yaml:"debug,omitempty"`
		Info  []string `yaml:"info,omitempty"`
		Warn  []string `yaml:"warn,omitempty"`
		Error []string `yaml:"error,omitempty"`
		Fatal []string `yaml:"fatal,omitempty"`
	} `yaml:"severity"`
}

// ReadDistributionFormatYAMLFromFile loads DistributionFormat data from a YAML file
func ReadDistributionFormatYAMLFromFile(path string) (format DistributionFormat, err error) {
	data, err := ioutil.ReadFile(path)
	if err == nil {
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
