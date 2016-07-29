package report

import (
	"bufio"
	"crypto/tls"
	"database/sql"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strconv"

	log "github.com/Sirupsen/logrus"
	"gopkg.in/gomail.v2"
)

// Handler interface for anything able to consume a report
type Handler interface {
	HandleReport(io.Reader) error
}

// PrintHandler initilization should contain any variables used for report
type PrintHandler struct{}

// HandleReport consumes ReportReader output, prints to stdout
func (pr PrintHandler) HandleReport(reader io.Reader) (err error) {
	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		fmt.Printf("%s\n", scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		log.Error(err)
	}

	return err
}

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

// PGHandler initilization with sql connection
type PGHandler struct {
	Cxn *sql.DB
}

// HandleReport consumes ReportReader output, writes to postgres db
func (pg PGHandler) HandleReport(reader io.Reader) (err error) {

	reportBytes, err := ioutil.ReadAll(reader)
	if err != nil {
		log.Error(err)
		return err
	}
	reportString := string(reportBytes)
	result, err := pg.Cxn.Exec(reportString)
	if err != nil {
		log.Error("query failed: ", err)
	} else {
		rows, _ := result.RowsAffected()
		log.Info(rows, " Row(s) Affected")

	}

	return
}
