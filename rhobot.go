package main

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/cfpb/rhobot/config"
	"github.com/cfpb/rhobot/database"
	"github.com/cfpb/rhobot/gocd"
	"github.com/cfpb/rhobot/healthcheck"
	"github.com/cfpb/rhobot/report"
	"github.com/codegangsta/cli"
)

func main() {

	app := cli.NewApp()
	app.Name = "Rhobot"
	app.Usage = "Rhobot is a database development tool that uses DevOps best practices."
	app.EnableBashCompletion = true

	config := config.NewConfig()

	gocdHostFlag := cli.StringFlag{
		Name:  "host",
		Value: "",
		Usage: "host of the GoCD server",
	}
	reportFileFlag := cli.StringFlag{
		Name:  "report",
		Value: "",
		Usage: "path to the healthcheck report",
	}
	dburiFlag := cli.StringFlag{
		Name:  "dburi",
		Value: "",
		Usage: "database uri postgres://user:password@host:port/database",
	}
	emailListFlag := cli.StringFlag{
		Name:  "email",
		Value: "",
		Usage: "yaml file containing email distribution list",
	}

	app.Commands = []cli.Command{
		{
			Name:    "run",
			Aliases: []string{},
			Usage:   "healthchecks|pipeline|tbd",
			Subcommands: []cli.Command{
				{
					Name:  "healthchecks",
					Usage: "HEALTHCHECK_FILE [--dburi DATABASE_URI] [--report REPORT_FILE] [--email DISTRIBUTION_FILE]",
					Flags: []cli.Flag{
						reportFileFlag,
						dburiFlag,
						emailListFlag,
					},
					Action: func(c *cli.Context) {

						// variables to be populated by cli args
						var healthcheckPath string
						var reportPath string
						var emailListPath string

						if c.Args().Get(0) != "" {
							healthcheckPath = c.Args().Get(0)
						} else {
							fmt.Println("You must provide the path to the healthcheck file.")
						}

						if c.String("dburi") != "" {
							config.SetDBURI(c.String("dburi"))
						}
						if c.String("report") != "" {
							reportPath = c.String("report")
						}
						if c.String("email") != "" {
							emailListPath = c.String("email")
						}

						healthcheckRunner(config, healthcheckPath, reportPath, emailListPath)
					},
				},
				{
					Name:    "pipeline",
					Aliases: []string{},
					Usage:   "Interact with GoCD pipeline",
					Subcommands: []cli.Command{
						{
							Name:  "push",
							Usage: "PATH [PIPELINE_GROUP]",
							Flags: []cli.Flag{
								gocdHostFlag,
							},
							Action: func(c *cli.Context) {
								if c.String("host") != "" {
									config.SetGoCDHost(c.String("host"))
								}

								if len(c.Args()) > 0 {
									path := c.Args()[0]
									group := c.Args().Get(1)
									gocd.Push(config.GoCDURL(), path, group)
								} else {
									fmt.Println("PATH is required for push command.")
								}
							},
						},
						{
							Name:  "pull",
							Usage: "PATH",
							Flags: []cli.Flag{
								gocdHostFlag,
							},
							Action: func(c *cli.Context) {
								if c.String("host") != "" {
									config.SetGoCDHost(c.String("host"))
								}

								path := c.Args()[0]

								gocd.Pull(config.GoCDURL(), path)
							},
						},
						{
							Name:  "clone",
							Usage: "PIPELINE_NAME PATH",
							Flags: []cli.Flag{
								gocdHostFlag,
							},
							Action: func(c *cli.Context) {
								if c.String("host") != "" {
									config.SetGoCDHost(c.String("host"))
								}

								name := c.Args()[0]
								path := c.Args()[1]

								gocd.Clone(config.GoCDURL(), path, name)
							},
						},
					},
				},
			},
		},
	}

	app.Run(os.Args)
}

func healthcheckRunner(config *config.Config, healthcheckPath string, reportPath string, emailListPath string) {
	fmt.Println("DB_URI: ", config.DBURI())
	fmt.Println("PATH: ", healthcheckPath)

	healthChecks := healthcheck.ReadYamlFromFile(healthcheckPath)
	cxn := database.GetPGConnection(config.DBURI())
	results, _ := healthcheck.PreformHealthChecks(healthChecks, cxn)
	var elements []report.Element
	for _, val := range results {
		elements = append(elements, val)
	}

	// Make Templated report
	metadata := map[string]interface{}{
		"name":      healthChecks.Name,
		"db_name":   config.PgDatabase,
		"footer":    healthcheck.FooterHealthcheck,
		"timestamp": time.Now().UTC().String(),
	}
	prr := report.NewPongo2ReportRunnerFromString(healthcheck.TemplateHealthcheck)
	rs := report.Set{Elements: elements, Metadata: metadata}
	reader, _ := prr.ReportReader(rs)

	// Write report to file
	if reportPath != "" {
		fhr := report.FileHandler{Filename: reportPath}
		_ = fhr.HandleReport(reader)
	}

	SMTPPortInt, _ := strconv.Atoi(config.SMTPPort)

	// Email report
	if emailListPath != "" {
		ehr := report.EmailHandler{
			SMTPHost:    config.SMTPHost,
			SMTPPort:    SMTPPortInt,
			SenderEmail: "-",
			SenderName:  "-",
			Subject:     "-",
			Recipients:  []string{"-"},
			HTML:        true,
		}
		_ = ehr.HandleReport(reader)
	}
}
