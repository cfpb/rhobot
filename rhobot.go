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

	app.Commands = []cli.Command{
		{
			Name:    "run",
			Aliases: []string{},
			Usage:   "healthchecks|pipeline|tbd",
			Subcommands: []cli.Command{
				{
					Name:  "healthchecks",
					Usage: "HEALTHCHECK_FILE [DATABASE_URI]",
					Action: func(c *cli.Context) {

						// The path argument is required, but URI is optional
						if len(c.Args()) > 1 {
							config.SetDBURI(c.Args()[1])
						}
						fmt.Println("DB_URI: ", config.DBURI())

						var path string
						if len(c.Args()) > 0 {
							path = c.Args()[0]
						} else {
							fmt.Println("You must provide the path to the healthcheck file.")
						}
						fmt.Println("PATH: ", path)

						healthChecks := healthcheck.ReadYamlFromFile(path)
						cxn := database.GetPGConnection(config.DBURI())
						results, _ := healthcheck.PreformHealthChecks(healthChecks, cxn)
						footerString := healthcheck.FooterHealthcheck
						metadata := map[string]interface{}{
							"name":      healthChecks.Name,
							"db_name":   config.DBURI(),
							"footer":    footerString,
							"timestamp": time.Now().UTC().String(),
						}

						var elements []report.Element
						for _, val := range results {
							elements = append(elements, val)
						}

						prr := report.NewPongo2ReportRunnerFromString(healthcheck.TemplateHealthcheck)
						fhr := report.FileHandler{Filename: "healthcheckResult.html"}
						rs := report.Set{Elements: elements, Metadata: metadata}
						reader, _ := prr.ReportReader(rs)
						_ = fhr.HandleReport(reader)

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
