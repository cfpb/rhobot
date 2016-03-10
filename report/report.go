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
	elements []ReportableElement
	metadata map[string]interface{}
}

func (rs *ReportSet) GetReportMap() map[string]interface{} {

	elements := make([]map[string]interface{}, len(rs.elements))
	for i, element := range rs.elements {
		elementMap := make(map[string]interface{})
		for _, header := range element.GetHeaders() {
			elementMap[header] = element.GetValue(header)
		}
		elements[i] = elementMap
	}

	reportSetMap := make(map[string]interface{})
	reportSetMap["elements"] = elements
	reportSetMap["metadata"] = rs.metadata
	return reportSetMap
}

func (rs *ReportSet) GetElementArray() []ReportableElement {
	return rs.elements
}

func (rs *ReportSet) GetMetadata() map[string]interface{} {
	return rs.metadata
}
