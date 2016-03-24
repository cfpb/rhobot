package main

import (
	"fmt"
	"os"
	"strconv"
	"time"

	log "github.com/Sirupsen/logrus"
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

	logLevelFlag := cli.StringFlag{
		Name:  "loglevel, lvl",
		Value: "",
		Usage: "sets the log level for Rhobot",
	}
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
						logLevelFlag,
						reportFileFlag,
						dburiFlag,
						emailListFlag,
					},
					Action: func(c *cli.Context) {
						if c.String("loglevel") != "" {
							config.SetLogLevel(c.String("loglevel"))
						}

						// variables to be populated by cli args
						var healthcheckPath string
						var reportPath string
						var emailListPath string

						if c.Args().Get(0) != "" {
							healthcheckPath = c.Args().Get(0)
						} else {
							log.Error("You must provide the path to the healthcheck file.")
							return
						}

						log.Info("Running health checks from ", healthcheckPath)

						if c.String("dburi") != "" {
							config.SetDBURI(c.String("dburi"))
						}

						log.Debug("DB_URI: ", config.DBURI())

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
								logLevelFlag,
								gocdHostFlag,
							},
							Action: func(c *cli.Context) {
								if c.String("loglevel") != "" {
									config.SetLogLevel(c.String("loglevel"))
								}

								if c.String("host") != "" {
									log.Debug("Setting GoCD host: ", c.String("host"))
									config.SetGoCDHost(c.String("host"))
								}

								if len(c.Args()) > 0 {
									path := c.Args()[0]
									group := c.Args().Get(1)
									log.Infof("Pushing config from %v to pipeline group %v...", path, group)
									if err := gocd.Push(config.GoCDURL(), path, group); err != nil {
										log.Error(err)
										log.Fatal("Failed to push pipeline configuration!")
									}
									log.Info("Success!")
								} else {
									log.Fatal("PATH is required for the 'push' command.")
								}
							},
						},
						{
							Name:  "pull",
							Usage: "PATH",
							Flags: []cli.Flag{
								logLevelFlag,
								gocdHostFlag,
							},
							Action: func(c *cli.Context) {
								if c.String("loglevel") != "" {
									config.SetLogLevel(c.String("loglevel"))
								}

								if c.String("host") != "" {
									config.SetGoCDHost(c.String("host"))
								}

								if len(c.Args()) > 0 {
									path := c.Args()[0]
									log.Infof("Pulling config from %v to %v...", config.GoCDURL(), path)
									if err := gocd.Pull(config.GoCDURL(), path); err != nil {
										log.Fatal("Failed to pull pipeline config: ", err)
									}
									log.Info("Success!")
								} else {
									log.Fatal("A path to pull the pipeline config to is required.")
								}
							},
						},
						{
							Name:  "clone",
							Usage: "PIPELINE_NAME PATH",
							Flags: []cli.Flag{
								logLevelFlag,
								gocdHostFlag,
							},
							Action: func(c *cli.Context) {
								if c.String("loglevel") != "" {
									config.SetLogLevel(c.String("loglevel"))
								}

								if c.String("host") != "" {
									config.SetGoCDHost(c.String("host"))
								}

								if len(c.Args()) > 1 {
									name := c.Args()[0]
									path := c.Args()[1]
									log.Infof("Cloning pipeline %v to %v...", name, path)
									if err := gocd.Clone(config.GoCDURL(), path, name); err != nil {
										log.Fatal("Failed to clone pipeline config: ", err)
									}
									log.Info("Success!")
								} else {
									log.Fatal("A pipeline name and a path to clone to are required.")
								}
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
