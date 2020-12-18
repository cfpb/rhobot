package healthcheck

// SQLHealthCheck is a data type for storing the definition
// and results of a SQL based health check
type SQLHealthCheck struct {
	Expected  string `yaml:"expected"`
	Query     string `yaml:"query"`
	Title     string `yaml:"title"`
	Severity  string `yaml:"severity"`
	Operation string `yaml:"operation,omitempty"`
	Passed    bool
	Actual    string
	Equal     bool
}

// Format is for unmarshiling a healthcheck file
// and contains control information for a set of SQLHealthChecks
type Format struct {
	Name         string           `yaml:"name"`
	Distribution []string         `yaml:"distribution"`
	Tests        []SQLHealthCheck `yaml:"tests"`
}

// HCError is a error helper for knowing to exit early on a failed healthcheck
type HCError struct {
	Err  string
	Exit bool
}
