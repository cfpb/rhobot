package main

import (
	"fmt"
	"os"

	"github.com/cfpb/rhobot/database"
	"github.com/cfpb/rhobot/gocd"
	"github.com/cfpb/rhobot/healthcheck"
	"github.com/codegangsta/cli"
)

func main() {

	app := cli.NewApp()
	app.Name = "Rhobot"
	app.Usage = "Rhobot is your friend."
	app.EnableBashCompletion = true
	// app.Action = func(c *cli.Context) {
	// 	fmt.Printf(healthcheck.Hello())
	// }

	app.Commands = []cli.Command{
		{
			Name:    "run",
			Aliases: []string{},
			Usage:   "healthchecks|gocd|tbd",
			Subcommands: []cli.Command{
				{
					Name:  "healthchecks",
					Usage: "[database uri] [path to healthcheck file]",
					Action: func(c *cli.Context) {
						dburi := c.Args()[0]
						path := c.Args()[1]
						fmt.Println("DB_URI: ", dburi)
						fmt.Println("PATH: ", path)
						fmt.Println(healthcheck.Hello())

						healthChecks := healthcheck.ReadYamlFromFile(path)
						cxn := database.GetPGConnection(dburi)

						healthcheck.RunHealthChecks(healthChecks, cxn)

					},
				},
				{
					Name:    "gocd",
					Aliases: []string{},
					Usage:   "Interact with GoCD pipeline",
					Subcommands: []cli.Command{
						{
							Name:  "push",
							Usage: "[gocd host] [path to pipeline config] [pipeline group]",
							Action: func(c *cli.Context) {
								gocdhost := c.Args()[0]
								path := c.Args()[1]
								group := c.Args()[2]

								gocd.Push(gocdhost, path, group)

							},
						},
						{
							Name:  "pull",
							Usage: "[gocd host] [path to pipeline config]",
							Action: func(c *cli.Context) {
								gocdhost := c.Args()[0]
								path := c.Args()[1]

								gocd.Pull(gocdhost, path)
							},
						},
						{
							Name:  "clone",
							Usage: "[gocd host] [path to pipeline config] [pipeline name]",
							Action: func(c *cli.Context) {
								gocdhost := c.Args()[0]
								path := c.Args()[1]
								name := c.Args()[2]

								gocd.Clone(gocdhost, path, name)
							},
						},
					},
				},
			},
		},
	}

	app.Run(os.Args)
}
