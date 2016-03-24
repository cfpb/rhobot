package report

import (
	"bufio"
	"io"
	"os"
)

// FileHandler initilization should contain any variables used for report
type FileHandler struct {
	Filename string
}

// HandleReport consumes ReportReader output, writes to file
func (fr FileHandler) HandleReport(reader io.Reader) error {

	f, err := os.Create(fr.Filename)
	w := bufio.NewWriter(f)

	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		_, err := w.WriteString(scanner.Text() + "\n")
		if err != nil {
			// TODO: print Writing error to logger
		}
	}
	if err := scanner.Err(); err != nil {
		// TODO: print Scanning error to logger
	}

	w.Flush()

	return err
}
