package report

import (
	"io"
	"strings"

	log "github.com/Sirupsen/logrus"
)

// Element interface for anything that is contained in a report
type Element interface {
	GetHeaders() []string
	GetValue(key string) string
}

// Runner interface for anything able to generate a report
type Runner interface {
	ReportReader(Set) (io.Reader, error)
}

// Set structure for containing elements and metadata for a report
type Set struct {
	Elements []Element
	Metadata map[string]interface{}
}

// GetReportMap converts a ReportSet to generic go map
func (rs *Set) GetReportMap() map[string]interface{} {
	elements := make([]map[string]interface{}, len(rs.Elements))
	for i, element := range rs.Elements {
		elementMap := make(map[string]interface{})
		for _, header := range element.GetHeaders() {
			elementMap[header] = element.GetValue(header)
		}
		elements[i] = elementMap
	}

	reportSetMap := make(map[string]interface{})
	reportSetMap["elements"] = elements
	reportSetMap["metadata"] = rs.Metadata
	return reportSetMap
}

// GetElementArray getter for Elements
func (rs *Set) GetElementArray() []Element {
	return rs.Elements
}

// GetMetadata getter for Metadata
func (rs *Set) GetMetadata() map[string]interface{} {
	return rs.Metadata
}

// Handler interface for anything able to consume a report
type Handler interface {
	HandleReport(io.Reader) error
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
