package report

import ()

type ReportableElement interface {
	GetHeaders() []string
	GetValue(key string) string
}

type ReportRunner interface {
	WriteReport(ReportSet) error
}

type ReportSet struct {
	Elements []ReportableElement
	Metadata map[string]interface{}
}

func (rs *ReportSet) GetReportMap() map[string]interface{} {

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

func (rs *ReportSet) GetElementArray() []ReportableElement {
	return rs.Elements
}

func (rs *ReportSet) GetMetadata() map[string]interface{} {
	return rs.Metadata
}
