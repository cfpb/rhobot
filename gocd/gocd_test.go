package gocd

import (
	"bytes"
	"encoding/json"
	"fmt"
	//"github.com/davecgh/go-spew/spew"
	"io"
	"os"
	"reflect"
	"testing"
)

var gocd_pipeline_config []byte

func init() {
	buf := bytes.NewBuffer(nil)
	f, _ := os.Open("./test.json")
	io.Copy(buf, f)
	f.Close()
	gocd_pipeline_config = buf.Bytes()
}

func TestMarshalJSONHAL(t *testing.T) {
	fmt.Println("TestMarshalJSONHAL")
	pipeline, err := ReadPipelineJSONFromFile("./test.json")
    if err != nil {
        t.Error(err)
	}
	//spew.Dump(pipeline)
	fmt.Printf("Pipeline Name: %+v\n", pipeline.Name)
	fmt.Printf("Pipeline Git URL: %v:%v\n", pipeline.Materials[0].Attributes.URL, pipeline.Materials[0].Attributes.Branch)
}

// TestUnmarshalFidelityLoss checks that data can be reserielized without fidelity loss
func TestUnmarshalFidelityLoss(t *testing.T) {
    fmt.Println("TestUnmarshalFidelityLoss")
	data, err1 := UnmarshalPipeline(gocd_pipeline_config)
	if err1 != nil {
        t.Error(err1)
	}

	gocd_pipeline_config2, _ := json.Marshal(data)
	data2, err2 := UnmarshalPipeline(gocd_pipeline_config2)
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
	_, err := pipelineConfigPOST("http://localhost:8153", pipelineConfig)
	if err != nil {
        t.Error(err)
	}
}

func TestGocdGET(t *testing.T) {
	fmt.Println("TestGocdGET")
	_, err := pipelineGET("http://localhost:8153", "test")
	if err != nil {
        t.Error(err)
	}
}
