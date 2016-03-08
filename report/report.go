package report

import ()

type ReportableElement interface {
	GetHeaders() []string
	GetValue(key string) string
}
