package report

import (
	"crypto/tls"
	"fmt"
	"io"
	"io/ioutil"

	"gopkg.in/gomail.v2"
)

// EmailHandler initilization should contain any variables used for report
type EmailHandler struct {
	SMTPHost    string
	SMTPPort    int
	SenderEmail string
	SenderName  string
	Recipients  []string
	Subject     string
	HTML        bool
}

// HandleReport consumes ReportReader output, writes to file
func (eh EmailHandler) HandleReport(reader io.Reader) error {

	msg := gomail.NewMessage()

	msg.SetHeaders(map[string][]string{
		"To":      eh.Recipients,
		"Subject": {eh.Subject},
	})

	if eh.SenderName != "" {
		msg.SetHeader("From", msg.FormatAddress(eh.SenderEmail, eh.SenderName))
	} else {
		msg.SetHeader("From", eh.SenderEmail)
	}

	var bodyType string
	if eh.HTML {
		bodyType = "text/html"
	} else {
		bodyType = "text/plain"
	}

	reportBytes, err := ioutil.ReadAll(reader)
	if err != nil {
		fmt.Println("print Reading error to logger")
		// TODO: print Reading error to logger
		return err
	}
	reportString := string(reportBytes)
	msg.SetBody(bodyType, reportString)

	//dialer := gomail.NewDialer(eh.SMTPHost, eh.SMTPPort, "", "")
	dialer := gomail.Dialer{Host: eh.SMTPHost, Port: eh.SMTPPort}
	dialer.TLSConfig = &tls.Config{InsecureSkipVerify: true}
	if err := dialer.DialAndSend(msg); err != nil {
		fmt.Println("print Dial error to logger")
		// TODO: print Dial error to logger
		return err
	}

	return err
}
