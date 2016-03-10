package main

import (
	"fmt"
	"os"

	"github.com/cfpb/rhobot/config"
	"github.com/cfpb/rhobot/database"
	"github.com/cfpb/rhobot/gocd"
	"github.com/cfpb/rhobot/healthcheck"
	"github.com/codegangsta/cli"
)

func main() {

	app := cli.NewApp()
	app.Name = "Rhobot"
	app.Usage = "Rhobot is a database development tool that uses DevOps best practices."
	app.EnableBashCompletion = true

	config := config.NewConfig()

	app.Commands = []cli.Command{
		{
			Name:    "run",
			Aliases: []string{},
			Usage:   "healthchecks|gocd|tbd",
			Subcommands: []cli.Command{
				{
					Name:  "healthchecks",
					Usage: "HEALTHCHECK_FILE [DATABASE_URI]",
					Action: func(c *cli.Context) {

						// The path argument is required, but URI is optional
						if len(c.Args()) > 1 {
							config.SetDBURI(c.Args()[1])
						}

						var path string
						if len(c.Args()) > 0 {
							path = c.Args()[0]
						} else {
							fmt.Println("You must provide the path to the healthcheck file.")
						}

						fmt.Println("DB_URI: ", config.GetDBURI())
						fmt.Println("PATH: ", path)

						healthChecks := healthcheck.ReadYamlFromFile(path)
						cxn := database.GetPGConnection(config.GetDBURI())
						healthcheck.PreformHealthChecks(healthChecks, cxn)

						//TODO: turn results into report
					},
				},
				{
					Name:    "pipeline",
					Aliases: []string{},
					Usage:   "Interact with GoCD pipeline",
					Flags: []cli.Flag{
						cli.StringFlag{
							Name:  "host, h",
							Value: "",
							Usage: "host of the GoCD server",
						},
					},
					Subcommands: []cli.Command{
						{
							Name:  "push",
							Usage: "PATH [PIPELINE_GROUP]",
							Action: func(c *cli.Context) {
								if c.String("host") != "" {
									config.SetGOCDHost(c.String("host"))
								}

								path := c.Args()[0]
								group := c.Args().Get(1)

								gocd.Push(config.GetGOCDHost(), path, group)
							},
						},
						{
							Name:  "pull",
							Usage: "PATH",
							Action: func(c *cli.Context) {
								if c.String("host") != "" {
									config.SetGOCDHost(c.String("host"))
								}

								path := c.Args()[0]

								gocd.Pull(config.GetGOCDHost(), path)
							},
						},
						{
							Name:  "clone",
							Usage: "PIPELINE_NAME PATH",
							Action: func(c *cli.Context) {
								if c.String("host") != "" {
									config.SetGOCDHost(c.String("host"))
								}

								name := c.Args()[0]
								path := c.Args()[1]

								gocd.Clone(config.GetGOCDHost(), path, name)
							},
						},
					},
				},
			},
		},
	}

	app.Run(os.Args)
}
