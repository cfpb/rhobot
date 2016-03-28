package report

import (
	"bufio"
	"io"

	log "github.com/Sirupsen/logrus"
)

// PrintHandler initilization should contain any variables used for report
type PrintHandler struct{}

// HandleReport consumes ReportReader output, prints to stdout
func (pr PrintHandler) HandleReport(reader io.Reader) (err error) {
	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		log.Infof("%s\n", scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		log.Error(err)
	}

	return err
}
