package main

import (
	"fmt"
	"os"

	"github.com/cfpb/rhobot/database"
	"github.com/cfpb/rhobot/healthcheck"
	"github.com/codegangsta/cli"
)

func main() {

	app := cli.NewApp()
	app.Name = "Rhobot"
	app.Usage = "Rhobot is your friend."
	app.EnableBashCompletion = true

	app.Commands = []cli.Command{
		{
			Name:    "run",
			Aliases: []string{},
			Usage:   "healthchecks|tbd",
			Subcommands: []cli.Command{
				{
					Name:  "healthchecks",
					Usage: "[database uri] [path to healthcheck file]",
					Action: func(c *cli.Context) {
						dburi := c.Args()[0]
						path := c.Args()[1]
						fmt.Println("DB_URI: ", dburi)
						fmt.Println("PATH: ", path)

						healthChecks := healthcheck.ReadYamlFromFile(path)
						cxn := database.GetPGConnection(dburi)

						healthcheck.RunHealthChecks(healthChecks, cxn)

					},
				},
			},
		},
	}

	app.Run(os.Args)
}
