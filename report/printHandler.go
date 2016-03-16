package report

import (
	"bufio"
	"fmt"
	"io"
)

//PrintHandler initilization should contain any variables used for report
type PrintHandler struct {
}

//HandleReport consumes ReportReader output, prints to stdout
func (pr PrintHandler) HandleReport(reader io.Reader) error {
	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		fmt.Printf("%s\n", scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		//print error to logger
		//fmt.Fprintln(os.Stderr, "There was an error with the scanner in attached container", err)
	}
	return nil
}
