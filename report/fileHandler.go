package report

import (
	"bufio"
	"io"
	"os"

	log "github.com/Sirupsen/logrus"
)

// FileHandler initilization should contain any variables used for report
type FileHandler struct {
	Filename string
}

// HandleReport consumes ReportReader output, writes to file
func (fr FileHandler) HandleReport(reader io.Reader) (err error) {

	f, err := os.Create(fr.Filename)
	w := bufio.NewWriter(f)

	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		_, err := w.WriteString(scanner.Text() + "\n")
		if err != nil {
			log.Error(err)
		}
	}

	if err := scanner.Err(); err != nil {
		log.Error(err)
	}

	w.Flush()

	return err
}
