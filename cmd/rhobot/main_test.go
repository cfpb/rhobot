package main

import (
	"os"
	"testing"
	"time"

	log "github.com/Sirupsen/logrus"

	"github.com/cfpb/rhobot/internal/config"
	"github.com/cfpb/rhobot/internal/database"
	"github.com/cfpb/rhobot/internal/healthcheck"
	"github.com/cfpb/rhobot/internal/report"
)

var conf *config.Config

func init() {
	conf = config.NewConfig()
	conf.SetLogLevel("debug")
}

func TestPostgresHealthCheckReporting(t *testing.T) {
	cxn := database.GetPGConnection(conf.DBURI())
	healthChecks, _ := healthcheck.ReadHealthCheckYAMLFromFile("healthcheck/healthchecksAll.yml")
	results, _ := healthChecks.PreformHealthChecks(cxn)
	var elements []report.Element
	for _, val := range results {
		elements = append(elements, val)
	}
	metadata := map[string]interface{}{
		"name":      healthChecks.Name,
		"schema":    "public",
		"table":     "healthchecks",
		"timestamp": time.Now().Format(time.ANSIC),
	}
	rs := report.Set{Elements: elements, Metadata: metadata}

	prr := report.NewPongo2ReportRunnerFromString(healthcheck.TemplateHealthcheckPostgres, false)
	pgr := report.PGHandler{Cxn: cxn}
	reader, err := prr.ReportReader(rs)
	err = pgr.HandleReport(reader)
	if err != nil {
		t.Fatalf("error writing report to PG database\n%s", err)
	}
}

func TestTemplateHealthCheckReporting(t *testing.T) {

	args := []string{
		"rhobot",
		"healthchecks", "healthcheck/healthchecksAll.yml",
		"-template", "healthcheck/templateHealthcheck.html",
		"-report", "testReportAll.html"}
	os.Args = args
	// TODO: need to find a better way to test fatal exits
	// assert.Panics(t, main, "The healthchecks did not cause a panic")
}

func TestPostgresHealthCheckEscape(t *testing.T) {
	cxn := database.GetPGConnection(conf.DBURI())
	healthChecks, _ := healthcheck.ReadHealthCheckYAMLFromFile("healthcheck/healthchecksEscape.yml")
	results, _ := healthChecks.PreformHealthChecks(cxn)
	var elements []report.Element
	for _, val := range results {
		elements = append(elements, val)
	}
	metadata := map[string]interface{}{
		"name":      healthChecks.Name,
		"schema":    "public",
		"table":     "healthchecks",
		"timestamp": time.Now().Format(time.ANSIC),
	}
	rs := report.Set{Elements: elements, Metadata: metadata}

	prr := report.NewPongo2ReportRunnerFromString(healthcheck.TemplateHealthcheckPostgres, false)
	pgr := report.PGHandler{Cxn: cxn}
	reader, err := prr.ReportReader(rs)
	err = pgr.HandleReport(reader)
	if err != nil {
		t.Fatalf("error writing report to PG database\n%s", err)
	}
}

func TestLogLevelingFiltering(t *testing.T) {
	cxn := database.GetPGConnection(conf.DBURI())
	healthChecks, _ := healthcheck.ReadHealthCheckYAMLFromFile("healthcheck/healthchecksAll.yml")
	results, _ := healthChecks.PreformHealthChecks(cxn)
	var elements []report.Element
	for _, val := range results {
		elements = append(elements, val)
	}
	metadata := map[string]interface{}{
		"name": healthChecks.Name,
	}
	originalReportSet := report.Set{Elements: elements, Metadata: metadata}
	TestLogLevelResults := map[string]int{
		"Debug": 6,
		"Info":  5,
		"Warn":  4,
		"Error": 2,
		"Fatal": 1,
	}

	for _, level := range report.LogLevelArray {
		log.Infof("Report Filter LogLevel %v", level)
		logFilteredSet := report.FilterReportSet(originalReportSet, level)
		prr := report.JSONReportRunner{}
		reader, _ := prr.ReportReader(logFilteredSet)
		phr := report.PrintHandler{}
		_ = phr.HandleReport(reader)

		if len(logFilteredSet.GetElementArray()) != TestLogLevelResults[level] {
			t.Fatalf("wrong number of healthchecks in report")
		}
	}

}

//The following are rhobot CLI test

func TestCLI_Healthchecks(t *testing.T) {
	//TODO create fail verification test

	args := []string{"rhobot", "healthchecks", "healthcheck/healthchecksTest.yml"}
	os.Args = args
	main()
}

func TestCLI_PG_Healthchecks(t *testing.T) {

	//clear hctest
	cxn := database.GetPGConnection(conf.DBURI())
	result, err := cxn.Exec("DROP TABLE IF EXISTS public.hctest")
	if err != nil {
		t.Fatalf("error dropping public.hctest\n%s", err)
	}

	args := []string{"rhobot", "healthchecks", "healthcheck/healthchecksTest.yml",
		"--schema", "public", "--table", "hctest"}
	os.Args = args
	main()

	//make sure hctest only has 3 rows
	row, err := cxn.Query("SELECT count(*) FROM public.hctest;")
	defer row.Close()
	if err != nil {
		t.Fatalf("error selecting public.hctest\n%s", err)
	} else {

		if row.Next() {
			var count int
			err = row.Scan(&count)

			if count != 3 {
				t.Fatal("public.hctest count should have been 3")
			}

		} else {
			t.Fatalf("no results in selecting public.hctest\n%s", err)
		}

	}

	log.Info(result)
}

func TestCLI(t *testing.T) {
	//All the following should exit with 0

	os.Args = []string{"rhobot"}
	main()

	os.Args = []string{"rhobot", "-V"}
	main()

	os.Args = []string{"rhobot", "pipeline"}
	main()
}
