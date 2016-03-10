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
	fmt.Println("--------------")
	fmt.Println(reportSet.GetReportMap())

	// Output: Hello Florian!
	// reportSet := make([]map[string]interface{}, len(elements))
	// for i, element := range elements {
	// 	elementMap := make(map[string]interface{})
	// 	for _, header := range element.GetHeaders() {
	// 		elementMap[header] = element.GetValue(header)
	// 	}
	// 	reportSet[i] = elementMap
	// }
	//
	// reportJSON, err := json.MarshalIndent(reportSet, "", "    ")
	// if err != nil {
	// 	return err
	// }
	// err = ioutil.WriteFile(jrr.OutputFilePath, reportJSON, 0666)
	// return err
	return nil
}
