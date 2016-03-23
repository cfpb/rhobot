package report

import (
	"bufio"
	"io"
	"net/smtp"
)

// EmailHandler initilization should contain any variables used for report
type EmailHandler struct {
	StmpHost  string
	Sender    string
	Recipient string
	Html bool
}

// HandleReport consumes ReportReader output, writes to file
func (eh EmailHandler) HandleReport(reader io.Reader) error {

	client, err := smtp.Dial(eh.StmpHost)
	if err != nil {
		// TODO: print Dial error to logger
		return err
	}
	defer client.Close()

	client.Mail(eh.Sender)
	client.Rcpt(eh.Recipient)
	wc, err := client.Data()
	if err != nil {
		// TODO: print Data Writer error to logger
		return err
	}
	defer wc.Close()

	writer := bufio.NewWriter(wc)
	scanner := bufio.NewScanner(reader)

	if(eh.Html){
	writer.WriteString("Content-Type:text/html\r\n")
	}

	for scanner.Scan() {
		_, err := writer.WriteString(scanner.Text() + "\n")
		if err != nil {
			// TODO: print Writing error to logger
			return err
		}
	}
	if err := scanner.Err(); err != nil {
		// TODO: print Scanning error to logger
		return err
	}

	writer.Flush()
	return err
}
