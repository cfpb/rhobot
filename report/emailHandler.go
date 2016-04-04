package report

import (
	"crypto/tls"
	"io"
	"io/ioutil"
	"strconv"

	log "github.com/Sirupsen/logrus"
	"gopkg.in/gomail.v2"
)

// EmailHandler initilization should contain any variables used for report
type EmailHandler struct {
	SMTPHost    string
	SMTPPort    string
	SenderEmail string
	SenderName  string
	Recipients  []string
	Subject     string
	HTML        bool
}

// HandleReport consumes ReportReader output, writes to file
func (eh EmailHandler) HandleReport(reader io.Reader) (err error) {

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
		log.Error(err)
		return err
	}
	reportString := string(reportBytes)
	msg.SetBody(bodyType, reportString)

	SMTPPortInt, _ := strconv.Atoi(eh.SMTPPort)
	dialer := gomail.Dialer{Host: eh.SMTPHost, Port: SMTPPortInt}
	dialer.TLSConfig = &tls.Config{InsecureSkipVerify: true}
	if err := dialer.DialAndSend(msg); err != nil {
		log.Error(err)
	}

	return err
}
