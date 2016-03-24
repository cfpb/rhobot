package main

import (
	"fmt"
	"os"
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

	app.Commands = []cli.Command{
		{
			Name:    "run",
			Aliases: []string{},
			Usage:   "healthchecks|pipeline|tbd",
			Subcommands: []cli.Command{
				{
					Name:  "healthchecks",
					Usage: "HEALTHCHECK_FILE [--dburi DATABASE_URI] [--report REPORT_FILE]",
					Flags: []cli.Flag{
						reportFileFlag,
						dburiFlag,
					},
					Action: func(c *cli.Context) {

						// variables to be populated by cli args
						var healthcheckPath string
						var reportPath string

						if c.String("dburi") != "" {
							config.SetDBURI(c.String("dburi"))
						}

						if c.Args().Get(0) != "" {
							healthcheckPath = c.Args().Get(0)
						} else {
							fmt.Println("You must provide the path to the healthcheck file.")
						}

						if c.String("report") != "" {
							reportPath = c.String("report")
						}

						fmt.Println("DB_URI: ", config.DBURI())
						fmt.Println("PATH: ", healthcheckPath)

						healthChecks := healthcheck.ReadYamlFromFile(healthcheckPath)
						cxn := database.GetPGConnection(config.DBURI())
						results, _ := healthcheck.PreformHealthChecks(healthChecks, cxn)
						metadata := map[string]interface{}{
							"name":      healthChecks.Name,
							"db_name":   config.PgDatabase,
							"footer":    healthcheck.FooterHealthcheck,
							"timestamp": time.Now().UTC().String(),
						}

						var elements []report.Element
						for _, val := range results {
							elements = append(elements, val)
						}

						if reportPath != "" {
							prr := report.NewPongo2ReportRunnerFromString(healthcheck.TemplateHealthcheck)
							fhr := report.FileHandler{Filename: reportPath}
							rs := report.Set{Elements: elements, Metadata: metadata}
							reader, _ := prr.ReportReader(rs)
							_ = fhr.HandleReport(reader)
						}

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
