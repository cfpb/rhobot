package gocd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"reflect"
	"testing"
	"time"
)

var gocdPipelineConfig []byte
var gocdHost string
var gocdPort string
var gocdURL string

func init() {
	buf := bytes.NewBuffer(nil)
	f, _ := os.Open("./test.json")
	io.Copy(buf, f)
	f.Close()
	gocdPipelineConfig = buf.Bytes()

	gocdHost = os.Getenv("GOCD_HOST")
	if gocdHost == "" {
		gocdHost = "http://localhost"
	}

	gocdPort = os.Getenv("GOCD_PORT")
  if gocdPort == "" {
	  gocdPort = "8153"
  }

	gocdURL = fmt.Sprintf("%s:%s", gocdHost, gocdPort)
}

func TestMarshalJSONHAL(t *testing.T) {
	fmt.Println("TestMarshalJSONHAL")
	pipeline, err := ReadPipelineJSONFromFile("./test.json")
	if err != nil {
		t.Error(err)
	}
	fmt.Printf("Pipeline Name: %+v\n", pipeline.Name)
	fmt.Printf("Pipeline Git URL: %v:%v\n", pipeline.Materials[0].Attributes.URL, pipeline.Materials[0].Attributes.Branch)
}

// TestUnmarshalFidelityLoss checks that data can be reserielized without fidelity loss
func TestUnmarshalFidelityLoss(t *testing.T) {
	fmt.Println("TestUnmarshalFidelityLoss")
	data, err1 := unmarshalPipeline(gocdPipelineConfig)
	if err1 != nil {
		t.Error(err1)
	}

	gocdPipelineConfig2, _ := json.Marshal(data)
	data2, err2 := unmarshalPipeline(gocdPipelineConfig2)
	if err2 != nil {
		t.Error(err2)
	}

	if !reflect.DeepEqual(data, data2) {
		t.Error("not the same")
	}
}

func TestGocdPOST(t *testing.T) {
	fmt.Println("TestGocdPOST")
	pipeline, _ := ReadPipelineJSONFromFile("./test.json")
	pipelineConfig := PipelineConfig{"Dev", pipeline}
	_, err := pipelineConfigPOST(gocdURL, pipelineConfig)
	if err != nil {
		t.Error(err)
	}
}

func TestGocdGET(t *testing.T) {
	fmt.Println("TestGocdGET")
	_, _, err := pipelineGET(gocdURL, "test")
	if err != nil {
		t.Error(err)
	}
}

func TestExist(t *testing.T) {
	fmt.Println("TestExist")
	etag := Exist(gocdURL, "test")
	if etag == "" {
		t.Error("test does not exist as a gocd pipeline")
	}
}

func TestGocdPUT(t *testing.T) {
	fmt.Println("TestGocdPUT")
	pipeline, etag, _ := pipelineGET(gocdURL, "test")

	// The Index of the STRANGE Environment Variable could potentially change between update
	strangeIndex := -1
	for i, envVar := range pipeline.EnvironmentVariables {
		if envVar.Name == "STRANGE" {
			strangeIndex = i
			break
		}
	}
	if strangeIndex == -1 {
		fmt.Printf("EnvironmentVariables: %+v\n", pipeline.EnvironmentVariables)
		t.Fatal("STRANGE environment variable not found")
	}

	//Update Original Value to Time Value
	pipeline, etag, _ = pipelineGET(gocdURL, "test")
	strangeEnvVarA := pipeline.EnvironmentVariables[strangeIndex]
	pipeline.EnvironmentVariables[strangeIndex].Value = time.Now().UTC().String()
	pipeline, _ = pipelineConfigPUT(gocdURL, pipeline, etag)

	//Update Time Value to Original Value
	pipeline, etag, _ = pipelineGET(gocdURL, "test")
	strangeEnvVarB := pipeline.EnvironmentVariables[strangeIndex]
	pipeline.EnvironmentVariables[strangeIndex].Value = strangeEnvVarA.Value
	pipeline, _ = pipelineConfigPUT(gocdURL, pipeline, etag)

	pipeline, etag, _ = pipelineGET(gocdURL, "test")
	strangeEnvVarC := pipeline.EnvironmentVariables[strangeIndex]
	fmt.Printf("STRANGE VALUE A: %+v\n", strangeEnvVarA)
	fmt.Printf("STRANGE VALUE B: %+v\n", strangeEnvVarB)
	fmt.Printf("STRANGE VALUE C: %+v\n", strangeEnvVarC)

	if strangeEnvVarA == strangeEnvVarB {
		t.Error("STRANGE environment variable was not changed")
	}

	if strangeEnvVarA != strangeEnvVarC {
		t.Error("STRANGE environment variable was not reset")
	}
}
