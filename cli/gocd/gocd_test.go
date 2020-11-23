package gocd

import (
	"bytes"
	"encoding/json"
	"io"
	"os"
	"reflect"
	"testing"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/cfpb/rhobot/cli/config"
)

var gocdPipelineConfig []byte
var conf *config.Config
var server *Server

func init() {
	conf = config.NewConfig()
	conf.SetLogLevel("debug")

	// use no authentication for test
	conf.GOCDUser = ""
	conf.GOCDPassword = ""

	server = NewServerConfig(conf.GOCDHost, conf.GOCDPort, conf.GOCDUser, conf.GOCDPassword, conf.GOCDTimeout)

	buf := bytes.NewBuffer(nil)
	f, _ := os.Open("./test.json")
	io.Copy(buf, f)
	f.Close()
	gocdPipelineConfig = buf.Bytes()
}

func TestMarshalJSONHAL(t *testing.T) {
	pipeline, err := readPipelineJSONFromFile("./test.json")
	if err != nil {
		t.Error(err)
	}
	log.Debug("Pipeline Name: %+v\n", pipeline.Name)
	log.Debug("Pipeline Git URL: %v:%v\n", pipeline.Materials[0].Attributes.URL, pipeline.Materials[0].Attributes.Branch)
}

// TestUnmarshalFidelityLoss checks that data can be reserielized without fidelity loss
func TestUnmarshalFidelityLoss(t *testing.T) {
	var data interface{}
	err1 := json.Unmarshal(gocdPipelineConfig, &data)
	if err1 != nil {
		t.Error(err1)
	}

	gocdPipelineConfig2, _ := json.Marshal(data)
	var data2 interface{}
	err2 := json.Unmarshal(gocdPipelineConfig2, &data2)
	if err2 != nil {
		t.Error(err2)
	}

	if !reflect.DeepEqual(data, data2) {
		t.Error("not the same")
	}
}

func TestGocdPOST(t *testing.T) {
	etag, _, err := Exist(server, "test")
	if err == nil && etag != "" {
		log.Info("Cannot run TestGoCDPOST, 'test' pipeline already exists.")
		t.SkipNow()
	}

	pipeline, _ := readPipelineJSONFromFile("./test.json")
	pipelineConfig := PipelineConfig{"Dev", pipeline}

	_, err = server.pipelineConfigPOST(pipelineConfig)
	if err != nil {
		t.Error(err)
	}
}

func TestGocdFindPipeline(t *testing.T) {
	environment, err := server.environmentGET()
	if err != nil {
		t.Error(err)
	}

	environmentName := findPipelineInEnvironment(environment, "test")

	if environmentName != "" {
		log.Debug("Pipeline in environment with name: %+v", environmentName)
	} else {
		log.Debug("Pipeline not found in an environment")
	}
}

func TestGocdEnvironmentGET(t *testing.T) {
	_, err := server.environmentGET()
	if err != nil {
		t.Error(err)
	}
}

func TestGocdGET(t *testing.T) {
	_, _, err := server.pipelineGET("test")
	if err != nil {
		t.Error(err)
	}
}

// //The following 3 tests require a pipeline to have a run history
// //thus is commented out for testing on TravisCI
// //Future scaffolding will be needed on travis to add an agent,
// //add the agent and pipeline to an environment,
// //unpause, and run a pipeline to get a run history
// //Uncomment for testing on local machine
//
// func TestGocdHistoryGET(t *testing.T) {
// 	counterMap, err := History(server, "test")
// 	spew.Dump(counterMap)
// 	if err != nil {
// 		t.Error(err)
// 	}
// }
//
// func TestGocdArtifactGET(t *testing.T) {
// 	runsIDMap, err := History(server, "test")
// 	if err != nil {
// 		t.Error(err)
// 	}
//
// 	artifactBuffer, err := Artifact(
// 		server,
// 		"test", runsIDMap["p_test"],
// 		"hello", runsIDMap["s_hello"],
// 		"world", "cruise-output/console.log")
// 	if err != nil {
// 		t.Error(err)
// 	}
//
// 	artifactBuffer.WriteTo(os.Stdout)
// }
//
// func TestGocdEnvironmentPATCH(t *testing.T) {
// 	err = server.environmentPATCH("test", "future_named_env")
// 	if err != nil {
// 		t.Error(err)
// 	}
// }

func TestExist(t *testing.T) {
	etag, _, err := Exist(server, "test")

	if err != nil {
		t.Error(err)
	}

	if etag == "" {
		t.Error("test does not exist as a gocd pipeline")
	}
}

func TestGocdPUT(t *testing.T) {
	pipeline, etag, _ := server.pipelineGET("test")

	// The Index of the STRANGE Environment Variable could potentially change between update
	strangeIndex := -1
	for i, envVar := range pipeline.EnvironmentVariables {
		if envVar.Name == "STRANGE" {
			strangeIndex = i
			break
		}
	}
	if strangeIndex == -1 {
		log.Debugf("EnvironmentVariables: %+v\n", pipeline.EnvironmentVariables)
		t.Fatal("STRANGE environment variable not found")
	}

	//Update Original Value to Time Value
	pipeline, etag, _ = server.pipelineGET("test")
	strangeEnvVarA := pipeline.EnvironmentVariables[strangeIndex]
	pipeline.EnvironmentVariables[strangeIndex].Value = time.Now().UTC().String()
	pipeline, _ = server.pipelineConfigPUT(pipeline, etag)

	//Update Time Value to Original Value
	pipeline, etag, _ = server.pipelineGET("test")
	strangeEnvVarB := pipeline.EnvironmentVariables[strangeIndex]
	pipeline.EnvironmentVariables[strangeIndex].Value = strangeEnvVarA.Value
	pipeline, _ = server.pipelineConfigPUT(pipeline, etag)

	pipeline, _, _ = server.pipelineGET("test")
	strangeEnvVarC := pipeline.EnvironmentVariables[strangeIndex]
	log.Debugf("STRANGE VALUE A: %+v\n", strangeEnvVarA)
	log.Debugf("STRANGE VALUE B: %+v\n", strangeEnvVarB)
	log.Debugf("STRANGE VALUE C: %+v\n", strangeEnvVarC)

	if strangeEnvVarA == strangeEnvVarB {
		t.Error("STRANGE environment variable was not changed")
	}

	if strangeEnvVarA != strangeEnvVarC {
		t.Error("STRANGE environment variable was not reset")
	}
}

func TestGocdPUTEncrypt(t *testing.T) {
	pipeline, etag, _ := server.pipelineGET("test")

	// The Index of the STRANGE Environment Variable could potentially change between update
	strangeIndex := -1
	for i, envVar := range pipeline.EnvironmentVariables {
		if envVar.Name == "STRANGE" {
			strangeIndex = i
			break
		}
	}
	if strangeIndex == -1 {
		log.Debugf("EnvironmentVariables: %+v\n", pipeline.EnvironmentVariables)
		t.Fatal("STRANGE environment variable not found")
	}

	//Update Original Value to Time Value
	pipeline, etag, _ = server.pipelineGET("test")
	strangeEnvVarA := pipeline.EnvironmentVariables[strangeIndex]
	pipeline.EnvironmentVariables[strangeIndex].Value = "something"
	pipeline.EnvironmentVariables[strangeIndex].EncryptedValue = ""
	pipeline.EnvironmentVariables[strangeIndex].Secure = true
	pipeline, _ = server.pipelineConfigPUT(pipeline, etag)

	//Update Time Value to Original Value
	pipeline, etag, _ = server.pipelineGET("test")
	strangeEnvVarB := pipeline.EnvironmentVariables[strangeIndex]
	pipeline.EnvironmentVariables[strangeIndex].Value = strangeEnvVarA.Value
	pipeline.EnvironmentVariables[strangeIndex].EncryptedValue = ""
	pipeline.EnvironmentVariables[strangeIndex].Secure = false
	pipeline, err := server.pipelineConfigPUT(pipeline, etag)

	if err != nil {
		log.Debugf("Put Error: %+v\n", err.Error())
		t.Error("Put failed")
	}

	pipeline, _, _ = server.pipelineGET("test")
	strangeEnvVarC := pipeline.EnvironmentVariables[strangeIndex]
	log.Debugf("STRANGE VALUE A: %+v\n", strangeEnvVarA)
	log.Debugf("STRANGE VALUE B: %+v\n", strangeEnvVarB)
	log.Debugf("STRANGE VALUE C: %+v\n", strangeEnvVarC)

	if strangeEnvVarA == strangeEnvVarB {
		t.Error("STRANGE environment variable was not changed")
	}

	if strangeEnvVarA != strangeEnvVarC {
		t.Error("STRANGE environment variable was not reset")
	}
}

func TestGocdTimeout(t *testing.T) {

	serverA := NewServerConfig(conf.GOCDHost, conf.GOCDPort, conf.GOCDUser, conf.GOCDPassword, "120")
	serverB := &Server{
		Host:     serverA.Host,
		Port:     serverA.Port,
		User:     serverA.User,
		Password: serverA.Password,
		Timeout:  time.Duration(1), //1ns
	}

	etag, _, err := Exist(serverA, "test")
	if etag == "" {
		t.Error("test does not exist as a gocd pipeline")
	}
	if err != nil {
		t.Error("threw an error but should not have", err)
	}

	etag, _, err = Exist(serverB, "test")
	if etag != "" {
		t.Error("got an etag but should not have", etag)
	}
	if err == nil {
		t.Error("did not Timeout but should have", err)
	}

}

func TestGocdDELETE(t *testing.T) {
	_, err := server.pipelineDELETE("test")
	if err != nil {
		t.Error(err)
	}
}
