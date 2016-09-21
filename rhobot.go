package main

import (
	"os"
	"strconv"
	"strings"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/cfpb/rhobot/config"
	"github.com/cfpb/rhobot/database"
	"github.com/cfpb/rhobot/gocd"
	"github.com/cfpb/rhobot/healthcheck"
	"github.com/cfpb/rhobot/report"
	"github.com/davecgh/go-spew/spew"
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

				if c.String("email") != "" {
					emailListPath = c.String("email")
					log.Infof("Emailing report to %v", emailListPath)
				}

				if c.String("schema") != "" && c.String("table") != "" {
					schema = c.String("schema")
					table = c.String("table")
					log.Infof("Saving healthchecks to %v.%v", schema, table)
				}

				healthcheckRunner(conf, healthcheckPath, reportPath, emailListPath, schema, table)
				log.Info("Success!")
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
			},
		},
	}

	app.Run(os.Args)
}

func updateLogLevel(c *cli.Context, config *config.Config) {
	if c.GlobalString("loglevel") != "" {
		config.SetLogLevel(c.GlobalString("loglevel"))
	}
}

func updateGOCDHost(c *cli.Context, config *config.Config) (gocdServer *gocd.Server) {
	if c.String("host") != "" {
		config.SetGoCDHost(c.String("host"))
	}
	gocdServer = gocd.NewServerConfig(config.GOCDHost, config.GOCDPort, config.GOCDUser, config.GOCDPassword, config.GOCDTimeout)
	return
}

func healthcheckRunner(config *config.Config, healthcheckPath string, reportPath string, emailListPath string, hcSchema string, hcTable string) {
	healthChecks, err := healthcheck.ReadHealthCheckYAMLFromFile(healthcheckPath)
	if err != nil {
		log.Fatal("Failed to read healthchecks: ", err)
	}
	cxn := database.GetPGConnection(config.DBURI())

	results, HCerrs := healthChecks.PreformHealthChecks(cxn)
	numErrors := 0
	fatal := false
	for _, hcerr := range HCerrs {
		if strings.Contains(strings.ToUpper(hcerr.Err), "FATAL") {
			fatal = true
		}
		if strings.Contains(strings.ToUpper(hcerr.Err), "ERROR") {
			numErrors = numErrors + 1
		}
	}

	var elements []report.Element
	for _, val := range results {
		elements = append(elements, val)
	}

	// Make Templated report
	metadata := map[string]interface{}{
		"name":      healthChecks.Name,
		"db_name":   config.PgDatabase,
		"footer":    healthcheck.FooterHealthcheck,
		"timestamp": time.Now().Format(time.ANSIC),
		"status":    healthcheck.StatusHealthchecks(numErrors, fatal),
		"schema":    hcSchema,
		"table":     hcTable,
	}
	rs := report.Set{Elements: elements, Metadata: metadata}

	// Write report to file
	if reportPath != "" {
		prr := report.NewPongo2ReportRunnerFromString(healthcheck.TemplateHealthcheckHTML)
		reader, _ := prr.ReportReader(rs)
		fhr := report.FileHandler{Filename: reportPath}
		err = fhr.HandleReport(reader)
		if err != nil {
			log.Error("error writing report to PG database: ", err)
		}
	}

	// Email report
	if emailListPath != "" {
		prr := report.NewPongo2ReportRunnerFromString(healthcheck.TemplateHealthcheckHTML)
		df, err := report.ReadDistributionFormatYAMLFromFile(emailListPath)
		if err != nil {
			log.Fatal("Failed to read distribution format: ", err)
		}

		for _, level := range report.LogLevelArray {

			subjectStr := healthcheck.SubjectHealthcheck(healthChecks.Name, config.PgDatabase, config.PgHost, level, numErrors, fatal)

			logFilteredSet := report.FilterReportSet(rs, level)
			reader, _ := prr.ReportReader(logFilteredSet)
			recipients := df.GetEmails(level)

			if recipients != nil && len(recipients) != 0 && len(logFilteredSet.Elements) != 0 {
				log.Infof("Send %s to: %v", subjectStr, recipients)
				ehr := report.EmailHandler{
					SMTPHost:    config.SMTPHost,
					SMTPPort:    config.SMTPPort,
					SenderEmail: config.SMTPEmail,
					SenderName:  config.SMTPName,
					Subject:     subjectStr,
					Recipients:  recipients,
					HTML:        true,
				}
				err = ehr.HandleReport(reader)
				if err != nil {
					log.Error("Failed to email report: ", err)
				}
			}
		}
	}

	if hcSchema != "" && hcTable != "" {
		prr := report.NewPongo2ReportRunnerFromString(healthcheck.TemplateHealthcheckPostgres)
		pgr := report.PGHandler{Cxn: cxn}
		reader, err := prr.ReportReader(rs)
		err = pgr.HandleReport(reader)
		if err != nil {
			log.Errorf("Failed to save healthchecks to PG database\n%s", err)
		}
	}

	// Bad Exit
	if HCerrs != nil {
		log.Fatal("Healthchecks Failed:\n", spew.Sdump(HCerrs))
	}
}

func getArtifact(gocdServer *gocd.Server, pipeline string, stage string, job string,
	pipelineRun string, stageRun string, artifactPath string, artifactSavePath string) {

	//parse or get latest run numbers for pipeline and stage
	var pipelineRunNum, stageRunNum int = 0, 0
	var pipelineOk, stageOk bool = true, true
	var pipelineErr, stageErr error

	counterMap, err := gocd.History(gocdServer, pipeline)
	if err != nil {
		log.Fatalf("Could not find run history for pipeline: %v", pipeline)
	}
	log.Debug(spew.Sdump(counterMap))

	if pipelineRun == "0" {
		pipelineRunNum, pipelineOk = counterMap["p_"+pipeline]
	} else {
		pipelineRunNum, pipelineErr = strconv.Atoi(pipelineRun)
	}

	if stageRun == "0" {
		stageRunNum, stageOk = counterMap["s_"+stage]
	} else {
		stageRunNum, stageErr = strconv.Atoi(stageRun)
	}

	if !pipelineOk && !stageOk {
		log.Fatalf("Pipeline: \"%v\" and Stage: \"%v\" not found in pipeline history", pipeline, stage)
	}
	if pipelineErr != nil || stageErr != nil {
		log.Fatalf("Pipeline: %v and Stage: %v could not be parsed to integers", pipelineRun, stageRun)
	}

	//fetch artifact
	log.Infof("getting GoCD Artifact - pipeline:%v , pipelineRunNum:%v , stage:%v , stageRunNum:%v , job:%v , artifactPath:%v",
		pipeline, pipelineRunNum, stage, stageRunNum, job, artifactPath)
	artifactBuffer, err := gocd.Artifact(gocdServer, pipeline, pipelineRunNum, stage, stageRunNum, job, artifactPath)
	if err != nil {
		log.Fatalf("Failed to fetch artifact: %v", artifactPath)
	}

	//write to file or log
	if artifactSavePath == "" {
		artifactBuffer.WriteTo(os.Stdout)
	} else {
		f, err := os.Create(artifactSavePath)
		if err != nil {
			log.Fatalf("Failed to create file: %v", artifactSavePath)
		}

		_, err = artifactBuffer.WriteTo(f)
		if err != nil {
			log.Fatalf("Failed to write to file: %v", artifactSavePath)
		}
		defer f.Close()
	}

}
