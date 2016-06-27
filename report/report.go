package report

import "io"

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
