package gocd

import (
	"bytes"
	"encoding/json"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"

	log "github.com/Sirupsen/logrus"
)

// History gets the run history of a pipeline of a given name exist. returns map
func History(server *Server, name string) (latestRuns map[string]int, err error) {

	//Get pipeline history if it exist
	historyJSON, err := server.historyGET(name)
	if err != nil {
		log.Fatalf("Could not find run history for pipeline: %v", name)
	}

	//setup temp vatiables
	var responseMap map[string]*json.RawMessage
	var pipelineArr []map[string]*json.RawMessage
	var pipelineLatest map[string]*json.RawMessage
	var pipelineCounter int
	var stages []map[string]*json.RawMessage

	//unmarshal raw json into temp variables.
	_ = json.Unmarshal(historyJSON.Bytes(), &responseMap)
	_ = json.Unmarshal(*responseMap["pipelines"], &pipelineArr)
	pipelineLatest = pipelineArr[0]
	_ = json.Unmarshal(*pipelineLatest["counter"], &pipelineCounter)

	//setup return variable
	latestRuns = make(map[string]int)
	latestRuns["p_"+name] = pipelineCounter

	//setup temp variables for loop
	_ = json.Unmarshal(*pipelineLatest["stages"], &stages)
	var stageName string
	var stageCounterStr string
	var stageCounterInt int

	//loop through and parse stages JSON to get "counter" vatiables
	for _, stage := range stages {
		_ = json.Unmarshal(*stage["name"], &stageName)
		_ = json.Unmarshal(*stage["counter"], &stageCounterStr)
		stageCounterInt, err = strconv.Atoi(stageCounterStr)
		latestRuns["s_"+stageName] = stageCounterInt
	}

	return
}

// Artifact gets an Areifact from a pipeline / stage / job
func Artifact(server *Server, pipelineName string, pipelineID int, stageName string, stageID int, jobName string, artifactPath string) (fileBytes *bytes.Buffer, err error) {
	fileBytes, err = server.artifactGET(
		pipelineName, pipelineID,
		stageName, stageID,
		jobName, artifactPath)
	return
}

// Push takes a pipeline from a file and sends it to GoCD
func Push(server *Server, path string, group string) (err error) {
	localPipeline, err := readPipelineJSONFromFile(path)
	if err != nil {
		return
	}

	etag, remotePipeline, err := Exist(server, localPipeline.Name)
	if err != nil {
		log.Info(err.Error())
	}

	Compare(localPipeline, remotePipeline, path)

	if etag == "" {
		pipelineConfig := PipelineConfig{group, localPipeline}
		_, err = server.pipelineConfigPOST(pipelineConfig)
	} else {
		_, err = server.pipelineConfigPUT(localPipeline, etag)
	}
	return
}

// Pull reads pipeline from a file, finds it on GoCD, and updates the file
func Pull(server *Server, path string) (err error) {
	localPipeline, err := readPipelineJSONFromFile(path)
	if err != nil {
		return
	}

	name := localPipeline.Name
	remotePipeline, err := Clone(server, path, name)

	Compare(localPipeline, remotePipeline, path)

	return
}

// Attempt to delete pipeline if name matches existing
func Delete(server *Server, pipelineName string) (pipeline Pipeline, err error) {
	environment, err := server.environmentGET()
	if err != nil {
		return
	}

	environmentName := findPipelineInEnvironment(environment, pipelineName)

	if environmentName != "" {
		log.Infof("Pipeline found in environment, removing from environment: %v", environmentName)

		err = server.environmentPATCH(pipelineName, environmentName)
		if err != nil {
			log.Info("Environment not patched")
			return
		}
	}

	pipeline, err = server.pipelineDELETE(pipelineName)
	if err != nil {
		return
	}

	return
}

// Exist checks if a pipeline of a given name exist, returns it's etag or an empty string
func Exist(server *Server, name string) (etag string, pipeline Pipeline, err error) {
	pipeline, etag, err = server.pipelineGET(name)
	return
}

// Clone finds a pipeline by name on GoCD and saves it to a file
func Clone(server *Server, path string, name string) (pipeline Pipeline, err error) {
	pipeline, _, err = server.pipelineGET(name)
	if err != nil {
		return
	}

	err = writePipeline(path, pipeline)
	return
}

// Compare saves copies of the local and remote pipeline if different
func Compare(localPipeline Pipeline, remotePipeline Pipeline, path string) {

	if !reflect.DeepEqual(localPipeline, remotePipeline) {
		log.Warn("Local and Remote are different")

		filepath := strings.TrimSuffix(path, filepath.Ext(path))
		localBakPath := filepath + ".local.bak.json"
		remoteBakPath := filepath + ".remote.bak.json"

		log.Info("Saving Local Backup: ", localBakPath)
		errLocal := writePipeline(localBakPath, localPipeline)
		log.Info("Saving Remote Backup: ", remoteBakPath)
		errRemote := writePipeline(remoteBakPath, remotePipeline)

		if errLocal != nil {
			log.Warn("Error while writing backup for local pipeline: ", errLocal)
		}
		if errRemote != nil {
			log.Warn("Error while writing backup for local pipeline: ", errLocal)
		}

	}
}
