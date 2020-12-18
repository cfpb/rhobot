package main

import (
	"os"

	log "github.com/Sirupsen/logrus"
	"github.com/cfpb/rhobot/internal/config"
	"github.com/cfpb/rhobot/internal/gocd"
	"github.com/urfave/cli"
)

func main() {

	cli.VersionFlag = cli.BoolFlag{
		Name:  "print-version, V",
		Usage: "print only the version",
	}

	app := cli.NewApp()
	app.Name = "Rhobot"
	app.Usage = "Rhobot is a database development tool that uses DevOps best practices."
	app.Version = Version
	app.EnableBashCompletion = true

	conf := config.NewConfig()
	gocdServer := gocd.NewServerConfig(conf.GOCDHost, conf.GOCDPort, conf.GOCDUser, conf.GOCDPassword, conf.GOCDTimeout)

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
	templateFileFlag := cli.StringFlag{
		Name:  "template",
		Value: "",
		Usage: "path to the report template",
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
	schemaFlag := cli.StringFlag{
		Name:  "schema",
		Value: "",
		Usage: "which schema the healthchecks should be put into",
	}
	tableFlag := cli.StringFlag{
		Name:  "table",
		Value: "",
		Usage: "which table the healthchecks should be put into",
	}
	pipelineRunFlag := cli.StringFlag{
		Name:  "pipeline-run",
		Value: "0",
		Usage: "which run of the pipeline history",
	}
	stageRunFlag := cli.StringFlag{
		Name:  "stage-run",
		Value: "0",
		Usage: "which run of the stage history",
	}
	artifactPathFlag := cli.StringFlag{
		Name:  "artifact-path",
		Value: "cruise-output/console.log",
		Usage: "which artifact to get",
	}
	artifactSavePathFlag := cli.StringFlag{
		Name:  "save",
		Value: "",
		Usage: "where to save the fetched artifact",
	}

	app.Flags = []cli.Flag{logLevelFlag}
	app.Commands = []cli.Command{
		{
			Name: "healthchecks",
			Usage: "HEALTHCHECK_FILE " +
				"[--dburi DATABASE_URI] " +
				"[--report REPORT_FILE] [--email DISTRIBUTION_FILE]" +
				"[--schema SCHEMA] [--table TABLE]",
			Flags: []cli.Flag{
				reportFileFlag,
				templateFileFlag,
				dburiFlag,
				emailListFlag,
				schemaFlag,
				tableFlag,
			},
			Action: func(c *cli.Context) {
				updateLogLevel(c, conf)

				// variables to be populated by cli args
				var healthcheckPath string
				var reportPath string
				var templatePath string
				var emailListPath string
				var schema string
				var table string

				if c.Args().Get(0) != "" {
					healthcheckPath = c.Args().Get(0)
				} else {
					log.Error("You must provide the path to the healthcheck file.")
					return
				}
				log.Info("Running health checks from ", healthcheckPath)

				if c.String("dburi") != "" {
					conf.SetDBURI(c.String("dburi"))
				}
				log.Debug("DB_URI: ", conf.DBURI())

				if c.String("report") != "" {
					reportPath = c.String("report")
					log.Infof("Generating report at %v", reportPath)
				}

				if c.String("template") != "" {
					templatePath = c.String("template")
					log.Infof("Using template at %v", templatePath)
				}

				if c.String("email") != "" {
					emailListPath = c.String("email")
					log.Infof("Emailing report to %v", emailListPath)
				}

				if c.String("schema") != "" && c.String("table") != "" {
					schema = c.String("schema")
					table = c.String("table")
					log.Infof("Saving healthchecks to %v.%v", schema, table)
				}

				err := healthcheckRunner(conf, healthcheckPath, reportPath, templatePath, emailListPath, schema, table)
				if err != nil {
					log.Fatal(err)
				}
				log.Info("Healthchecks Success!")

			},
		},
		{
			Name: "artifacts",
			Usage: "PIPELINE STAGE JOB" +
				"[--pipeline-run PIPELINE_RUN] [--stage-run STAGE_RUN] " +
				"[--artifact-path ARTIFACT_PATH] [--save SAVE_PATH]",
			Flags: []cli.Flag{
				pipelineRunFlag,
				stageRunFlag,
				artifactPathFlag,
				artifactSavePathFlag,
			},
			Action: func(c *cli.Context) {
				updateLogLevel(c, conf)
				gocdServer = updateGOCDHost(c, conf)

				// variables to be populated by cli args
				var pipeline string
				var stage string
				var job string

				//optional cli flags
				var pipelineRun string
				var stageRun string
				var artifactPath string
				var artifactSavePath string

				if c.Args().Get(0) != "" && c.Args().Get(1) != "" && c.Args().Get(2) != "" {
					pipeline = c.Args().Get(0)
					stage = c.Args().Get(1)
					job = c.Args().Get(2)
				} else {
					log.Error("You must provide the PIPELINE STAGE and JOB of the artifact")
					return
				}

				if c.String(pipelineRunFlag.GetName()) != "" {
					pipelineRun = c.String(pipelineRunFlag.GetName())
					log.Debugf("%v: %v", pipelineRunFlag.GetName(), pipelineRun)
				}

				if c.String(stageRunFlag.GetName()) != "" {
					stageRun = c.String(stageRunFlag.GetName())
					log.Debugf("%v: %v", stageRunFlag.GetName(), stageRun)
				}

				if c.String(artifactPathFlag.GetName()) != "" {
					artifactPath = c.String(artifactPathFlag.GetName())
					log.Debugf("%v: %v", artifactPathFlag.GetName(), artifactPath)
				}

				if c.String(artifactSavePathFlag.GetName()) != "" {
					artifactSavePath = c.String(artifactSavePathFlag.GetName())
					log.Debugf("%v: %v", artifactSavePathFlag.GetName(), artifactSavePath)
				}

				getArtifact(gocdServer, pipeline, stage, job, pipelineRun, stageRun, artifactPath, artifactSavePath)
				log.Info("Success!")
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
						updateLogLevel(c, conf)
						gocdServer = updateGOCDHost(c, conf)

						if len(c.Args()) > 0 {
							path := c.Args()[0]
							group := c.Args().Get(1)
							log.Infof("Pushing config from %v to pipeline group %v...", path, group)
							if err := gocd.Push(gocdServer, path, group); err != nil {
								log.Fatal("Failed to push pipeline config: ", err)
							}
							log.Info("Success!")
						} else {
							log.Fatal("A path to the pipeline config to push is required.")
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
						updateLogLevel(c, conf)
						gocdServer = updateGOCDHost(c, conf)

						if len(c.Args()) > 0 {
							path := c.Args()[0]
							log.Infof("Pulling config from %v to %v...", gocdServer.URL(), path)
							if err := gocd.Pull(gocdServer, path); err != nil {
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
						gocdHostFlag,
					},
					Action: func(c *cli.Context) {
						updateLogLevel(c, conf)
						gocdServer = updateGOCDHost(c, conf)

						if len(c.Args()) > 1 {
							name := c.Args()[0]
							path := c.Args()[1]
							log.Infof("Cloning pipeline %v to %v...", name, path)
							if _, err := gocd.Clone(gocdServer, path, name); err != nil {
								log.Fatal("Failed to clone pipeline config: ", err)
							}
							log.Info("Success!")
						} else {
							log.Fatal("A pipeline name and a path to clone to are required.")
						}
					},
				},
				{
					Name:  "delete",
					Usage: "PIPELINE_NAME",
					Flags: []cli.Flag{
						gocdHostFlag,
					},
					Action: func(c *cli.Context) {
						updateLogLevel(c, conf)
						gocdServer = updateGOCDHost(c, conf)

						if len(c.Args()) > 0 {
							name := c.Args()[0]
							log.Infof("Deleteing pipeline %v...", name)
							if _, err := gocd.Delete(gocdServer, name); err != nil {
								log.Fatal("Failed to delete pipeline: ", err)
							}
							log.Info("Success!")
						} else {
							log.Fatal("A pipeline name is required.")
						}
					},
				},
			},
		},
	}

	app.Run(os.Args)
}
