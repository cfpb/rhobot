package healthcheck

type HealthCheckReport struct {
	SQLHealthCheck
}

var HealthCheckReportHeaders = []string{"Title", "Query", "Test Passed", "Expected", "Actual"}

func (hcr HealthCheckReport) GetHeaders() []string {
	return HealthCheckReportHeaders[0:]
}

func (hcr HealthCheckReport) GetValue(key string) string {
	//return key+"_val"

	switch key {
	case HealthCheckReportHeaders[0]:
		return hcr.Title
	case HealthCheckReportHeaders[1]:
		return hcr.Query
	case HealthCheckReportHeaders[2]:
		if hcr.Passed {
			return "Succeed"
		}
		return "Fail"
	case HealthCheckReportHeaders[3]:
		return hcr.Expected
	case HealthCheckReportHeaders[4]:
		return hcr.Actual
	}
	return ""
}
