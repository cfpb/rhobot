package report

import (
	"io"
	"io/ioutil"

	"gopkg.in/gomail.v2"
)

// EmailHandler initilization should contain any variables used for report
type EmailHandler struct {
	SMTPHost   string
	SMTPPort   int
	Sender     string
	Recipients []string
	Subject    string
	HTML       bool
}

// HandleReport consumes ReportReader output, writes to file
func (eh EmailHandler) HandleReport(reader io.Reader) error {

	msg := gomail.NewMessage()

	msg.SetHeaders(map[string][]string{
		"From":    {eh.Sender},
		"To":      eh.Recipients,
		"Subject": {eh.Subject},
	})

	reportBytes, err := ioutil.ReadAll(reader)
	if err != nil {
		// TODO: print Reading error to logger
		return err
	}
	reportString := string(reportBytes)

	var bodyType string
	if eh.HTML {
		bodyType = "text/html"
	} else {
		bodyType = "text/plain"
	}

	msg.SetBody(bodyType, reportString)

	//dialer := gomail.NewDialer(eh.SMTPHost, eh.SMTPPort, "", "")
	dialer := gomail.Dialer{Host: eh.SMTPHost, Port: eh.SMTPPort}
	if err := dialer.DialAndSend(msg); err != nil {
		// TODO: print Dial error to logger
		return err
	}

	return err
}
