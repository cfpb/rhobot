package report

import (
	"bufio"
	"io"

	log "github.com/Sirupsen/logrus"
)

// PrintHandler initilization should contain any variables used for report
type PrintHandler struct{}

// HandleReport consumes ReportReader output, prints to stdout
func (pr PrintHandler) HandleReport(reader io.Reader) error {
	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		log.Debugf("%s\n", scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		return err
	}
	return nil
}
