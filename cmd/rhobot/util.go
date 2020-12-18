package main

import (
	"errors"
	"io/ioutil"
	"os"
	"strconv"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/cfpb/rhobot/internal/config"
	"github.com/cfpb/rhobot/internal/database"
	"github.com/cfpb/rhobot/internal/gocd"
	"github.com/cfpb/rhobot/internal/healthcheck"
	"github.com/cfpb/rhobot/internal/report"
	"github.com/davecgh/go-spew/spew"
	"github.com/urfave/cli"
)

//TODO: Fatal exits should happen in the main loop, not the utils, pass errs up

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

func healthcheckRunner(config *config.Config, healthcheckPath string, reportPath string, templatePath string, emailListPath string, hcSchema string, hcTable string) (err error) {
	healthChecks, err := healthcheck.ReadHealthCheckYAMLFromFile(healthcheckPath)
	if err != nil {
		log.Fatal("Failed to read healthchecks: ", err)
	}
	cxn := database.GetPGConnection(config.DBURI())

	results, HCerrs := healthChecks.PreformHealthChecks(cxn)
	numErrors, numWarnings, fatal := healthcheck.EvaluateHCErrors(HCerrs)

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
		"status":    healthcheck.StatusHealthchecks(numErrors, numWarnings, fatal),
		"schema":    hcSchema,
		"table":     hcTable,
	}
	rs := report.Set{Elements: elements, Metadata: metadata}

	// Load template if provided
	var template string
	if templatePath != "" {
		data, readErr := ioutil.ReadFile(templatePath)
		if readErr != nil {
			log.Fatal("Failed to read template: ", err)
		}
		template = string(data)
	} else {
		template = healthcheck.TemplateHealthcheckHTML
	}

	// Write report to file
	if reportPath != "" {
		prr := report.NewPongo2ReportRunnerFromString(template, true)
		reader, _ := prr.ReportReader(rs)
		fhr := report.FileHandler{Filename: reportPath}
		err = fhr.HandleReport(reader)
		if err != nil {
			log.Error("error writing report to file: ", err)
		}
	}

	// Email report
	if emailListPath != "" {
		prr := report.NewPongo2ReportRunnerFromString(template, true)
		df, err := report.ReadDistributionFormatYAMLFromFile(emailListPath)
		if err != nil {
			log.Fatal("Failed to read distribution format: ", err)
		}

		for _, level := range report.LogLevelArray {

			subjectStr := healthcheck.SubjectHealthcheck(healthChecks.Name, config.PgDatabase, config.PgHost, level, numErrors, numWarnings, fatal)

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
		prr := report.NewPongo2ReportRunnerFromString(healthcheck.TemplateHealthcheckPostgres, false)
		pgr := report.PGHandler{Cxn: cxn}
		reader, err := prr.ReportReader(rs)
		err = pgr.HandleReport(reader)
		if err != nil {
			log.Errorf("Failed to save healthchecks to PG database\n%s", err)
		}
	}

	if numErrors > 0 || fatal == true {
		// log.Panic("Healthchecks Failed:\n", spew.Sdump(HCerrs))
		err := errors.New("Healthchecks Failed")
		return err
	}
	return nil
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
